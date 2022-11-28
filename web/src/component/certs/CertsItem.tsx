import {TextSnippetOutlined} from "@mui/icons-material";
import {Box} from "@mui/material";
import React from "react";
import {Cert} from "../../app/types";
import {useTheme} from "../../provider/ThemeProvider";
import {DeleteIconButton} from "../view/IconButtons";
import {useMutation} from "@tanstack/react-query";
import {certApi} from "../../app/api";
import {useMutationOptions} from "../../app/hooks";

const SX = {
    item: { display: "flex", alignItems: "center", padding: "5px 30px", margin: "5px 10px", borderRadius: "5px" },
    body: { flexGrow: 1, display: "flex" },
    name: { textAlign: "center" },
}

type Props = {
    cert: Cert
}

export function CertsItem(props: Props) {
    const { cert } = props
    const { info } = useTheme()

    const deleteOptions = useMutationOptions(["certs"])
    const deleteCert = useMutation(certApi.delete, deleteOptions)

    return (
        <Box sx={{...SX.item, border: `1px solid ${info?.palette.divider}`}} gap={1}>
            <Box sx={SX.body} gap={1}>
                <TextSnippetOutlined />
                <Box sx={SX.name}>{cert.fileName}</Box>
            </Box>
            <DeleteIconButton loading={deleteCert.isLoading} onClick={() => deleteCert.mutate(cert.fileId)} />
        </Box>
    )
}
