import {Box, Divider, ToggleButton, ToggleButtonGroup} from "@mui/material"
import {memo} from "react"

import {SxPropsMap} from "../../../shared/helper/HelperType"
import {CertOptions, VaultOptions} from "../../../shared/helper/HelperUtils"
import {CertType} from "../../cert/api/CertType"
import {Feature} from "../../Feature"
import {ManageAccessBox} from "../../management/component/ManageAccess"
import {VaultType} from "../../vault/api/VaultType"
import {Options, Plugins} from "../api/ClusterType"
import {ClusterOptionsCert} from "./ClusterOptionsCert"
import {ClusterOptionsPlugins} from "./ClusterOptionsPlugins"
import {ClusterOptionsTags} from "./ClusterOptionsTags"
import {ClusterOptionsVault} from "./ClusterOptionsVault"

// NOTE: each field is a wrap item with a min width, so the widget lays fields
// out in a dynamic grid when its container is wide (e.g. the overview page)
// but collapses to a single, wider column when the container is narrow
// (e.g. the deploy/detect dialogs)
const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "row", flexWrap: "wrap", gap: 1.5, alignItems: "flex-start", marginTop: "10px"},
    divider: {flexBasis: "100%"},
    field: {flex: "1 1 var(--size-field)", minWidth: "var(--size-field)"},
    tls: {width: "50%", fontWeight: "bold"},
}

type Props = {
    options: Options,
    onUpdate: (options: Options) => void,
    disablePlugins?: boolean,
}

export const ClusterOptions = memo(function ClusterOptions(props: Props) {
    const {onUpdate, options, disablePlugins = false} = props
    const {vaults, tags, certs, tls, plugins} = options

    return (
        <ManageAccessBox sx={SX.box} feature={Feature.ManageClusterUpdate}>
            <ClusterOptionsPlugins plugins={plugins} onUpdate={handlePluginsUpdate} disabled={disablePlugins}/>
            <Divider sx={SX.divider}/>
            <Box sx={SX.field}><ClusterOptionsVault type={VaultType.DATABASE_PASSWORD} selected={vaults.databaseId} onUpdate={handleVaultUpdate}/></Box>
            <Box sx={SX.field}><ClusterOptionsVault type={VaultType.KEEPER_PASSWORD} selected={vaults.keeperId} onUpdate={handleVaultUpdate}/></Box>
            <Box sx={SX.field}><ClusterOptionsVault type={VaultType.SSH_KEY} selected={vaults.sshKeyId} onUpdate={handleVaultUpdate}/></Box>
            <Divider sx={SX.divider}/>
            <ToggleButtonGroup sx={SX.field} size={"small"} fullWidth>
                <ToggleButton onClick={handleTlsKeeperUpdate} selected={tls.keeper} value={"keeper"}>Keeper</ToggleButton>
                <ToggleButton sx={SX.tls} disabled={true} value={"tls"}>TLS</ToggleButton>
                <ToggleButton onClick={handleTlsDatabaseUpdate} selected={tls.database} value={"database"}>Database</ToggleButton>
            </ToggleButtonGroup>
            <Divider sx={SX.divider}/>
            <Box sx={SX.field}><ClusterOptionsCert type={CertType.CLIENT_CA} selected={certs.clientCAId} onUpdate={handleCertUpdate}/></Box>
            <Box sx={SX.field}><ClusterOptionsCert type={CertType.CLIENT_CERT} selected={certs.clientCertId} onUpdate={handleCertUpdate}/></Box>
            <Box sx={SX.field}><ClusterOptionsCert type={CertType.CLIENT_KEY} selected={certs.clientKeyId} onUpdate={handleCertUpdate}/></Box>
            <Divider sx={SX.divider}/>
            <Box sx={SX.field}><ClusterOptionsTags selected={tags} onUpdate={handleTagsUpdate}/></Box>
        </ManageAccessBox>
    )

    function handlePluginsUpdate(plugins: Plugins) {
        onUpdate({...options, plugins})
    }

    function handleVaultUpdate(t: VaultType, s?: string) {
        onUpdate({...options, vaults: {...options.vaults, [VaultOptions[t].key]: s}})
    }

    function handleCertUpdate(t: CertType, s?: string) {
        onUpdate({...options, certs: {...options.certs, [CertOptions[t].key]: s}})
    }

    function handleTagsUpdate(tags: string[]) {
        onUpdate({...options, tags})
    }

    function handleTlsKeeperUpdate() {
        onUpdate({...options, tls: {...options.tls, keeper: !tls.keeper}})
    }

    function handleTlsDatabaseUpdate() {
        onUpdate({...options, tls: {...options.tls, database: !tls.database}})
    }
})
