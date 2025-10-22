import {PageStartupBox} from "../../view/box/PageStartupBox";
import {Alert, Button} from "@mui/material";
import {KeyEnterInput} from "../../view/input/KeyEnterInput";
import {useState} from "react";
import {useRouterLogin, useRouterLogout} from "../../../api/auth/hook";
import {SxPropsMap} from "../../../app/type";
import {AuthType} from "../../../api/auth/type";

const SX: SxPropsMap = {
    alert: {width: "100%", padding: "0 20px", justifyContent: "center"},
}

type Props = {
    supported: AuthType[],
    error?: string,
}

export function LoginBody(props: Props) {
    const {supported, error} = props
    const [username, setUsername] = useState("")
    const [password, setPass] = useState("")

    const login = useRouterLogin()
    const logoutRouter = useRouterLogout()

    return (
        <PageStartupBox header={"Authentication"} renderFooter={renderButtons()}>
            {error && <Alert sx={SX.alert} severity={"warning"} icon={false}>{error}</Alert>}
            {renderCreds()}
        </PageStartupBox>
    )

    function renderCreds() {
        if (![AuthType.BASIC, AuthType.LDAP].some(r => supported.includes(r))) return null

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

    function renderButtons() {
        if (!supported.length) return <Button onClick={handleLogout}>Sign out</Button>

        return supported.map(type => (
            <Button key={type} onClick={() => handleLogin(type)} loading={login.isPending}>Sign in with {AuthType[type]}</Button>
        ))
    }

    function handleLogout() {
        logoutRouter.mutate()
    }

    function handleLogin(type: AuthType) {
        const subject = type === AuthType.OIDC ? undefined : {username, password}
        login.mutate({type, subject})
    }
}
