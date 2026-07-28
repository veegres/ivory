import {Edit, Preview} from "@mui/icons-material"
import {Box, TextField, ToggleButton, ToggleButtonGroup, Tooltip} from "@mui/material"

import {SxPropsMap} from "../../helper/HelperType"
import {Code} from "../box/Code"

const SX: SxPropsMap = {
    between: {display: "flex", justifyContent: "space-between", alignItems: "center", gap: 1},
    note: {
        display: "flex", justifyContent: "center", alignItems: "center",
        color: "text.disabled", fontSize: 12, flexWrap: "wrap", gap: 0.5,
    },
}

type Props = {
    imageUri: string,
    onImageUriChange: (value: string) => void,
    preview: boolean,
    onPreviewChange: (value: boolean) => void,
    placeholderKeys: string[],
}

// DeployImageHeader renders the image URI field, the preview/edit toggle and
// the interpolation placeholder hints shared by every keeper deploy dialog
// (ClusterDeploy and ContainerKeeperDeploy).
export function DeployImageHeader(props: Props) {
    const {imageUri, onImageUriChange, preview, onPreviewChange, placeholderKeys} = props
    return (
        <>
            <Box sx={SX.between}>
                <TextField
                    fullWidth={true}
                    size={"small"}
                    label={"Image"}
                    value={imageUri}
                    onChange={v => onImageUriChange(v.target.value)}
                />
                <ToggleButtonGroup value={preview} exclusive={true} size={"small"} onChange={(_, v) => onPreviewChange(v)}>
                    <Tooltip title={"Preview"} placement={"top"}>
                        <ToggleButton value={true}><Preview/></ToggleButton>
                    </Tooltip>
                    <Tooltip title={"Edit"} placement={"top"}>
                        <ToggleButton value={false}><Edit/></ToggleButton>
                    </Tooltip>
                </ToggleButtonGroup>
            </Box>
            <Box sx={SX.note}>
                <Box>Use interpolated options to automatically populate values</Box>
                <Box sx={SX.note}>
                    {placeholderKeys.map(k => (
                        <Code key={k} sx={{fontSize: "11px"}}>{k}</Code>
                    ))}
                </Box>
            </Box>
        </>
    )
}
