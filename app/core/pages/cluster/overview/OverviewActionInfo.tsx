import {Box} from "@mui/material"
import {useMemo} from "react"

import {CertType} from "../../../../features/cert/api/type"
import {Cluster} from "../../../../features/cluster/api/type"
import {VaultType} from "../../../../features/vault/api/type"
import {InfoBoxList} from "../../../../shared/component/box/InfoBoxList"
import {SxPropsMap} from "../../../../shared/helper/type"
import {CertOptions, VaultOptions} from "../../../../shared/helper/utils"

const SX: SxPropsMap = {
    box: {display: "flex", alignItems: "center", gap: 1},
}

type Props = {
    cluster: Cluster,
}

export function OverviewActionInfo(props: Props) {
    const {cluster} = props
    const clusterItems = useMemo(handleClusterItemsMemo, [cluster])

    return (
        <Box sx={SX.box}>
            <InfoBoxList items={clusterItems} label={"Configured Cluster Options"}/>
        </Box>
    )

    function handleClusterItemsMemo() {
        return [
            {...VaultOptions[VaultType.SSH_KEY], active: !!cluster.vaults.sshKeyId},
            {...VaultOptions[VaultType.KEEPER_PASSWORD], active: !!cluster.vaults.keeperId},
            {...VaultOptions[VaultType.DATABASE_PASSWORD], active: !!cluster.vaults.databaseId},
            {...CertOptions[CertType.CLIENT_CA], active: !!cluster.certs.clientCAId},
            {...CertOptions[CertType.CLIENT_CERT], active: !!cluster.certs.clientCertId},
            {...CertOptions[CertType.CLIENT_KEY], active: !!cluster.certs.clientKeyId}
        ]
    }
}
