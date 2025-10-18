import {ChangeEvent} from "react";
import {AlertColor, Box, FormControlLabel, Radio, RadioGroup} from "@mui/material";
import {ConfigBox} from "../../../view/box/ConfigBox";
import {KeyEnterInput} from "../../../view/input/KeyEnterInput";
import {AuthConfig, AuthType} from "../../../../api/config/type";
import {SxPropsMap} from "../../../../app/type";

const SX: SxPropsMap = {
    radio: {display: "flex"},
    basic: {display: "flex", gap: 1},
    oidc: {display: "flex", flexDirection: "column", gap: 1},
    ldap: {display: "flex", flexDirection: "column", gap: 1},
}

const Description = {
    [AuthType.NONE]: <>
        No authentications means that anyone will have direct access to <b>Ivory</b> and
        to all the information provided like certs, password, etc. Anyone who has a link to
        Ivory will be able to perform any action with database like switchover, reinit, etc.
    </>,
    [AuthType.BASIC]: <>
        Basic authentication helps to protect <b>Ivory</b> data from direct access
        with general username and password. It provides sign in form to get access
        to <b>Ivory</b>. It will generate unique token for each sign in. Each token
        has expiration time. You need to provide general username and password
        that will be used below.
    </>,
    [AuthType.OIDC]: <>
        OIDC authentication helps to protect <b>Ivory</b> data from direct access
        with OIDC provider. It provides sign in form to get access to <b>Ivory</b>.
        It will generate a unique token for each sign in. Each token has expiration time.
    </>,
    [AuthType.LDAP]: <>
        LDAP authentication helps to protect <b>Ivory</b> data from direct access
        with LDAP provider. It provides sign in form to get access to <b>Ivory</b>.
        It will generate unique token for each sign in. Each token has expiration time.
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
                <KeyEnterInput label={"Issuer URL"} value={auth.oidc?.issuerUrl ?? ""} onChange={handleConfigChange("oidc", "issuerUrl")}/>
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
                <KeyEnterInput label={"URL"} value={auth.ldap?.url ?? ""} onChange={handleConfigChange("ldap", "url")}/>
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