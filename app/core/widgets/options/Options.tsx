import {Divider, ToggleButton, ToggleButtonGroup} from "@mui/material"
import {memo} from "react"

import {CertType} from "../../../features/cert/api/type"
import {Options as ClusterOptions, Plugins} from "../../../features/cluster/api/type"
import {Feature} from "../../../features/feature"
import {ManageAccessBox} from "../../../features/management/component/ManageAccess"
import {VaultType} from "../../../features/vault/api/type"
import {SxPropsMap} from "../../../shared/helper/type"
import {CertOptions, VaultOptions} from "../../../shared/helper/utils"
import {OptionsCert} from "./OptionsCert"
import {OptionsPlugins} from "./OptionsPlugins"
import {OptionsTags} from "./OptionsTags"
import {OptionsVault} from "./OptionsVault"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    tls: {width: "50%", fontWeight: "bold"},
}

type Props = {
    options: ClusterOptions,
    onUpdate: (options: ClusterOptions) => void,
}

export const Options = memo(function Options(props: Props) {
    const {onUpdate, options} = props
    const {vaults, tags, certs, tls, plugins} = options

    return (
        <ManageAccessBox sx={SX.box} feature={Feature.ManageClusterUpdate}>
            <OptionsPlugins plugins={plugins} onUpdate={handlePluginsUpdate}/>
            <Divider variant={"middle"}/>
            <OptionsVault type={VaultType.DATABASE_PASSWORD} selected={vaults.databaseId} onUpdate={handleVaultUpdate}/>
            <OptionsVault type={VaultType.KEEPER_PASSWORD} selected={vaults.keeperId} onUpdate={handleVaultUpdate}/>
            <OptionsVault type={VaultType.SSH_KEY} selected={vaults.sshKeyId} onUpdate={handleVaultUpdate}/>
            <Divider variant={"middle"}/>
            <ToggleButtonGroup size={"small"} fullWidth>
                <ToggleButton onClick={handleTlsKeeperUpdate} selected={tls.keeper} value={"keeper"}>Keeper</ToggleButton>
                <ToggleButton sx={SX.tls} disabled={true} value={"tls"}>TLS</ToggleButton>
                <ToggleButton onClick={handleTlsDatabaseUpdate} selected={tls.database} value={"database"}>Database</ToggleButton>
            </ToggleButtonGroup>
            <Divider variant={"middle"}/>
            <OptionsCert type={CertType.CLIENT_CA} selected={certs.clientCAId} onUpdate={handleCertUpdate}/>
            <OptionsCert type={CertType.CLIENT_CERT} selected={certs.clientCertId} onUpdate={handleCertUpdate}/>
            <OptionsCert type={CertType.CLIENT_KEY} selected={certs.clientKeyId} onUpdate={handleCertUpdate}/>
            <Divider variant={"middle"}/>
            <OptionsTags selected={tags} onUpdate={handleTagsUpdate}/>
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
