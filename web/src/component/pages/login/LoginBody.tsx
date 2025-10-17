import {PageStartupBox} from "../../view/box/PageStartupBox";
import {Alert, Button} from "@mui/material";
import {KeyEnterInput} from "../../view/input/KeyEnterInput";
import {useState} from "react";
import {useRouterLogin, useRouterLogout} from "../../../api/auth/hook";
import {AuthType} from "../../../api/config/type";
import {SxPropsMap} from "../../../app/type";

const SX: SxPropsMap = {
    alert: {width: "100%", padding: "0 20px", justifyContent: "center"},
}

type Props = {
    type: AuthType,
    error: string,
}

export function LoginBody(props: Props) {
    const {type, error} = props
    const [username, setUsername] = useState("")
    const [password, setPass] = useState("")

    const login = useRouterLogin(type)
    const logoutRouter = useRouterLogout()

    return (
        <PageStartupBox header={"Authentication"} renderFooter={renderLogin()}>
            {error && <Alert sx={SX.alert} severity={"warning"} icon={false}>{error}</Alert>}
            {renderBody()}
        </PageStartupBox>
    )

    function renderBody() {
        switch (type) {
            case AuthType.NONE: return null
            case AuthType.BASIC: return (
                <>
                    <KeyEnterInput label={"username"} onChange={(e) => setUsername(e.target.value)}/>
                    <KeyEnterInput
                        label={"password"}
                        hidden
                        onChange={(e) => setPass(e.target.value)}
                        onEnterPress={handleLogin}
                    />
                </>
            )
            case AuthType.LDAP: return (
                <>
                    <KeyEnterInput label={"username"} onChange={(e) => setUsername(e.target.value)}/>
                    <KeyEnterInput
                        label={"password"}
                        hidden
                        onChange={(e) => setPass(e.target.value)}
                        onEnterPress={handleLogin}
                    />
                </>
            )
            case AuthType.OIDC: return null
        }
    }

    function renderLogin() {
        switch (type) {
            case AuthType.NONE: return <Button onClick={handleLogout}>Sign out</Button>
            case AuthType.BASIC: return <Button onClick={handleLogin} loading={login.isPending}>Sign in</Button>
            case AuthType.LDAP: return <Button onClick={handleLogin} loading={login.isPending}>Sign in</Button>
            case AuthType.OIDC: return <Button onClick={handleOidcLogin}>Sign in with SSO</Button>
        }
    }

    function handleLogout() {
        logoutRouter.mutate()
    }

    function handleOidcLogin() {
        login.mutate(undefined)
    }

    function handleLogin() {
        login.mutate({username, password})
    }
}
