import {Box} from "@mui/material"
import {useMemo} from "react"

import {CertType} from "../../../../api/cert/type"
import {Cluster, Node} from "../../../../api/cluster/type"
import {VaultType} from "../../../../api/vault/type"
import {SxPropsMap} from "../../../../app/type"
import {CertOptions, getDetectionItems,VaultOptions} from "../../../../app/utils"
import {useStore} from "../../../../provider/StoreProvider"
import {InfoBox} from "../../../view/box/InfoBox"
import {InfoBoxList} from "../../../view/box/InfoBoxList"
import {InfoColorBoxList} from "../../../view/box/InfoColorBoxList"

const SX: SxPropsMap = {
    box: {display: "flex", alignItems: "center", gap: 1},
}

type Props = {
    cluster: Cluster,
    mainNode: [string?, Node?],
}

export function OverviewActionInfo(props: Props) {
    const {cluster, mainNode} = props
    const manualKeeper = useStore(s => s.manualKeeper)

    const clusterItems = useMemo(handleClusterItemsMemo, [cluster])

    const detectionItems = getDetectionItems(mainNode, !!manualKeeper)
    const node = detectionItems[1]

    return (
        <Box sx={SX.box}>
            <InfoBox tooltip={<InfoColorBoxList items={detectionItems} label={"Cluster Detection"}/>}>
                <Box sx={{color: node.bgColor}}>{node.label.toUpperCase()}</Box>
            </InfoBox>
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
