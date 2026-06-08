import {ArrowDropDown, ArrowDropUp} from "@mui/icons-material"
import {Box, Collapse, Tooltip} from "@mui/material"
import {useState} from "react"

import {useRouterActivity} from "../../../features/query/hook"
import {QueryApi} from "../../../features/query/router"
import {Connection} from "../../../features/query/type"
import {SimpleButton} from "../../../shared/component/button/SimpleButton"
import {SxPropsMap} from "../../../shared/helper/type"
import {Refresher} from "../refresher/Refresher"
import {QueryTable} from "./QueryTable"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", border: 1, borderRadius: 1, borderColor: "divider", padding: "8px"},
    head: {display: "flex", alignItems: "center", gap: 1},
    label: {display: "flex", justifyContent: "start", fontSize: "14px", fontFamily: "monospace", color: "text.secondary", padding: "0px 3px"},
    help: {fontSize: "9px", fontWeight: "normal", color: "text.disabled", textAlign: "center"},
    action: {display: "flex", justifyContent: "end", gap: 1, color: "text.secondary", cursor: "pointer"},
    info: {
        display: "flex", justifyContent: "center", alignItems: "center", padding: "0 15px",
        fontSize: "11px", color: "text.secondary", textAlign: "center",
    },
    icon: {fontSize: "16px", color: "text.secondary"},
    error: {color: "error.light"},
    collapse: {marginTop: "10px", display: "flex", flexDirection: "column", gap: 1},
}

type Props = {
    connection: Connection,
}

export function QueryActivity(props: Props) {
    const {connection} = props
    const [open, setOpen] = useState(false)
    const {data, isError, error, refetch} = useRouterActivity(connection)
    const table = isError ? undefined : data
    return (
        <Box sx={SX.box}>
            <Box sx={SX.head}>
                <Box sx={SX.label} flex={1}>Active Session Queries</Box>
                {open && <Box sx={SX.help} flex={1}>[ hold shift for horizontal scrolling ]</Box>}
                <Box sx={SX.action} flex={1}>
                    <Tooltip title={"COUNT"} placement={"top"}>
                        <Box>[ {error || !data ? "-" : data.rows.length - 1} ]</Box>
                    </Tooltip>
                    <Refresher queryKeys={[QueryApi.activity.key()]} defaultPeriod={["5s", 5000]}/>
                    <Tooltip title={!open ? "Show Queries" : "Hide Queries"} placement={"top"} disableInteractive>
                        <SimpleButton onClick={() => setOpen(!open)}>
                            {open ? <ArrowDropUp sx={SX.icon}/> : <ArrowDropDown sx={SX.icon}/>}
                        </SimpleButton>
                    </Tooltip>
                </Box>
            </Box>
            <Collapse in={open}>
                <Box sx={SX.collapse}>
                    <QueryTable
                        connection={connection}
                        refetch={refetch}
                        height={100}
                        data={table}
                        error={error}
                        showIndexColumn={false}
                    />
                    <Box sx={SX.info}>
                        This section displays queries submitted within your current session in Ivory
                    </Box>
                </Box>
            </Collapse>
        </Box>
    )
}
