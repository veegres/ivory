import {Box, ToggleButton, ToggleButtonGroup, Tooltip} from "@mui/material";
import {useState} from "react";
import {CertsTabPath} from "./CertsTabPath";
import {CertsTabUpload} from "./CertsTabUpload";
import {FileUsageOptions} from "../../../../app/utils";
import {CertType, FileUsageType} from "../../../../api/cert/type";
import {SxPropsMap} from "../../../../app/type";

const SX: SxPropsMap = {
    box: {display: "flex", height: "130px"},
    group: {padding: "5px"},
    button: {flexGrow: 1, borderRadius: "15px"},
    tab: {flexGrow: 1}
}

type Props = {
    type: CertType,
}

export function CertsNew(props: Props) {
    const [tab, setTab] = useState(FileUsageType.UPLOAD)

    const upload = FileUsageOptions[FileUsageType.UPLOAD]
    const path = FileUsageOptions[FileUsageType.PATH]

    return (
        <Box sx={SX.box}>
            <ToggleButtonGroup sx={SX.group} orientation={"vertical"} size={"small"} exclusive value={tab}>
                <ToggleButton sx={SX.button} onClick={() => setTab(FileUsageType.UPLOAD)} value={FileUsageType.UPLOAD}>
                    <Tooltip placement={"right"} title={upload.label}>{upload.icon}</Tooltip>
                </ToggleButton>
                <ToggleButton sx={SX.button} onClick={() => setTab(FileUsageType.PATH)} value={FileUsageType.PATH}>
                    <Tooltip placement={"right"} title={path.label}>{path.icon}</Tooltip>
                </ToggleButton>
            </ToggleButtonGroup>
            <Box sx={SX.tab}>{renderTab()}</Box>
        </Box>
    )

    function renderTab() {
        switch (tab) {
            case FileUsageType.UPLOAD: return <CertsTabUpload {...props} />
            case FileUsageType.PATH: return <CertsTabPath {...props} />
        }
    }
}
