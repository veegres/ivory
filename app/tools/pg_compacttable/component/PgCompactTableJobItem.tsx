import {Clear, Stop} from "@mui/icons-material"
import {Box, CircularProgress, Divider, IconButton, Tooltip} from "@mui/material"
import {SvgIconProps} from "@mui/material"
import {cloneElement, ReactElement, useState} from "react"

import {Feature} from "../../../features/Feature"
import {ManageAccess} from "../../../features/management/component/ManageAccess"
import {Logs} from "../../../shared/component/box/Logs"
import {PaperBlue} from "../../../shared/component/box/PaperBlue"
import {OpenIcon} from "../../../shared/component/icon/OpenIcon"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {getShortUuid} from "../../../shared/helper/HelperUtils"
import select from "../../../shared/style/select.module.css"
import {useRouterPgCompactTableJob} from "../api/job/PgCompactTableJobHook"
import {useRouterPgCompactTableDelete, useRouterPgCompactTableStop} from "../api/PgCompactTableHook"
import {PgCompactTable} from "../api/PgCompactTableType"

const SX: SxPropsMap = {
    paper: {fontSize: "13px", width: "100%", padding: "6px 12px", border: 1, borderColor: "divider"},
    header: {display: "flex", flexDirection: "column", cursor: "pointer"},
    headerLine: {display: "flex", justifyContent: "space-between", alignItems: "center", flexWrap: "nowrap", height: "20px"},
    headerTitle: {fontWeight: "bold"},
    logs: {colorScheme: "dark"},
    button: {padding: "1px"},
    tooltipBox: {marginLeft: "4px", width: "25px", display: "flex", alignItems: "center", justifyContent: "center"},
    separator: {display: "flex", alignItems: "start", marginLeft: "10px"},
    vault: {display: "inline", color: "text.secondary", marginLeft: "5px"},
    divider: {margin: "5px 0", fontFamily: "monospace"},
}

type Props = {
    cluster: string,
    item: PgCompactTable,
    refetchList: () => void,
}

export function PgCompactTableJobItem(props: Props) {
    const {item, cluster, refetchList} = props
    const {uuid, status: initStatus, command, vaultId} = item
    const [open, setOpen] = useState(false)
    const {isFetching, logs, status} = useRouterPgCompactTableJob(uuid, initStatus, open, refetchList)

    const deleteJob = useRouterPgCompactTableDelete(uuid, cluster)
    const stopJob = useRouterPgCompactTableStop()

    return (
        <PaperBlue sx={SX.paper}>
            {renderHeader()}
            {renderBody()}
        </PaperBlue>
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
                <Box sx={[SX.headerLine, SX.headerTitle]}>
                    <Box>Command</Box>
                    <Box sx={SX.separator}>
                        <Box sx={{color: status.color}}>{status.name}</Box>
                        <ManageAccess feature={Feature.ManageToolPgCompactTableJob}>
                            {status.active ?
                                renderJobButton("Stop", <Stop/>, () => stopJob.mutate(uuid), stopJob.isPending) :
                                renderJobButton("Delete", <Clear/>, () => deleteJob.mutate(uuid), deleteJob.isPending)
                            }
                        </ManageAccess>
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
                        <ManageAccess feature={Feature.ViewToolPgCompactTableLogs}>
                            <Tooltip title={"Open"}>
                                <Box sx={SX.tooltipBox}>
                                    <IconButton sx={SX.button} size={"small"}>
                                        <OpenIcon open={open} size={18}/>
                                    </IconButton>
                                </Box>
                            </Tooltip>
                        </ManageAccess>
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
