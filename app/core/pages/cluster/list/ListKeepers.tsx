import {ArrowDropDown, ArrowDropUp} from "@mui/icons-material"
import {Box, Button, ClickAwayListener, MenuItem, MenuList, Paper, Popper, Tooltip} from "@mui/material"
import {useRef, useState} from "react"

import {KeeperPlugin, ReleaseStage} from "../../../../features/node/api/NodeType"
import {EnumOptions, SxPropsMap} from "../../../../shared/helper/HelperType"
import {KeeperPluginOptions, ReleaseStageOptions} from "../../../../shared/helper/HelperUtils"
import {useStore, useStoreAction} from "../../../../shared/provider/StoreProvider"

const SX: SxPropsMap = {
    group: {display: "flex", alignItems: "stretch", minWidth: 0},
    text: {
        padding: "7px 10px", fontSize: "12px", textTransform: "uppercase", color: "common.white",
        fontFamily: "monospace", letterSpacing: 1,
        border: "1px solid", borderColor: "divider", borderRight: "none", borderRadius: "4px 0 0 4px",
        minWidth: 0, overflow: "hidden", whiteSpace: "nowrap", textOverflow: "ellipsis", lineHeight: 1,
    },
    arrow: {
        height: "28px", width: "32px", minWidth: "unset", flexShrink: 0, padding: 0,
        borderRadius: "0 4px 4px 0", borderColor: "divider",
    },
    icon: {fontSize: "18px", color: "common.white"},
    popper: {zIndex: 3},
    menu: {padding: "4px 0px", minWidth: "170px"},
    item: {
        display: "flex", alignItems: "center", justifyContent: "space-between", gap: 1.5, padding: "5px 10px",
        fontFamily: "monospace", letterSpacing: 1, fontSize: "12px", textTransform: "uppercase",
    },
    stage: {fontSize: "11px", color: "text.secondary", fontFamily: "monospace"},
}

export function ListKeepers() {
    const keeper = useStore(s => s.activeClusterKeeperPlugin)
    const {setClusterKeeperPlugin} = useStoreAction
    const [open, setOpen] = useState(false)
    const anchorRef = useRef(null)

    const option = KeeperPluginOptions[keeper]

    return (
        <>
            <Box sx={SX.group} ref={anchorRef}>
                <Box sx={SX.text}>{option.label}</Box>
                <Button variant={"outlined"} color={"inherit"} sx={SX.arrow} onClick={() => setOpen(!open)}>
                    {open ? <ArrowDropUp sx={SX.icon}/> : <ArrowDropDown sx={SX.icon}/>}
                </Button>
            </Box>
            <Popper sx={SX.popper} open={open} anchorEl={anchorRef.current} placement={"bottom-start"}>
                <Paper>
                    <ClickAwayListener onClickAway={() => setOpen(false)}>
                        <MenuList sx={SX.menu}>
                            {Object.entries(KeeperPluginOptions).map(([key, value]) => renderItem(key as KeeperPlugin, value))}
                        </MenuList>
                    </ClickAwayListener>
                </Paper>
            </Popper>
        </>
    )

    function renderItem(key: KeeperPlugin, value: EnumOptions & {stage: ReleaseStage}) {
        const itemStage = ReleaseStageOptions[value.stage]
        return (
            <MenuItem key={key} sx={SX.item} selected={key === keeper} onClick={() => handleClick(key)}>
                <Box>{value.label}</Box>
                <Tooltip title={itemStage.description} placement={"top"}>
                    <Box sx={SX.stage}>{itemStage.label.toLowerCase()}</Box>
                </Tooltip>
            </MenuItem>
        )
    }

    function handleClick(key: KeeperPlugin) {
        setClusterKeeperPlugin(key)
        setOpen(false)
    }
}
