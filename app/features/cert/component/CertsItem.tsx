import {Box, Tooltip} from "@mui/material"

import {DeleteIconButton} from "../../../shared/component/button/IconButtons"
import {StylePropsMap, SxPropsMap} from "../../../shared/helper/type"
import {FileUsageOptions} from "../../../shared/helper/utils"
import {Feature} from "../../feature"
import {ManageAccess} from "../../management/component/ManageAccess"
import {useRouterCertDelete} from "../api/hook"
import {Cert} from "../api/type"

const SX: SxPropsMap = {
    item: {
        display: "flex", alignItems: "center", padding: "5px 10px", margin: "5px 10px",
        borderRadius: "5px", gap: 1, border: 1, borderColor: "divider", height: "42px",
    },
    name: {flexBasis: "150px"},
    path: {flexBasis: "400px", fontSize: "13px", color: "text.disabled"},
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
            <Tooltip placement={"top-start"} title={cert.fileName}>
                <Box sx={SX.name} style={style.break}>{cert.fileName}</Box>
            </Tooltip>
            <Tooltip placement={"top-start"} title={cert.path}>
                <Box sx={SX.path} style={style.break}>{cert.path}</Box>
            </Tooltip>
            <ManageAccess feature={Feature.ManageCertDelete}>
                <DeleteIconButton loading={deleteCert.isPending} onClick={() => deleteCert.mutate(uuid)}/>
            </ManageAccess>
        </Box>
    )
}
