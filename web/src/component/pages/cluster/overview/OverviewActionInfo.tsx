import {CertOptions, CredentialOptions, getDomain, InstanceColor} from "../../../../app/utils";
import {PasswordType} from "../../../../api/password/type";
import {CertType} from "../../../../api/cert/type";
import {purple} from "@mui/material/colors";
import {Box} from "@mui/material";
import {InfoBoxList} from "../../../view/box/InfoBoxList";
import {InfoBox} from "../../../view/box/InfoBox";
import {InfoColorBoxList} from "../../../view/box/InfoColorBoxList";
import {ActiveCluster} from "../../../../api/cluster/type";
import {SxPropsMap} from "../../../../api/management/type";

const SX: SxPropsMap = {
    box: {display: "flex", alignItems: "center", gap: 1},
}

type Props = {
    cluster: ActiveCluster,
}

export function OverviewActionInfo(props: Props) {
    const {cluster, defaultInstance, detection} = props.cluster

    const infoItems = [
        {...CredentialOptions[PasswordType.POSTGRES], active: !!cluster.credentials.postgresId},
        {...CredentialOptions[PasswordType.PATRONI], active: !!cluster.credentials.patroniId},
        {...CertOptions[CertType.CLIENT_CA], active: !!cluster.certs.clientCAId},
        {...CertOptions[CertType.CLIENT_CERT], active: !!cluster.certs.clientCertId},
        {...CertOptions[CertType.CLIENT_KEY], active: !!cluster.certs.clientKeyId}
    ]

    const detectionItems = [
        {title: "Detection", label: detection, bgColor: purple[400]},
        {title: "Main Instance", label: getDomain(defaultInstance.sidecar), bgColor: InstanceColor[defaultInstance.role]}
    ]

    return (
        <Box sx={SX.box}>
            <InfoBox tooltip={<InfoColorBoxList items={detectionItems} label={"Cluster Detection"}/>}>
                <Box sx={{color: InstanceColor[defaultInstance.role]}}>
                    {defaultInstance.sidecar.host.toUpperCase()}
                </Box>
            </InfoBox>
            <InfoBoxList items={infoItems} label={"Configured Cluster Options"}/>
        </Box>
    )
}
