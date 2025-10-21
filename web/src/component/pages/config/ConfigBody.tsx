import {PageStartupBox} from "../../view/box/PageStartupBox";
import {KeyEnterInput} from "../../view/input/KeyEnterInput";
import {useState} from "react";
import {ConfigQuery} from "./query/ConfigQuery";
import {ConfigAuth} from "./auth/ConfigAuth";
import {Alert, Box, Button} from "@mui/material";
import {useRouterConfigSet} from "../../../api/config/hook";
import {AuthType} from "../../../api/auth/type";
import {AuthConfig} from "../../../api/config/type";
import {SxPropsMap} from "../../../app/type";

const SX: SxPropsMap = {
    alert: {width: "100%", padding: "0 20px", justifyContent: "center"},
    center: {display: "flex", justifyContent: "center"},
}

type Props = {
    configured: boolean,
    error?: string,
}

export function ConfigBody(props: Props) {
    const {configured, error} = props
    const [company, setCompany] = useState("")
    const [query, setQuery] = useState(false)
    const [secret, setSecret] = useState("")
    const [auth, setAuth] = useState<AuthConfig>({type: AuthType.NONE})
    const config = useRouterConfigSet()
    const isConfigBroken = configured && error

    return (
        <PageStartupBox header={"Configuration"} renderFooter={renderFooter()}>
            {isConfigBroken && renderError()}
            <KeyEnterInput label={"Company"} onChange={(e) => setCompany(e.target.value)}/>
            <ConfigQuery onChange={setQuery}/>
            <ConfigAuth onChange={setAuth} auth={auth}/>
        </PageStartupBox>
    )

    function renderError() {
        return (<>
            <Alert sx={SX.alert} severity={"info"} icon={false}>
                The configuration could not be initialized due to an unexpected issue.
                Please set it up again from scratch. Since a previous configuration exists,
                you’ll need to provide the <b>secret word</b> to continue working with sensitive data.
            </Alert>
            <Alert sx={SX.alert} severity={"error"} icon={false}>
                <Box sx={SX.center}><b>CONFIGURATION ISSUE</b></Box>
                <Box sx={SX.center}>{error}</Box>
            </Alert>
        </>)
    }

    function renderFooter() {
        return (<>
            {isConfigBroken && (
                <KeyEnterInput
                    label={"Secret word"}
                    onChange={(e) => setSecret(e.target.value)}
                    onEnterPress={handleClick}
                    hidden
                />
            )}
            <Button
                variant={"contained"}
                loading={config.isPending}
                onClick={handleClick}
            >
                Set
            </Button>
        </>)
    }

    function handleClick() {
        config.mutate({
            secret: isConfigBroken ? secret : undefined,
            appConfig: {
                company,
                auth,
                availability: {manualQuery: query},
            }
        })
    }
}
