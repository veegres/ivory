import {Box, Paper, Switch, Tab, Tabs} from "@mui/material"
import {ChangeEvent, useState} from "react"

import {AuthType} from "../../../../api/auth/type"
import {BasicConfig, LdapConfig, OidcConfig} from "../../../../api/config/type"
import {SxPropsMap} from "../../../../app/type"
import {ConfigBox} from "../../../view/box/ConfigBox"
import {KeyEnterInput} from "../../../view/input/KeyEnterInput";
import {ConfigAuthBasic} from "./ConfigAuthBasic"
import {ConfigAuthLdap} from "./ConfigAuthLdap"
import {ConfigAuthOidc} from "./ConfigAuthOidc"

const SX: SxPropsMap = {
    body: {display: "flex", flexDirection: "column", gap: 2},
    tab: {padding: "10px", minHeight: "40px"},
    tabs: {minHeight: "40px"},
}

type Props = {
    ldapConfig: LdapConfig,
    onLdapChange: (config: LdapConfig) => void,
    basicConfig: BasicConfig,
    onBasicChange: (config: BasicConfig) => void,
    oidcConfig: OidcConfig,
    onOidcChange: (config: OidcConfig) => void,
    onDefaultChange: () => void,
    onAdminsChange: (v: string) => void,
}

export function ConfigAuth(props: Props) {
    const {ldapConfig, oidcConfig, basicConfig, onOidcChange, onBasicChange, onLdapChange, onDefaultChange, onAdminsChange} = props
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
            <Paper variant={"outlined"}>
                <Tabs sx={SX.tabs} variant={"fullWidth"} value={authTypeOpen} onChange={(_, v) => setAuthTypeOpen(v)}>
                    <Tab sx={SX.tab} value={AuthType.BASIC} label={AuthType[AuthType.BASIC]}/>
                    <Tab sx={SX.tab} value={AuthType.LDAP} label={AuthType[AuthType.LDAP]}/>
                    <Tab sx={SX.tab} value={AuthType.OIDC} label={AuthType[AuthType.OIDC]}/>
                </Tabs>
            </Paper>
        )
    }

    function renderBody() {
        return (
            <Box sx={SX.body}>
                <KeyEnterInput
                    label={"Superusers"}
                    required={true}
                    helperText={"at least one superuser should be provided, use ', ' for several"}
                    onChange={(e) => onAdminsChange(e.target.value)}
                />
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