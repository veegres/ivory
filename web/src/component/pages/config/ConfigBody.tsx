import {Alert, Box, Button} from "@mui/material"
import {useState} from "react"

import {useRouterConfigSet} from "../../../api/config/hook"
import {BasicConfig, LdapConfig, OidcConfig} from "../../../api/config/type"
import {SxPropsMap} from "../../../app/type"
import {PageStartupBox} from "../../view/box/PageStartupBox"
import {KeyEnterInput} from "../../view/input/KeyEnterInput"
import {ConfigAuth} from "./auth/ConfigAuth"

const SX: SxPropsMap = {
    alert: {width: "100%", padding: "0 20px", justifyContent: "center"},
    center: {display: "flex", justifyContent: "center"},
}

const defaultBasicConfig: BasicConfig = {password: "", username: ""}
const defaultLdapConfig: LdapConfig = {bindPass: "", bindDN: "", baseDN: "", filter: "(uid=%s)", url: ""}
// NOTE: we set base url when configure env variables in the backend
const defaultOidcConfig: OidcConfig = {clientSecret: "", clientId: "", issuerUrl: "", redirectUrl: `${document.baseURI}api/oidc/callback`}

type Props = {
    configured: boolean,
    error?: string,
}

export function ConfigBody(props: Props) {
    const {configured, error} = props
    const [company, setCompany] = useState("")
    const [secret, setSecret] = useState("")
    const [admins, setAdmins] = useState("")
    const [basic, setBasic] = useState<BasicConfig>(defaultBasicConfig)
    const [ldap, setLdap] = useState<LdapConfig>(defaultLdapConfig)
    const [oidc, setOidc] = useState<OidcConfig>(defaultOidcConfig)
    const config = useRouterConfigSet()
    const isConfigBroken = configured && error

    return (
        <PageStartupBox header={"Configuration"} renderFooter={renderFooter()} position={"start"}>
            {isConfigBroken && renderError()}
            <KeyEnterInput label={"Company"} onChange={(e) => setCompany(e.target.value)}/>
            <ConfigAuth
                basicConfig={basic} ldapConfig={ldap} oidcConfig={oidc}
                onDefaultChange={handleDefaultChange} onAdminsChange={setAdmins}
                onOidcChange={setOidc} onBasicChange={setBasic} onLdapChange={setLdap}
            />
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
            <Button variant={"contained"} loading={config.isPending} onClick={handleClick}>Set</Button>
        </>)
    }

    function handleDefaultChange() {
        setBasic(defaultBasicConfig)
        setLdap(defaultLdapConfig)
        setOidc(defaultOidcConfig)
    }

    function handleClick() {
        config.mutate({
            secret: isConfigBroken ? secret : undefined,
            appConfig: {
                company,
                auth: {
                    superusers: admins.split(", "),
                    basic: JSON.stringify(basic) === JSON.stringify(defaultBasicConfig) ? undefined : basic,
                    ldap: JSON.stringify(ldap) === JSON.stringify(defaultLdapConfig) ? undefined : ldap,
                    oidc: JSON.stringify(oidc) === JSON.stringify(defaultOidcConfig) ? undefined : oidc,
                },
            }
        })
    }
}
