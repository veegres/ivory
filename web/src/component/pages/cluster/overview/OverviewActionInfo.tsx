import {Box} from "@mui/material"

import {CertType} from "../../../../api/cert/type"
import {useRouterClusterOverview} from "../../../../api/cluster/hook"
import {Cluster} from "../../../../api/cluster/type"
import {VaultType} from "../../../../api/vault/type"
import {SxPropsMap} from "../../../../app/type"
import {CertOptions, getDetectionItems,VaultOptions} from "../../../../app/utils"
import {InfoBox} from "../../../view/box/InfoBox"
import {InfoBoxList} from "../../../view/box/InfoBoxList"
import {InfoColorBoxList} from "../../../view/box/InfoColorBoxList"

const SX: SxPropsMap = {
    box: {display: "flex", alignItems: "center", gap: 1},
}

type Props = {
    cluster: Cluster,
    manualKeeper?: string,
}

export function OverviewActionInfo(props: Props) {
    const {cluster, manualKeeper} = props

    const infoItems = [
        {...VaultOptions[VaultType.DATABASE_PASSWORD], active: !!cluster.vaults.databaseId},
        {...VaultOptions[VaultType.KEEPER_PASSWORD], active: !!cluster.vaults.keeperId},
        {...CertOptions[CertType.CLIENT_CA], active: !!cluster.certs.clientCAId},
        {...CertOptions[CertType.CLIENT_CERT], active: !!cluster.certs.clientCertId},
        {...CertOptions[CertType.CLIENT_KEY], active: !!cluster.certs.clientKeyId}
    ]

    const overview = useRouterClusterOverview(cluster.name, false)
    const {nodes, detectedDomain} = overview.data ?? {}
    const detectionItems = getDetectionItems(nodes, detectedDomain, manualKeeper)
    const node = detectionItems[1]

    return (
        <Box sx={SX.box}>
            <InfoBox tooltip={<InfoColorBoxList items={detectionItems} label={"Cluster Detection"}/>}>
                <Box sx={{color: node.bgColor}}>
                    {node.label.toUpperCase()}
                </Box>
            </InfoBox>
            <InfoBoxList items={infoItems} label={"Configured Cluster Options"}/>
        </Box>
    )
}
