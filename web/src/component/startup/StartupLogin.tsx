import {StartupBlock} from "./StartupBlock";
import {Alert, Button} from "@mui/material";
import {SecretInput} from "../shared/secret/SecretInput";
import {useMutation, useQueryClient} from "@tanstack/react-query";
import {generalApi} from "../../app/api";
import {useMutationOptions} from "../../hook/QueryCustom";
import {useState} from "react";
import {SxPropsMap} from "../../type/common";
import {LoadingButton} from "@mui/lab";
import {useAuth} from "../../provider/AuthProvider";

const SX: SxPropsMap = {
    alert: {width: "100%", padding: "0 20px", justifyContent: "center"},
}

type Props = {
    type: "none" | "basic",
    error: string,
}

export function StartupLogin(props: Props) {
    const {type, error} = props
    const {setToken, logout} = useAuth()
    const [username, setUsername] = useState("")
    const [password, setPass] = useState("")

    const queryClient = useQueryClient();
    const loginOption = useMutationOptions([], handleSuccess)
    const login = useMutation(generalApi.login, loginOption)

    return (
        <StartupBlock header={"Authentication"} renderFooter={renderLogin()}>
            {error && <Alert sx={SX.alert} severity={"warning"} icon={false}>{error}</Alert>}
            {renderBody()}
        </StartupBlock>
    )

    function renderBody() {
        switch (type) {
            case "none": return null
            case "basic": return (
                <>
                    <SecretInput label={"username"} onChange={(e) => setUsername(e.target.value)}/>
                    <SecretInput
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
            case "none": return <Button onClick={handleLogout}>Sign out</Button>
            case "basic": return <LoadingButton onClick={handleLogin} loading={login.isLoading}>Sign in</LoadingButton>
        }
    }

    function handleLogout() {
        logout()
        queryClient.refetchQueries(["info"]).then()
    }

    function handleLogin() {
        login.mutate({username, password})
    }

    function handleSuccess(data: any) {
        if (data && data.token) setToken(data.token)
        queryClient.refetchQueries(["info"]).then()
    }
}
