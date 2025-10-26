import {InfoOutlined} from "@mui/icons-material"
import {Alert, Box, Collapse, ToggleButton, Tooltip} from "@mui/material"
import {ReactNode, useState} from "react"

import {Database} from "../../../../api/query/type"
import {SxPropsMap} from "../../../../app/type"
import {DatabaseBox} from "../../../view/box/DatabaseBox"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column"},
    title: {display: "flex", justifyContent: "space-between", alignItems: "center", columnGap: 3, flexWrap: "wrap", alignContent: "stretch"},
    label: {fontWeight: "bold", fontSize: "30px", width: "125px"},
    toggle: {padding: "3px"},
    buttons: {display: "flex", justifyContent: "flex-end", alignItems: "center", gap: 1, height: "45px"},
    info: {display: "flex", justifyContent: "space-between", alignItems: "center", gap: 3},
}

type Props = {
    label: string,
    info: ReactNode,
    db: Database,
    renderActions?: ReactNode,
}

export function InstanceMainTitle(props: Props) {
    const {label, info, renderActions, db} = props
    const [alert, setAlert] = useState(false)

    return (
        <Box sx={SX.box}>
            <Box sx={SX.title}>
                <Box flexGrow={1} sx={SX.info}>
                    <Box sx={SX.label}>{label}</Box>
                    <DatabaseBox db={db}/>
                </Box>
                <Box flexGrow={1} sx={SX.buttons}>
                    {renderActions}
                    <ToggleButton sx={SX.toggle} value={"info"} size={"small"} selected={alert} onClick={() => setAlert(!alert)}>
                        <Tooltip title={"Description"} placement={"top"}>
                            <InfoOutlined/>
                        </Tooltip>
                    </ToggleButton>
                </Box>
            </Box>
            <Collapse in={alert}>
                <Alert severity={"info"} onClose={() => setAlert(false)}>{info}</Alert>
            </Collapse>
        </Box>
    )
}
