import {PageStartupBox} from "../../view/box/PageStartupBox";
import {LoadingButton} from "@mui/lab";
import {KeyEnterInput} from "../../view/input/KeyEnterInput";
import {useState} from "react";
import {AuthConfig, AuthType} from "../../../type/common";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {GeneralApi} from "../../../app/api";
import {ConfigQuery} from "./query/ConfigQuery";
import {ConfigAuth} from "./auth/ConfigAuth";

export function ConfigBody() {
    const [company, setCompany] = useState("")
    const [query, setQuery] = useState(false)
    const [auth, setAuth] = useState<AuthConfig>({type: AuthType.NONE, body: undefined})

    const configOptions = useMutationOptions([["info"]])
    const config = useMutation({mutationFn: GeneralApi.setConfig, ...configOptions})

    return (
        <PageStartupBox header={"Configuration"} renderFooter={renderFooter()}>
            <KeyEnterInput label={"Company"} onChange={(e) => setCompany(e.target.value)}/>
            <ConfigQuery onChange={setQuery}/>
            <ConfigAuth onChange={setAuth} auth={auth}/>
        </PageStartupBox>
    )

    function renderFooter() {
        return (
            <LoadingButton variant={"contained"} loading={config.isPending} onClick={handleClick}>
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
