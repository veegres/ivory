import {Box, Tooltip} from "@mui/material"

import {useRouterCertDelete} from "../../../../api/cert/hook"
import {Cert} from "../../../../api/cert/type"
import {Permission} from "../../../../api/permission/type"
import {StylePropsMap, SxPropsMap} from "../../../../app/type"
import {FileUsageOptions} from "../../../../app/utils"
import {DeleteIconButton} from "../../../view/button/IconButtons"
import {Access} from "../../access/Access"

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
            <Access permission={Permission.ManageCertDelete}>
                <DeleteIconButton loading={deleteCert.isPending} onClick={() => deleteCert.mutate(uuid)}/>
            </Access>
        </Box>
    )
}
