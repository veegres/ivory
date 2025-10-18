import {PageStartupBox} from "../../view/box/PageStartupBox";
import {KeyEnterInput} from "../../view/input/KeyEnterInput";
import {useState} from "react";
import {ConfigQuery} from "./query/ConfigQuery";
import {ConfigAuth} from "./auth/ConfigAuth";
import {Button} from "@mui/material";
import {useRouterConfigSet} from "../../../api/config/hook";
import {AuthConfig, AuthType} from "../../../api/config/type";
import {useRouterInfo} from "../../../api/management/hook";

export function ConfigBody() {
    const [company, setCompany] = useState("")
    const [query, setQuery] = useState(false)
    const [auth, setAuth] = useState<AuthConfig>({type: AuthType.NONE})
    const config = useRouterConfigSet()
    const info = useRouterInfo()

    return (
        <PageStartupBox header={"Configuration"} renderFooter={renderFooter()}>
            <KeyEnterInput label={"Company"} onChange={(e) => setCompany(e.target.value)}/>
            <ConfigQuery onChange={setQuery}/>
            <ConfigAuth onChange={setAuth} auth={auth} path={info.data?.path ?? ""}/>
        </PageStartupBox>
    )

    function renderFooter() {
        return (
            <Button variant={"contained"} loading={config.isPending} onClick={handleClick}>
                Set
            </Button>
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
