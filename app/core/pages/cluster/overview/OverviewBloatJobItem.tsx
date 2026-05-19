import {Clear, Stop} from "@mui/icons-material"
import {Box, CircularProgress, Divider, IconButton, Paper, Tooltip} from "@mui/material"
import {SvgIconProps} from "@mui/material"
import {cloneElement, ReactElement, useState} from "react"

import {useRouterBloatDelete, useRouterBloatStop} from "../../../../features/bloat/hook"
import {Bloat} from "../../../../features/bloat/type"
import {Feature} from "../../../../features/feature"
import {OpenIcon} from "../../../../shared/component/icon/OpenIcon"
import {SxPropsMap} from "../../../../shared/helper/type"
import {getShortUuid} from "../../../../shared/helper/utils"
import {useEventJob} from "../../../../shared/hook/EventJob"
import select from "../../../../shared/style/select.module.css"
import {Access} from "../../../widgets/access/Access"
import {Logs} from "../../../widgets/logs/Logs"

const SX: SxPropsMap = {
    paper: {fontSize: "13px", width: "100%", padding: "8px 15px"},
    header: {display: "flex", flexDirection: "column", fontWeight: "bold", cursor: "pointer"},
    headerLine: {display: "flex", justifyContent: "space-between", alignItems: "center", flexWrap: "nowrap", height: "20px"},
    logs: {colorScheme: "dark"},
    button: {padding: "1px"},
    tooltipBox: {marginLeft: "4px", width: "25px", display: "flex", alignItems: "center", justifyContent: "center"},
    separator: {display: "flex", alignItems: "start", marginLeft: "10px"},
    vault: {display: "inline", color: "text.secondary", marginLeft: "5px"},
    divider: {margin: "5px 0", fontFamily: "monospace"},
}

type Props = {
    cluster: string,
    item: Bloat,
    refetchList: () => void,
}

export function OverviewBloatJobItem(props: Props) {
    const {item, cluster, refetchList} = props
    const {uuid, status: initStatus, command, vaultId} = item
    const [open, setOpen] = useState(false)
    const {isFetching, logs, status} = useEventJob(uuid, initStatus, open, refetchList)

    const deleteJob = useRouterBloatDelete(uuid, cluster)
    const stopJob = useRouterBloatStop()

    return (
        <Paper sx={SX.paper} variant={"outlined"}>
            {renderHeader()}
            {renderBody()}
        </Paper>
    )

    function renderBody() {
        if (!open) return
        return (
            <>
                <Divider sx={SX.divider} textAlign={"left"}>LOGS</Divider>
                <Logs sx={SX.logs} logs={logs} auto={status.active && open} loading={isFetching}/>
            </>
        )
    }

    function renderHeader() {
        return (
            <Box sx={SX.header} onClick={() => setOpen(!open)} className={select.none}>
                <Box sx={SX.headerLine}>
                    <Box>Command</Box>
                    <Box sx={SX.separator}>
                        <Box sx={{color: status.color}}>{status.name}</Box>
                        <Access feature={Feature.ManageToolBloatJob}>
                            {status.active ?
                                renderJobButton("Stop", <Stop/>, () => stopJob.mutate(uuid), stopJob.isPending) :
                                renderJobButton("Delete", <Clear/>, () => deleteJob.mutate(uuid), deleteJob.isPending)
                            }
                        </Access>
                    </Box>
                </Box>
                <Box sx={SX.headerLine}>
                    <Box>
                        {command}
                        <Tooltip title={renderVaultTooltip()} placement={"top"}>
                            <Box sx={SX.vault}>
                                {vaultId ? `--username ${getShortUuid(vaultId)} --password ${getShortUuid(vaultId)}` : "--username postgres"}
                            </Box>
                        </Tooltip>
                    </Box>
                    <Box sx={SX.separator}>
                        <Tooltip title={`Job ID: ${uuid}`}>
                            <Box>{getShortUuid(uuid)}</Box>
                        </Tooltip>
                        <Access feature={Feature.ViewToolBloatLogs}>
                            <Tooltip title={"Open"}>
                                <Box sx={SX.tooltipBox}>
                                    <IconButton sx={SX.button} size={"small"}>
                                        <OpenIcon open={open} size={18}/>
                                    </IconButton>
                                </Box>
                            </Tooltip>
                        </Access>
                    </Box>
                </Box>
            </Box>
        )
    }

    function renderVaultTooltip() {
        return (
            <Box>
                <Box><b>Vault ID</b></Box>
                <Box>[ provided only for transparency, you need to provide real vault and username ]</Box>
            </Box>
        )
    }

    function renderJobButton(title: string, icon: ReactElement<SvgIconProps>, onClick: () => void, isLoading: boolean) {
        const fontSize = 18
        return (
            <Tooltip title={title} placement={"top"}>
                <Box sx={SX.tooltipBox}>
                    {isLoading ? <CircularProgress size={fontSize - 3}/> : (
                        <IconButton
                            sx={SX.button}
                            size={"small"}
                            onClick={(e) => {
                                e.stopPropagation()
                                onClick()
                            }}
                        >
                            {cloneElement(icon, {sx: {fontSize}})}
                        </IconButton>
                    )}
                </Box>
            </Tooltip>
        )
    }
}
