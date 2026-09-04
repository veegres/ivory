import {Box, Divider, Paper, Switch, Tab, Tabs} from "@mui/material"
import {ChangeEvent, useState} from "react"

import {AuthType} from "../../../../features/auth/api/AuthType"
import {LdapConfig, OidcConfig} from "../../../../features/config/api/ConfigType"
import {UserSetupRequest} from "../../../../features/user/api/UserType"
import {UserCreationForm} from "../../../../features/user/component/UserCreationForm"
import {AlertCentered} from "../../../../shared/component/box/AlertCentered"
import {ConfigBox} from "../../../../shared/component/box/ConfigBox"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {ConfigAuthBasic} from "./ConfigAuthBasic"
import {ConfigAuthLdap} from "./ConfigAuthLdap"
import {ConfigAuthOidc} from "./ConfigAuthOidc"

const SX: SxPropsMap = {
    body: {display: "flex", flexDirection: "column", gap: 2},
    gap: {display: "flex", flexDirection: "column", gap: 1},
    divider: {marginTop: 2},
}

type Props = {
    ldapConfig: LdapConfig,
    onLdapChange: (config: LdapConfig) => void,
    basicEnabled: boolean,
    onBasicChange: (enabled: boolean) => void,
    oidcConfig: OidcConfig,
    onOidcChange: (config: OidcConfig) => void,
    user: UserSetupRequest,
    onUserChange: (user: UserSetupRequest) => void,
    onDefaultChange: () => void,
}

export function ConfigAuth(props: Props) {
    const {ldapConfig, oidcConfig, basicEnabled, onOidcChange, onBasicChange, onLdapChange} = props
    const {user, onUserChange, onDefaultChange} = props
    const [authEnabled, setAuthEnabled] = useState(false)
    const [authTypeOpen, setAuthTypeOpen] = useState(AuthType.BASIC)

    return (
        <ConfigBox
            label={"Authentication"}
            renderAction={<Switch size={"small"} onChange={handleAuthOpen}/>}
            renderBody={renderBody()}
            showBody={authEnabled}
            description={<>
                Without authentication, anyone with access to <b>Ivory</b> can fully control it,
                including sensitive data, performing critical database or vm actions. <b>Ivory</b> provides
                support for <b>Basic</b>, <b>OIDC</b> and <b>LDAP</b> authentication methods.
            </>}
            recommendation={<>
                <b>No Auth:</b> Suitable for <b>local use</b> on your personal computer.<br/>
                <b>Basic Auth:</b> Recommended for <b>teams</b> or <b>small groups</b>.<br/>
                <b>OIDC / LDAP:</b> Ideal for <b>companies</b> using OIDC or LDAP provider.<br/>
            </>}
        />
    )

    function renderAuthSwitch() {
        return (
            <Paper variant={"outlined"}>
                <Tabs variant={"fullWidth"} value={authTypeOpen} onChange={(_, v) => setAuthTypeOpen(v)}>
                    <Tab value={AuthType.BASIC} label={AuthType[AuthType.BASIC]}/>
                    <Tab value={AuthType.LDAP} label={AuthType[AuthType.LDAP]}/>
                    <Tab value={AuthType.OIDC} label={AuthType[AuthType.OIDC]}/>
                </Tabs>
            </Paper>
        )
    }

    function renderBody() {
        return (
            <Box sx={SX.body}>
                <Box sx={SX.gap}>
                    <Divider>Configure Superuser</Divider>
                    {renderSuperuser()}
                </Box>
                <Box sx={SX.gap}>
                    <Divider>Configure Providers</Divider>
                    {renderAuthSwitch()}
                    {renderAuthBody()}
                </Box>
            </Box>
        )
    }

    function renderSuperuser() {
        return (
            <Box sx={SX.gap}>
                <AlertCentered text={renderSuperuserDescription()}/>
                <UserCreationForm setup value={user} onChange={onUserChange}/>
            </Box>
        )
    }

    function renderSuperuserDescription() {
        return (<>
            At least one <b>superuser</b> is required when authentication is enabled.
            Provide initial superuser credentials below. You can add more users after setup.
        </>)
    }

    function renderAuthBody() {
        switch (authTypeOpen) {
            case AuthType.BASIC: return <ConfigAuthBasic enabled={basicEnabled} onChange={onBasicChange}/>
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
