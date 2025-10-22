import {ChangeEvent, useState} from "react";
import {Box, Switch, Tab, Tabs} from "@mui/material";
import {ConfigBox} from "../../../view/box/ConfigBox";
import {SxPropsMap} from "../../../../app/type";
import {AuthType} from "../../../../api/auth/type";
import {BasicConfig, LdapConfig, OidcConfig} from "../../../../api/config/type";
import {ConfigAuthBasic} from "./ConfigAuthBasic";
import {ConfigAuthOidc} from "./ConfigAuthOidc";
import {ConfigAuthLdap} from "./ConfigAuthLdap";

const SX: SxPropsMap = {
    body: {display: "flex", flexDirection: "column", gap: 2},
}

type Props = {
    ldapConfig: LdapConfig,
    onLdapChange: (config: LdapConfig) => void,
    basicConfig: BasicConfig,
    onBasicChange: (config: BasicConfig) => void,
    oidcConfig: OidcConfig,
    onOidcChange: (config: OidcConfig) => void,
    onDefaultChange: () => void,
}

export function ConfigAuth(props: Props) {
    const {ldapConfig, oidcConfig, basicConfig, onOidcChange, onBasicChange, onLdapChange, onDefaultChange} = props
    const [authEnabled, setAuthEnabled] = useState(false)
    const [authTypeOpen, setAuthTypeOpen] = useState(AuthType.BASIC)

    return (
        <ConfigBox
            label={"Authentication"}
            renderAction={<Switch size={"small"} onChange={handleAuthOpen}/>}
            renderBody={renderBody()}
            showBody={authEnabled}
            description={<>
                Without authentication, anyone with access to <b>Ivory</b> can fully control the system,
                including using sensitive data and performing critical database actions like
                switchover or reinitialization. <b>Ivory</b> provides support
                for <b>Basic</b>, <b>OIDC</b> and <b>LDAP</b> authentication methods.
            </>}
            recommendation={<>
                <b>No Auth:</b> Suitable for <b>local use</b> on your personal computer.<br/>
                <b>Basic Auth:</b> Recommended for <b>teams</b> or <b>small groups</b>.<br/>
                <b>OIDC / LDAP:</b> Ideal for <b>companies</b> using an OIDC or LDAP provider.<br/>
            </>}
            severity={authEnabled ? "info" : "warning"}
        />
    )

    function renderAuthSwitch() {
        return (
            <Tabs variant={"fullWidth"} value={authTypeOpen} onChange={(_, v) => setAuthTypeOpen(v)}>
                <Tab value={AuthType.BASIC} label={AuthType[AuthType.BASIC]}/>
                <Tab value={AuthType.LDAP} label={AuthType[AuthType.LDAP]}/>
                <Tab value={AuthType.OIDC} label={AuthType[AuthType.OIDC]}/>
            </Tabs>
        )
    }

    function renderBody() {
        return (
            <Box sx={SX.body}>
                {renderAuthSwitch()}
                {renderAuthBody()}
            </Box>
        )
    }

    function renderAuthBody() {
        switch (authTypeOpen) {
            case AuthType.BASIC: return <ConfigAuthBasic config={basicConfig} onChange={onBasicChange}/>
            case AuthType.OIDC: return <ConfigAuthOidc config={oidcConfig} onChange={onOidcChange}/>
            case AuthType.LDAP: return <ConfigAuthLdap config={ldapConfig} onChange={onLdapChange}/>
            default: return null
        }
    }

    function handleAuthOpen(e: ChangeEvent<HTMLInputElement>) {
        const type = e.target.checked
        onDefaultChange()
        setAuthEnabled(type)
    }
}