import {PageStartupBox} from "../../view/box/PageStartupBox";
import {Alert, Button} from "@mui/material";
import {KeyEnterInput} from "../../view/input/KeyEnterInput";
import {useState} from "react";
import {AuthType, SxPropsMap} from "../../../type/general";
import {LoadingButton} from "@mui/lab";
import {useAuth} from "../../../provider/AuthProvider";
import {useRouterLogin} from "../../../router/auth";

const SX: SxPropsMap = {
    alert: {width: "100%", padding: "0 20px", justifyContent: "center"},
}

type Props = {
    type: AuthType,
    error: string,
}

export function LoginBody(props: Props) {
    const {type, error} = props
    const {setToken, logout} = useAuth()
    const [username, setUsername] = useState("")
    const [password, setPass] = useState("")

    const login = useRouterLogin(handleSuccess)

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
        }
    }

    function renderLogin() {
        switch (type) {
            case AuthType.NONE: return <Button onClick={logout}>Sign out</Button>
            case AuthType.BASIC: return <LoadingButton onClick={handleLogin} loading={login.isPending}>Sign in</LoadingButton>
        }
    }

    function handleLogin() {
        login.mutate({username, password})
    }

    function handleSuccess(data: any) {
        if (data && data.token) setToken(data.token)
    }
}
