import {Alert, Button, Divider, ToggleButton, ToggleButtonGroup} from "@mui/material"
import {useState} from "react"

import {useRouterLogin, useRouterLogout} from "../../../api/auth/hook"
import {AuthType} from "../../../api/auth/type"
import {SxPropsMap} from "../../../app/type"
import {PageStartupBox} from "../../view/box/PageStartupBox"
import {KeyEnterInput} from "../../view/input/KeyEnterInput"

const SX: SxPropsMap = {
    alert: {width: "100%", padding: "0 20px", justifyContent: "center"},
    label: {color: "text.secondary", fontWeight: "bold", fontSize: "25px"},
    types: {display: "flex", gap: 1, width: "100%", flexDirection: "column"},
}

type Props = {
    supported: AuthType[],
    error?: string,
}

export function LoginBody(props: Props) {
    const {supported, error} = props
    const [auth, setAuth] = useState<AuthType>(supported[0] ?? AuthType.BASIC)
    const [username, setUsername] = useState("")
    const [password, setPass] = useState("")

    const login = useRouterLogin()
    const logoutRouter = useRouterLogout()

    return (
        <PageStartupBox header={"Authentication"} renderFooter={renderFooter()} position={"start"}>
            {renderButtons()}
            <Divider/>
            {error && <Alert sx={SX.alert} severity={"warning"} icon={false}>{error}</Alert>}
            <Divider/>
            {renderCreds()}
        </PageStartupBox>
    )

    function renderCreds() {
        if ([AuthType.OIDC].includes(auth)) return null

        return (
            <>
                <KeyEnterInput label={"username"} onChange={(e) => setUsername(e.target.value)}/>
                <KeyEnterInput
                    label={"password"}
                    hidden
                    onChange={(e) => setPass(e.target.value)}
                    onEnterPress={() => handleLogin(AuthType.BASIC)}
                />
            </>
        )
    }

    function renderFooter() {
        if (!supported.length) return <Button onClick={handleLogout}>Sign out</Button>

        return (
            <Button onClick={() => handleLogin(auth)} loading={login.isPending} fullWidth={true}>
                SIGN IN
            </Button>
        )
    }

    function renderButtons() {
        if (supported.length < 2) return null

        return (
            <ToggleButtonGroup value={auth} exclusive={true} fullWidth={true} size={"small"}>
                {supported.map(type => (
                    <ToggleButton value={type} key={type} onClick={() => setAuth(type)}>
                        {AuthType[type]}
                    </ToggleButton>
                ))}
            </ToggleButtonGroup>
        )
    }

    function handleLogout() {
        logoutRouter.mutate()
    }

    function handleLogin(type: AuthType) {
        const subject = type === AuthType.OIDC ? undefined : {username, password}
        login.mutate({type, subject})
    }
}
