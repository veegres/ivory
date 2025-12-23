import {FileDownloadOutlined} from "@mui/icons-material"
import {Box, Button} from "@mui/material"

import {SxPropsMap} from "../../../app/type"

const SX: SxPropsMap = {
    wrap: {padding: "5px", height: "100%"},
    button: {
        minHeight: "120px", minWidth: "350px", width: "100%", height: "100%", display: "flex",
        alignItems: "center", flexDirection: "column", justifyContent: "space-evenly",
        padding: "10px 20px", border: "2px dashed", borderRadius: "15px", gap: 1,
    },
    label: {textTransform: "uppercase", fontSize: "12px"},
}

type Props = {
    onClick: () => void,
}

export function DownloadButton(props: Props) {
    return (
        <Box sx={SX.wrap}>
            <Button sx={SX.button} color={"inherit"} onClick={props.onClick}>
                <FileDownloadOutlined fontSize={"large"}/>
                <Box sx={SX.label}>Click here to download the file</Box>
            </Button>
        </Box>
    )
}
