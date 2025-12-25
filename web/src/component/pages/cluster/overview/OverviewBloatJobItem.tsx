import {Clear, Stop} from "@mui/icons-material"
import {Box, CircularProgress, Divider, IconButton, Paper, Tooltip} from "@mui/material"
import {SvgIconProps} from "@mui/material/SvgIcon/SvgIcon"
import {cloneElement, ReactElement, useState} from "react"

import {useRouterBloatDelete, useRouterBloatStop} from "../../../../api/bloat/hook"
import {Bloat} from "../../../../api/bloat/type"
import {Permission} from "../../../../api/permission/type"
import {SxPropsMap} from "../../../../app/type"
import {shortUuid} from "../../../../app/utils"
import {useEventJob} from "../../../../hook/EventJob"
import scroll from "../../../../style/scroll.module.css"
import select from "../../../../style/select.module.css"
import {OpenIcon} from "../../../view/icon/OpenIcon"
import {LinearProgressStateful} from "../../../view/progress/LinearProgressStateful"
import {DynamicRowVirtualizer} from "../../../view/scrolling/DynamicRowVirtualizer"
import {Access} from "../../../widgets/access/Access"

const SX: SxPropsMap = {
    paper: {fontSize: "13px", width: "100%", padding: "8px 15px"},
    console: {fontSize: "13px", width: "100%", padding: "8px 15px"},
    row: {"&:hover": {color: "primary.main"}},
    emptyLine: {textAlign: "center", textTransform: "uppercase"},
    header: {display: "flex", flexDirection: "column", fontWeight: "bold", cursor: "pointer"},
    headerLine: {display: "flex", justifyContent: "space-between", alignItems: "center", flexWrap: "nowrap", height: "20px"},
    loader: {margin: "10px 0 5px"},
    divider: {margin: "5px 0"},
    logs: {colorScheme: "dark"},
    button: {padding: "1px"},
    tooltipBox: {marginLeft: "4px", width: "25px", display: "flex", alignItems: "center", justifyContent: "center"},
    separator: {display: "flex", alignItems: "start", marginLeft: "10px"},
    credential: {display: "inline", color: "text.secondary", marginLeft: "5px"},
}

type Props = {
    cluster: string,
    item: Bloat,
    refetchList: () => void,
}

export function OverviewBloatJobItem(props: Props) {
    const {item, cluster, refetchList} = props
    const {uuid, status: initStatus, command, credentialId} = item
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

    function renderHeader() {
        return (
            <Box sx={SX.header} onClick={() => setOpen(!open)} className={select.none}>
                <Box sx={SX.headerLine}>
                    <Box>Command</Box>
                    <Box sx={SX.separator}>
                        <Box sx={{color: status.color}}>{status.name}</Box>
                        <Access permission={Permission.ManageBloatJob}>
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
                        <Tooltip title={renderPasswordTooltip()} placement={"top"}>
                            <Box sx={SX.credential}>
                                --username {shortUuid(credentialId)} --password {shortUuid(credentialId)}
                            </Box>
                        </Tooltip>
                    </Box>
                    <Box sx={SX.separator}>
                        <Tooltip title={`Job ID: ${uuid}`}>
                            <Box>{shortUuid(uuid)}</Box>
                        </Tooltip>
                        <Access permission={Permission.ViewBloatLogs}>
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

    function renderBody() {
        if (!open) return
        return (
            <>
                <Divider sx={SX.divider} textAlign={"left"}>LOGS</Divider>
                {logs.length === 0 ? isFetching ? (
                    <Box sx={SX.emptyLine}>Waiting for logs</Box>
                ) : (
                    <Box sx={SX.emptyLine}>No logs</Box>
                ) : (
                    <DynamicRowVirtualizer
                        sx={SX.logs}
                        auto={status.active && open}
                        className={scroll.small}
                        sxVirtualRow={SX.row}
                        height={350}
                        rows={logs}
                    />
                )}
                <LinearProgressStateful sx={SX.loader} loading={isFetching} color={"inherit"} line/>
            </>
        )
    }

    function renderPasswordTooltip() {
        return (
            <Box>
                <Box><b>Credential ID</b></Box>
                <Box>[ provided only for transparency, you need to provide real password and username ]</Box>
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
