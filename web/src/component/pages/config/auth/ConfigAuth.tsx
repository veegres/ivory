import {ChangeEvent} from "react";
import {AlertColor, Box, FormControlLabel, Radio, RadioGroup} from "@mui/material";
import {ConfigBox} from "../../../view/box/ConfigBox";
import {KeyEnterInput} from "../../../view/input/KeyEnterInput";
import {SxPropsMap} from "../../../../app/type";
import {AuthType} from "../../../../api/auth/type";
import {AuthConfig} from "../../../../api/config/type";

const SX: SxPropsMap = {
    radio: {display: "flex"},
    basic: {display: "flex", gap: 1},
    oidc: {display: "flex", flexDirection: "column", gap: 1},
    ldap: {display: "flex", flexDirection: "column", gap: 1},
}

const Description = {
    [AuthType.NONE]: <>
        Without authentication, anyone who can access <b>Ivory</b> will have full
        control over the system. This includes using sensitive data (certificates, passwords, etc.)
        and executing dangerous database operations like switchover or reinitialization.
    </>,
    [AuthType.BASIC]: <>
        Basic authentication helps protect <b>Ivory</b> from unauthorized access using a
        general username and password. It provides a sign-in form that requires these credentials
        before granting access to <b>Ivory</b>. You’ll need to specify the username and password below.
    </>,
    [AuthType.OIDC]: <>
        OIDC authentication protects <b>Ivory</b> data by integrating with an external OpenID Connect (OIDC) provider.
        It’s recommended to create a dedicated application for <b>Ivory</b> and restrict access to specific users only.
        When configuring OIDC, ensure that <b>Ivory</b>’s callback URL is allowed on the provider side.
        Also, make sure the Issuer URL exactly matches the value specified in the <i>/.well-known/openid-configuration</i> endpoint.
    </>,
    [AuthType.LDAP]: <>
        LDAP authentication protects <b>Ivory</b> by integrating with an LDAP provider.
        It provides a sign-in form that requires valid LDAP credentials to access <b>Ivory</b>.
        You can specify a filter below to restrict access to certain users.
    </>
}

const Recommendation = {
    [AuthType.NONE]: <>
        This type of authentication is recommended for <b>local usage</b>, when
        you use <b>Ivory</b> in your personal laptop / computer.
    </>,
    [AuthType.BASIC]: <>
        This type of authentication is recommended for <b>teams</b> or <b>small
        group of people</b>.
    </>,
    [AuthType.OIDC]: <>
        This type of authentication is recommended for <b>companies</b> that use
        OIDC provider.
    </>,
    [AuthType.LDAP]: <>
        This type of authentication is recommended for <b>companies</b> that use
        LDAP provider.
    </>
}

const Severity: { [key in AuthType]: AlertColor } = {
    [AuthType.NONE]: "warning",
    [AuthType.BASIC]: "info",
    [AuthType.OIDC]: "info",
    [AuthType.LDAP]: "info",
}

type Props = {
    auth: AuthConfig,
    onChange: (auth: AuthConfig) => void,
}

