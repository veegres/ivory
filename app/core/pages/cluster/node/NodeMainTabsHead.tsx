import {InfoOutlined} from "@mui/icons-material"
import {Alert, Box, Collapse, Tab, Tabs, ToggleButton, Tooltip} from "@mui/material"
import {ReactNode, useEffect, useState} from "react"

import {Feature} from "../../../../features/Feature"
import {useHasAccess} from "../../../../features/management/component/ManageAccess"
import {NodeTabType} from "../../../../features/node/api/NodeType"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {useStoreAction} from "../../../../shared/provider/StoreProvider"
import {NODE_TABS} from "./NodeMainTabs"

// NOTE: the buttons row grows into the space left by the tabs and keeps
// space-between inside, so the tab actions sit on the left of that space and
// the info toggle on the right — both inline and when the row wraps below
const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    title: {display: "flex", justifyContent: "space-between", alignItems: "center", columnGap: 3, flexWrap: "wrap", alignContent: "stretch"},
    // NOTE: the tabs grow factor dwarfs the buttons row's, so the tabs take
    // all the slack inline while the wrapped buttons row still fills its own
    // line and keeps its space-between
    tabs: {
        minWidth: 0, maxWidth: "100%", flexGrow: 100,
        "& .MuiTabs-scrollButtons.Mui-disabled": {opacity: 0.3},
    },
    toggle: {padding: "3px"},
    buttons: {display: "flex", justifyContent: "space-between", alignItems: "center", gap: 1, minHeight: "45px", flexGrow: 1},
    actions: {display: "flex", alignItems: "center", gap: 1, flexWrap: "wrap"},
}

type Props = {
    info: ReactNode,
    tab: NodeTabType,
    renderActions?: ReactNode,
}

export function NodeMainTabsHead(props: Props) {
    const {tab, info, renderActions} = props
    const [alert, setAlert] = useState(false)
    const {setNodeBody} = useStoreAction
    // NOTE: only "unsupported" hides the tab - a permission the user lacks is
    // still worth showing, so the tab body can say so
    const system = useHasAccess(Feature.ViewNodeSystem) !== "unsupported"

    useEffect(handleEffectSystemTab, [system, tab, setNodeBody])

    return (
        <Box sx={SX.box}>
            <Box sx={SX.title}>
                <Tabs sx={SX.tabs} value={tab} onChange={(_, e) => setNodeBody(e)} variant={"scrollable"} scrollButtons={"auto"} allowScrollButtonsMobile>
                    {system && <Tab value={NodeTabType.SYSTEM} label={NODE_TABS[NodeTabType.SYSTEM].label}/>}
                    <Tab value={NodeTabType.CONTAINER} label={NODE_TABS[NodeTabType.CONTAINER].label}/>
                    <Tab value={NodeTabType.KEEPER} label={NODE_TABS[NodeTabType.KEEPER].label}/>
                    <Tab value={NodeTabType.DATABASE} label={NODE_TABS[NodeTabType.DATABASE].label}/>
                    <Tab value={NodeTabType.TOOLS} label={NODE_TABS[NodeTabType.TOOLS].label}/>
                </Tabs>
                <Box sx={SX.buttons}>
                    <Box sx={SX.actions}>{renderActions}</Box>
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

    function handleEffectSystemTab() {
        if (!system && tab === NodeTabType.SYSTEM) setNodeBody(NodeTabType.CONTAINER)
    }
}
