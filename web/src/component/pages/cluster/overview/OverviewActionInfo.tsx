import {Box} from "@mui/material"

import {CertType} from "../../../../api/cert/type"
import {Cluster, Instance} from "../../../../api/cluster/type"
import {PasswordType} from "../../../../api/password/type"
import {SxPropsMap} from "../../../../app/type"
import {CertOptions, CredentialOptions, getDetectionItems} from "../../../../app/utils"
import {InfoBox} from "../../../view/box/InfoBox"
import {InfoBoxList} from "../../../view/box/InfoBoxList"
import {InfoColorBoxList} from "../../../view/box/InfoColorBoxList"

const SX: SxPropsMap = {
    box: {display: "flex", alignItems: "center", gap: 1},
}

type Props = {
    cluster: Cluster,
    detectBy?: Instance,
    mainInstance?: Instance,
}

export function OverviewActionInfo(props: Props) {
    const {mainInstance, cluster, detectBy} = props

    const infoItems = [
        {...CredentialOptions[PasswordType.POSTGRES], active: !!cluster.credentials.postgresId},
        {...CredentialOptions[PasswordType.PATRONI], active: !!cluster.credentials.patroniId},
        {...CertOptions[CertType.CLIENT_CA], active: !!cluster.certs.clientCAId},
        {...CertOptions[CertType.CLIENT_CERT], active: !!cluster.certs.clientCertId},
        {...CertOptions[CertType.CLIENT_KEY], active: !!cluster.certs.clientKeyId}
    ]

    const detectionItems = getDetectionItems(mainInstance, detectBy)
    const instance = detectionItems[1]

    return (
        <Box sx={SX.box}>
            <InfoBox tooltip={<InfoColorBoxList items={detectionItems} label={"Cluster Detection"}/>}>
                <Box sx={{color: instance.bgColor}}>
                    {instance.label.toUpperCase()}
                </Box>
            </InfoBox>
            <InfoBoxList items={infoItems} label={"Configured Cluster Options"}/>
        </Box>
    )
}
