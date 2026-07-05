import {InfoOutlined} from "@mui/icons-material"
import {Alert, Box, Collapse, Tab, Tabs, ToggleButton, Tooltip} from "@mui/material"
import {ReactNode, useState} from "react"

import {NodeTabType} from "../../../../features/node/api/NodeType"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {useStoreAction} from "../../../../shared/provider/StoreProvider"
import {NODE_TABS} from "./NodeMainTabs"

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
                    <Tab value={NodeTabType.PLATFORM} label={NODE_TABS[NodeTabType.PLATFORM].label}/>
                    <Tab value={NodeTabType.CONTAINER} label={NODE_TABS[NodeTabType.CONTAINER].label}/>
                    <Tab value={NodeTabType.KEEPER} label={NODE_TABS[NodeTabType.KEEPER].label}/>
                    <Tab value={NodeTabType.DATABASE} label={NODE_TABS[NodeTabType.DATABASE].label}/>
                    <Tab value={NodeTabType.TOOLS} label={NODE_TABS[NodeTabType.TOOLS].label}/>
                </Tabs>
                <Box sx={[SX.buttons, {flexGrow: 1}]}>
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
