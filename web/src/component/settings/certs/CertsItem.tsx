import {Box, Tooltip} from "@mui/material";
import {useAppearance} from "../../../provider/AppearanceProvider";
import {DeleteIconButton} from "../../view/button/IconButtons";
import {useMutation} from "@tanstack/react-query";
import {certApi} from "../../../app/api";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {FileUsageOptions} from "../../../app/utils";
import {StylePropsMap, SxPropsMap} from "../../../type/common";
import {Cert} from "../../../type/cert";

const SX: SxPropsMap = {
    item: { display: "flex", alignItems: "center", padding: "5px 10px", margin: "5px 10px", borderRadius: "5px", gap: 2 },
    body: { flexGrow: 1, display: "flex", gap: 2 },
    name: { flexBasis: "150px" },
    path: { flexBasis: "400px", fontSize: "13px" },
}


const style: StylePropsMap = {
    break: {textOverflow: "ellipsis", whiteSpace: 'nowrap', overflow: "hidden"},
}

type Props = {
    cert: Cert,
    uuid: string,
}

export function CertsItem(props: Props) {
    const { cert, uuid } = props
    const { info } = useAppearance()

    const deleteOptions = useMutationOptions([["certs"]])
    const deleteCert = useMutation(certApi.delete, deleteOptions)

    const fileUsage = FileUsageOptions[cert.fileUsageType]

    return (
        <Box sx={{...SX.item, border: `1px solid ${info?.palette.divider}`}}>
            <Tooltip placement={"top"} title={fileUsage.label}>
                {fileUsage.icon}
            </Tooltip>
            <Tooltip placement={"top-start"} title={cert.fileName}>
                <Box sx={SX.name} style={style.break}>{cert.fileName}</Box>
            </Tooltip>
            <Tooltip placement={"top-start"} title={cert.path}>
                <Box sx={{...SX.path, color: info?.palette.text.disabled}} style={style.break}>{cert.path}</Box>
            </Tooltip>
            <DeleteIconButton loading={deleteCert.isLoading} onClick={() => deleteCert.mutate(uuid)} />
        </Box>
    )
}
