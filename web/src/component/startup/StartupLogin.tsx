import {StartupBlock} from "./StartupBlock";
import {Alert} from "@mui/material";
import {SecretInput} from "../shared/secret/SecretInput";
import {useMutation} from "@tanstack/react-query";
import {infoApi} from "../../app/api";
import {useMutationOptions} from "../../hook/QueryCustom";
import {useState} from "react";
import {SxPropsMap} from "../../type/common";
import {LoadingButton} from "@mui/lab";
import {useAuth} from "../../provider/AuthProvider";

const SX: SxPropsMap = {
    alert: {width: "100%", padding: "0 20px", justifyContent: "center"},
}

type Props = {
    error: string,
}

export function StartupLogin(props: Props) {
    const {error} = props
    const {setToken} = useAuth()
    const [username, setUsername] = useState("")
    const [password, setPass] = useState("")

    const loginOption = useMutationOptions([], handleSuccess, [["info"]])
    const login = useMutation(infoApi.login, loginOption)

    return (
        <StartupBlock header={"Sign in"} renderFooter={renderButton()}>
            {error && <Alert sx={SX.alert} severity={"warning"} icon={false}>{error}</Alert>}
            <SecretInput label={"username"} onChange={(e) => setUsername(e.target.value)}/>
            <SecretInput
                label={"password"}
                hidden
                onChange={(e) => setPass(e.target.value)}
                onEnterPress={handleLogin}
            />
        </StartupBlock>
    )

    function renderButton() {
        return (
            <LoadingButton onClick={handleLogin} loading={login.isLoading}>Sign in</LoadingButton>
        )
    }

    function handleLogin() {
        login.mutate({username, password})
    }

    function handleSuccess(data: any) {
        if (data && data.token) setToken(data.token)
    }
}
