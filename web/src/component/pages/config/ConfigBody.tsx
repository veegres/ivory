import {PageStartupBox} from "../../view/box/PageStartupBox";
import {KeyEnterInput} from "../../view/input/KeyEnterInput";
import {useState} from "react";
import {ConfigQuery} from "./query/ConfigQuery";
import {ConfigAuth} from "./auth/ConfigAuth";
import {Alert, Button} from "@mui/material";
import {useRouterConfigSet} from "../../../api/config/hook";
import {AuthType} from "../../../api/auth/type";
import {AuthConfig} from "../../../api/config/type";
import {SxPropsMap} from "../../../app/type";

const SX: SxPropsMap = {
    alert: {width: "100%", padding: "0 20px", justifyContent: "center"},
}

type Props = {
    error: string,
}

export function ConfigBody(props: Props) {
    const {error} = props
    const [company, setCompany] = useState("")
    const [query, setQuery] = useState(false)
    const [auth, setAuth] = useState<AuthConfig>({type: AuthType.NONE})
    const config = useRouterConfigSet()

    return (
        <PageStartupBox header={"Configuration"} renderFooter={renderFooter()}>
            {error && <Alert sx={SX.alert} severity={"warning"} icon={false}>{error}</Alert>}
            <KeyEnterInput label={"Company"} onChange={(e) => setCompany(e.target.value)}/>
            <ConfigQuery onChange={setQuery}/>
            <ConfigAuth onChange={setAuth} auth={auth}/>
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
