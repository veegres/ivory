import {StartupBlock} from "./StartupBlock";
import {LoadingButton} from "@mui/lab";
import {KeyEnterInput} from "../view/input/KeyEnterInput";
import {StartupConfigAuth} from "./StartupConfigAuth";
import {StartupConfigQuery} from "./StartupConfigQuery";
import {useState} from "react";
import {AuthConfig, AuthType} from "../../type/common";
import {useMutationOptions} from "../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {generalApi} from "../../app/api";

export function StartupInitialConfig() {
    const [company, setCompany] = useState("")
    const [query, setQuery] = useState(false)
    const [auth, setAuth] = useState<AuthConfig>({type: AuthType.NONE, body: undefined})

    const configOptions = useMutationOptions([["info"]])
    const config = useMutation(generalApi.setConfig, configOptions)

    return (
        <StartupBlock header={"Configuration"} renderFooter={renderFooter()}>
            <KeyEnterInput label={"Company"} onChange={(e) => setCompany(e.target.value)}/>
            <StartupConfigQuery onChange={setQuery}/>
            <StartupConfigAuth onChange={setAuth} auth={auth}/>
        </StartupBlock>
    )

    function renderFooter() {
        return (
            <LoadingButton variant={"contained"} loading={config.isLoading} onClick={handleClick}>
                Set
            </LoadingButton>
        )
    }

    function handleClick() {
        config.mutate({
            company,
            auth,
            availability: {manualQuery: query},
        })
    }
}
