import {Box, Tooltip} from "@mui/material"

import {Code} from "../../../shared/component/box/Code"
import {DeleteIconButton} from "../../../shared/component/button/IconButtons"
import {StylePropsMap} from "../../../shared/helper/HelperType"
import {FileUsageOptions, getShortUuid} from "../../../shared/helper/HelperUtils"
import {Feature} from "../../Feature"
import {ManageAccess} from "../../management/component/ManageAccess"
import {useRouterCertDelete} from "../api/CertHook"
import {Cert} from "../api/CertType"

const SX = {
    item: {
        display: "flex", alignItems: "center", padding: "5px 10px", margin: "5px 10px",
        borderRadius: "5px", gap: 1, border: 1, borderColor: "divider",
    },
    info: {display: "flex", flexDirection: "column", flexGrow: 1, minWidth: 0},
    name: {lineHeight: 1.2},
    id: {fontSize: "12px", flexShrink: 0},
    path: {fontSize: "13px", color: "text.disabled", lineHeight: 1.2},
}

const style: StylePropsMap = {
    break: {textOverflow: "ellipsis", whiteSpace: "nowrap", overflow: "hidden"},
}

type Props = {
    cert: Cert,
    uuid: string,
}

export function CertsItem(props: Props) {
    const {cert, uuid} = props
    const deleteCert = useRouterCertDelete(cert.type)
    const fileUsage = FileUsageOptions[cert.fileUsageType]

    return (
        <Box sx={SX.item}>
            <Tooltip placement={"top"} title={fileUsage.label}>
                {fileUsage.icon}
            </Tooltip>
            <Box sx={SX.info}>
                <Tooltip placement={"top-start"} title={cert.fileName}>
                    <Box sx={SX.name} style={style.break}>{cert.fileName}</Box>
                </Tooltip>
                <Tooltip placement={"top-start"} title={cert.path}>
                    <Box sx={SX.path} style={style.break}>{cert.path}</Box>
                </Tooltip>
            </Box>
            <Tooltip placement={"top"} title={uuid}>
                <Code sx={SX.id}>{getShortUuid(uuid)}</Code>
            </Tooltip>
            <ManageAccess feature={Feature.ManageCertDelete}>
                <DeleteIconButton loading={deleteCert.isPending} onClick={() => deleteCert.mutate(uuid)}/>
            </ManageAccess>
        </Box>
    )
}
