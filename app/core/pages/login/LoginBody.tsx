import {Alert, Button, Chip, Divider, ToggleButton, ToggleButtonGroup} from "@mui/material"
import {useState} from "react"

import {useRouterLogin} from "../../../features/auth/api/AuthHook"
import {LogoutButton} from "../../../features/auth/component/LogoutButton"
import {UserAuthType} from "../../../features/user/api/UserType"
import {PageStartupBox} from "../../../shared/component/box/PageStartupBox"
import {KeyEnterInput} from "../../../shared/component/input/KeyEnterInput"
import {SxPropsMap} from "../../../shared/helper/HelperType"

const SX: SxPropsMap = {
    alert: {width: "100%", padding: "0 20px", justifyContent: "center"},
    label: {color: "text.secondary", fontWeight: "bold", fontSize: "25px"},
    types: {display: "flex", gap: 1, width: "100%", flexDirection: "column"},
}

type Props = {
    supported: UserAuthType[],
    error?: string,
}

export function LoginBody(props: Props) {
    const {supported, error} = props
    const [auth, setAuth] = useState<UserAuthType>(supported[0] ?? UserAuthType.BASIC)
    const [show, setShow] = useState(false)
    const [username, setUsername] = useState("")
    const [password, setPass] = useState("")
    const login = useRouterLogin()

    return (
        <PageStartupBox header={"Authentication"} renderFooter={renderFooter()} position={"start"}>
            {renderButtons()}
            <Divider flexItem={true}><Chip label={"✦"} size={"small"} onClick={() => setShow(!show)}/></Divider>
            {show && error && <Alert sx={SX.alert} severity={"warning"} icon={false}>{error}</Alert>}
            {renderCreds()}
        </PageStartupBox>
    )

    function renderCreds() {
        if ([UserAuthType.OIDC].includes(auth)) return null

        return (
            <>
                <KeyEnterInput
                    label={"Username"}
                    onChange={(e) => setUsername(e.target.value)}
                />
                <KeyEnterInput
                    label={"Password"}
                    hidden
                    onChange={(e) => setPass(e.target.value)}
                    onEnterPress={() => handleLogin(UserAuthType.BASIC)}
                />
            </>
        )
    }

    function renderFooter() {
        if (!supported.length) return <LogoutButton/>

        return (
            <Button onClick={() => handleLogin(auth)} loading={login.isPending} fullWidth={true}>
                SIGN IN
            </Button>
        )
    }

    function renderButtons() {
        if (supported.length < 2) return null

        return (
            <ToggleButtonGroup value={auth} exclusive={true} fullWidth={true}>
                {supported.map(type => (
                    <ToggleButton value={type} key={type} onClick={() => setAuth(type)}>
                        {type.toUpperCase()}
                    </ToggleButton>
                ))}
            </ToggleButtonGroup>
        )
    }

    function handleLogin(type: UserAuthType) {
        const subject = type === UserAuthType.OIDC ? undefined : {username, password}
        login.mutate({type, subject})
    }
}