export function ConfigAuth(props: Props) {
    const {onChange, auth} = props

    return (
        <ConfigBox
            label={"Authentication"}
            renderAction={renderAction()}
            renderBody={renderBody()}
            showBody={auth.type !== AuthType.NONE}
            description={Description[auth.type]}
            recommendation={Recommendation[auth.type]}
            severity={Severity[auth.type]}
        />
    )

    function renderAction() {
        return (
            <RadioGroup sx={SX.radio} row value={auth.type} onChange={handleAuthTypeChange}>
                <FormControlLabel
                    value={AuthType.NONE}
                    control={<Radio size={"small"}/>}
                    label={AuthType[AuthType.NONE].toLowerCase()}
                />
                <FormControlLabel
                    value={AuthType.BASIC}
                    control={<Radio size={"small"}/>}
                    label={AuthType[AuthType.BASIC].toLowerCase()}
                />
                <FormControlLabel
                    value={AuthType.LDAP}
                    control={<Radio size={"small"}/>}
                    label={AuthType[AuthType.LDAP].toLowerCase()}
                />
                <FormControlLabel
                    value={AuthType.OIDC}
                    control={<Radio size={"small"}/>}
                    label={AuthType[AuthType.OIDC].toLowerCase()}
                />
            </RadioGroup>
        )
    }

    function renderBody() {
        switch (auth.type) {
            case AuthType.BASIC: return renderBasic()
            case AuthType.OIDC: return renderOidc()
            case AuthType.LDAP: return renderLdap()
            default: return null
        }
    }

    function renderBasic() {
        return (
            <Box sx={SX.basic}>
                <KeyEnterInput label={"Username"} value={auth.basic?.username ?? ""} onChange={handleConfigChange("basic", "username")}/>
                <KeyEnterInput label={"Password"} value={auth.basic?.password ?? ""} onChange={handleConfigChange("basic", "password")} hidden />
            </Box>
        )
    }

    function renderOidc() {
        return (
            <Box sx={SX.oidc}>
                <KeyEnterInput
                    label={"Issuer URL"}
                    value={auth.oidc?.issuerUrl ?? ""}
                    helperText={"Specify the complete HTTP url, including the protocol prefix http:// or https://"}
                    onChange={handleConfigChange("oidc", "issuerUrl")}
                />
                <KeyEnterInput label={"Client ID"} value={auth.oidc?.clientId ?? ""} onChange={handleConfigChange("oidc", "clientId")}/>
                <KeyEnterInput label={"Client secret"} value={auth.oidc?.clientSecret ?? ""} onChange={handleConfigChange("oidc", "clientSecret")} hidden/>
                <KeyEnterInput
                    label={"Redirect URL"}
                    value={auth.oidc?.redirectUrl ?? ""}
                    helperText={"This URL must be allowed in your provider configuration"}
                    required={false}
                    disabled={true}
                    onChange={handleConfigChange("oidc", "redirectUrl")}
                />
            </Box>
        )
    }

    function renderLdap() {
        return (
            <Box sx={SX.ldap}>
                <KeyEnterInput
                    label={"URL"}
                    helperText={"Specify the complete LDAP url, including the protocol prefix ldap:// or ldaps://"}
                    value={auth.ldap?.url ?? ""}
                    onChange={handleConfigChange("ldap", "url")}
                />
                <KeyEnterInput label={"Bind DN"} value={auth.ldap?.bindDN ?? ""} onChange={handleConfigChange("ldap", "bindDN")}/>
                <KeyEnterInput label={"Bind password"} value={auth.ldap?.bindPass ?? ""} onChange={handleConfigChange("ldap", "bindPass")} hidden/>
                <KeyEnterInput label={"Base DN"} value={auth.ldap?.baseDN ?? ""} onChange={handleConfigChange("ldap", "baseDN")}/>
                <KeyEnterInput label={"Filter"} value={auth.ldap?.filter ?? ""} required={false} onChange={handleConfigChange("ldap", "filter")}/>
            </Box>
        )
    }

    function handleConfigChange(authType: "basic" | "ldap" | "oidc", key: string) {
        return (e: ChangeEvent<HTMLInputElement>) => {
            onChange({...auth, [authType]: {...auth[authType], [key]: e.target.value}})
        }
    }

    function handleAuthTypeChange(e: ChangeEvent<HTMLInputElement>) {
        const type = parseInt(e.target.value) as AuthType
        if (type === AuthType.OIDC) onChange({type, oidc: {issuerUrl: "", clientId: "", clientSecret: "", redirectUrl: getRedirectUrl()}})
        else onChange({type})
    }

    function getRedirectUrl() {
        // NOTE: we set base url when configure env variables in backend
        return `${document.baseURI}api/oidc/callback`
    }
}