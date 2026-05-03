import {Divider, ToggleButton, ToggleButtonGroup} from "@mui/material"

import {CertType} from "../../../api/cert/type"
import {Options as ClusterOptions} from "../../../api/cluster/type"
import {Feature} from "../../../api/feature"
import {VaultType} from "../../../api/vault/type"
import {SxPropsMap} from "../../../app/type"
import {CertOptions, VaultOptions} from "../../../app/utils"
import {AccessBox} from "../access/Access"
import {OptionsCert} from "./OptionsCert"
import {OptionsTags} from "./OptionsTags"
import {OptionsVault} from "./OptionsVault"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    tls: {width: "50%", fontWeight: "bold"},
}

type Props = {
    cluster: ClusterOptions,
    onUpdate: (cluster: ClusterOptions) => void,
}

export function Options(props: Props) {
    const {onUpdate, cluster} = props
    const {vaults, tags, certs, tls} = cluster

    return (
        <AccessBox sx={SX.box} feature={Feature.ManageClusterUpdate}>
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
        </AccessBox>
    )

    function handleVaultUpdate(t: VaultType, s?: string) {
        onUpdate({...cluster, vaults: {...cluster.vaults, [VaultOptions[t].key]: s}})
    }

    function handleCertUpdate(t: CertType, s?: string) {
        onUpdate({...cluster, certs: {...cluster.certs, [CertOptions[t].key]: s}})
    }

    function handleTagsUpdate(tags: string[]) {
        onUpdate({...cluster, tags})
    }

    function handleTlsKeeperUpdate() {
        onUpdate({...cluster, tls: {...cluster.tls, keeper: !tls.keeper}})
    }

    function handleTlsDatabaseUpdate() {
        onUpdate({...cluster, tls: {...cluster.tls, database: !tls.database}})
    }
}
