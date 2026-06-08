import {InfoOutlined} from "@mui/icons-material"
import {Alert, Box, Collapse, Tab, Tabs, ToggleButton, Tooltip} from "@mui/material"
import {ReactNode, useState} from "react"

import {NodeTabType} from "../../../../features/node/type"
import {SxPropsMap} from "../../../../shared/helper/type"
import {useStoreAction} from "../../../../shared/provider/StoreProvider"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column"},
    title: {display: "flex", justifyContent: "space-between", alignItems: "center", columnGap: 3, flexWrap: "wrap", alignContent: "stretch"},
    toggle: {padding: "3px"},
    buttons: {display: "flex", justifyContent: "flex-end", alignItems: "center", gap: 1, height: "45px"},
}

type Props = {
    info: ReactNode,
    tab: NodeTabType,
    renderActions?: ReactNode,
}

export function NodeMainTitle(props: Props) {
    const {tab, info, renderActions} = props
    const [alert, setAlert] = useState(false)
    const {setNodeBody} = useStoreAction

    return (
        <Box sx={SX.box}>
            <Box sx={SX.title}>
                <Tabs value={tab} onChange={(_, e) => setNodeBody(e)}>
                    <Tab value={NodeTabType.MONITOR} label={"Monitor"}/>
                    <Tab value={NodeTabType.QUERY} label={"Queries"}/>
                </Tabs>
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
