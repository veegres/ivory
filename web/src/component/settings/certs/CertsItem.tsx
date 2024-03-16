import {Box, Tooltip} from "@mui/material";
import {DeleteIconButton} from "../../view/button/IconButtons";
import {useMutation} from "@tanstack/react-query";
import {CertApi} from "../../../app/api";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {FileUsageOptions} from "../../../app/utils";
import {StylePropsMap, SxPropsMap} from "../../../type/common";
import {Cert} from "../../../type/cert";

const SX: SxPropsMap = {
    item: {display: "flex", alignItems: "center", padding: "5px 10px", margin: "5px 10px", borderRadius: "5px", gap: 2, border: 1, borderColor: "divider"},
    body: {flexGrow: 1, display: "flex", gap: 2},
    name: {flexBasis: "150px"},
    path: {flexBasis: "400px", fontSize: "13px", color: "text.disabled"},
}

const style: StylePropsMap = {
    break: {textOverflow: "ellipsis", whiteSpace: 'nowrap', overflow: "hidden"},
}

type Props = {
    cert: Cert,
    uuid: string,
}

export function CertsItem(props: Props) {
    const {cert, uuid} = props

    const deleteOptions = useMutationOptions([["certs"]])
    const deleteCert = useMutation({mutationFn: CertApi.delete, ...deleteOptions})

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
            <DeleteIconButton loading={deleteCert.isPending} onClick={() => deleteCert.mutate(uuid)}/>
        </Box>
    )
}
