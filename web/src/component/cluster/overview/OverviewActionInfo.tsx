import {CertOptions, CredentialOptions, getDomain, InstanceColor} from "../../../app/utils";
import {PasswordType} from "../../../type/password";
import {CertType} from "../../../type/cert";
import {WarningAmberRounded} from "@mui/icons-material";
import {orange, purple} from "@mui/material/colors";
import {Box} from "@mui/material";
import {InfoIcons} from "../../view/box/InfoIcons";
import {InfoBox} from "../../view/box/InfoBox";
import {InfoTitle} from "../../view/box/InfoTitle";
import {ActiveCluster} from "../../../type/cluster";
import {SxPropsMap} from "../../../type/common";

const SX: SxPropsMap = {
    box: {display: "flex", alignItems: "center", gap: 1},
}

type Props = {
    cluster: ActiveCluster,
}

export function OverviewActionInfo(props: Props) {
    const {cluster, defaultInstance, warning, detection} = props.cluster

    const infoItems = [
        {...CredentialOptions[PasswordType.POSTGRES], active: !!cluster.credentials.postgresId},
        {...CredentialOptions[PasswordType.PATRONI], active: !!cluster.credentials.patroniId},
        {...CertOptions[CertType.CLIENT_CA], active: !!cluster.certs.clientCAId},
        {...CertOptions[CertType.CLIENT_CERT], active: !!cluster.certs.clientCertId},
        {...CertOptions[CertType.CLIENT_KEY], active: !!cluster.certs.clientKeyId}
    ]
    const warningItems = [
        {icon: <WarningAmberRounded/>, label: "Warning", active: warning, iconColor: orange[500]}
    ]
    const roleTooltip = [
        {label: "Detection", value: detection, bgColor: purple[400]},
        {label: "Instance", value: getDomain(defaultInstance.sidecar), bgColor: InstanceColor[defaultInstance.role]}
    ]

    return (
        <Box sx={SX.box}>
            <InfoIcons items={warningItems}/>
            <InfoIcons items={infoItems}/>
            <InfoBox tooltip={<InfoTitle items={roleTooltip}/>} withPadding>
                <Box sx={{color: InstanceColor[defaultInstance.role]}}>
                    {defaultInstance.role.toUpperCase()}
                </Box>
            </InfoBox>
        </Box>
    )
}
