import {Alert, Box, Button} from "@mui/material"
import {useState} from "react"

import {useRouterConfigSet} from "../../../features/config/api/ConfigHook"
import {LdapConfig, OidcConfig} from "../../../features/config/api/ConfigType"
import {UserAuthType, UserSetupRequest} from "../../../features/user/api/UserType"
import {PageStartupBox} from "../../../shared/component/box/PageStartupBox"
import {KeyEnterInput} from "../../../shared/component/input/KeyEnterInput"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {ConfigAuth} from "./auth/ConfigAuth"

const SX: SxPropsMap = {
    alert: {width: "100%", padding: "0 20px", justifyContent: "center"},
    center: {display: "flex", justifyContent: "center"},
}

// NOTE: the first user is offered every way of signing in, basic included -
// unlike everybody after them, whose password arrives as a link, this one is
// typed here, and it is the only way into a fresh Ivory
const defaultUser: UserSetupRequest = {
    username: "",
    password: "",
    authTypes: [UserAuthType.BASIC, UserAuthType.LDAP, UserAuthType.OIDC],
}
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
    const [basicEnabled, setBasicEnabled] = useState(false)
    const [user, setUser] = useState<UserSetupRequest>(defaultUser)
    const [ldap, setLdap] = useState<LdapConfig>(defaultLdapConfig)
    const [oidc, setOidc] = useState<OidcConfig>(defaultOidcConfig)
    const config = useRouterConfigSet()
    const isConfigBroken = configured && error

    return (
        <PageStartupBox header={"Configuration"} renderFooter={renderFooter()} position={"start"}>
            {isConfigBroken && renderError()}
            <KeyEnterInput label={"Company"} onChange={(e) => setCompany(e.target.value)}/>
            <ConfigAuth
                basicEnabled={basicEnabled} ldapConfig={ldap} oidcConfig={oidc}
                user={user} onUserChange={setUser}
                onDefaultChange={handleDefaultChange}
                onOidcChange={setOidc} onBasicChange={setBasicEnabled} onLdapChange={setLdap}
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
        setBasicEnabled(false)
        setUser(defaultUser)
        setLdap(defaultLdapConfig)
        setOidc(defaultOidcConfig)
    }

    function isUserProvided() {
        return user.username !== "" || user.password !== ""
    }

    function handleClick() {
        config.mutate({
            secret: isConfigBroken ? secret : undefined,
            // NOTE: an untouched user is left out entirely, so re-running setup
            // where a superuser already exists does not have to name one again
            user: isUserProvided() ? user : undefined,
            appConfig: {
                company,
                auth: {
                    basic: basicEnabled ? {} : undefined,
                    ldap: JSON.stringify(ldap) === JSON.stringify(defaultLdapConfig) ? undefined : ldap,
                    oidc: JSON.stringify(oidc) === JSON.stringify(defaultOidcConfig) ? undefined : oidc,
                },
            }
        })
    }
}
