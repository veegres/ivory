import {Box, Table, TableCell, TableHead, TableRow, Tooltip} from "@mui/material";
import {Database, SxPropsMap} from "../../../type/common";
import {ErrorSmart} from "../../view/box/ErrorSmart";
import {TableBody} from "../../view/table/TableBody";
import scroll from "../../../style/scroll.module.css"
import {QueryRunResponse, QueryVariety} from "../../../type/query";
import {CancelIconButton, RefreshIconButton, TerminateIconButton} from "../../view/button/IconButtons";
import {InfoColorBox} from "../../view/box/InfoColorBox";
import {QueryVarietyOptions} from "../../../app/utils";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    body: {overflow: "auto", maxHeight: "300px"},
    table: {"tr td, th": {border: "1px solid", borderColor: "divider"}, padding: "0px 5px 5px 0px"},
    cell: {maxWidth: "600px", whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis"},
    type: {fontFamily: "monospace", color: "text.disabled"},
    name: {color: "text.secondary"},
    number: {
        width: "1%", whiteSpace: "nowrap", color: "text.secondary", position: "sticky", left: 0,
        bgcolor: "background.paper", zIndex: 2, textAlign: "center",
    },
    no: {
        display: "flex", alignItems: "center", justifyContent: "center", textTransform: "uppercase",
        padding: "8px 16px", border: "1px solid", borderColor: "divider", lineHeight: "1.54",
    },
    pid: {display: "flex", padding: "0 5px"},
    info: {display: "flex", alignItems: "center", justifyContent: "space-between", gap: 2},
    label: {color: "text.secondary", padding: "0 5px", cursor: "pointer", fontSize: "13.5px"},
    buttons: {display: "flex", alignItems: "center", gap: 1},
    headTitle: {display: "flex", gap: 1},
}

type Props = {
    loading: boolean,
    varieties?: QueryVariety[],
    db: Database,
    data?: QueryRunResponse,
    onRefresh: () => void,
    onTerminate: (pid: number) => void,
    onCancel: (pid: number) => void,
    error: unknown,
}

export function QueryItemRun(props: Props) {
    const {loading, data, error, db, varieties} = props
    const {onTerminate, onCancel, onRefresh} = props
    const pidIndex = data?.fields.findIndex(field => field.name === "pid") ?? -1

    return (
        <Box sx={SX.box}>
            {renderInfo()}
            {renderTable()}
        </Box>
    )

    function renderInfo() {
        return (
            <Box sx={SX.info}>
                <Tooltip title={"SENT TO"} placement={"right"} arrow={true}>
                    <Box sx={SX.label}>[ postgres://{db.host}:{db.port}/{db.database ? db.database : "postgres"} ]</Box>
                </Tooltip>
                <Box sx={SX.buttons}>
                    {varieties && varieties.map(v => {
                        const {badge, label, color} = QueryVarietyOptions[v]
                        return (
                            <InfoColorBox key={v} label={badge ?? "UNKNOWN"} title={label} bgColor={color} opacity={0.8}/>
                        )
                    })}
                    <RefreshIconButton color={"success"} disabled={loading} onClick={onRefresh}/>
                </Box>
            </Box>
        )
    }

    function renderTable() {
        if (error) return <ErrorSmart error={error}/>
        if (!loading && (!data || !data.fields.length || !data.rows.length)) {
            return <Box sx={SX.no}>Response is empty</Box>
        }

        return (
            <Box sx={SX.body} className={scroll.small}>
                <Table sx={SX.table} size={"small"} stickyHeader>
                    <TableHead>
                        {renderHead()}
                    </TableHead>
                    <TableBody isLoading={loading} rowCount={getRowCount()}>
                        {renderBody()}
                    </TableBody>
                </Table>
            </Box>
        )
    }

    function renderHead() {
        if (loading) return
        if (!data) return

        return (
            <TableRow>
                <TableCell sx={{...SX.number, zIndex: 3}}/>
                {data.fields.map(field => (
                    <TableCell key={field.name} sx={SX.cell}>
                        <Tooltip title={renderHeadTitle(field.name, field.dataType)} placement={"top-start"}>
                            <Box sx={SX.name}>{field.name}</Box>
                        </Tooltip>
                    </TableCell>
                ))}
                {pidIndex !== -1 && <TableCell/>}
            </TableRow>
        )
    }

    function renderHeadTitle(name: string, type: string) {
        return (
            <Box sx={SX.headTitle}>
                <Box>{name}</Box>
                <Box sx={SX.type}>({type})</Box>
            </Box>
        )
    }

    function renderBody() {
        if (!data) return

        return data.rows.map((rows, i) => (
            <TableRow key={i}>
                <TableCell sx={SX.number}>{i + 1}</TableCell>
                {rows.map((row, j) => {
                    const parsedRow = getParseRow(row)
                    return (
                        <TableCell key={j} sx={SX.cell}>
                            <Tooltip title={parsedRow} placement={"top-start"}>
                                <Box>{parsedRow}</Box>
                            </Tooltip>
                        </TableCell>
                    )
                })}
                {pidIndex !== -1 && (
                    <TableCell sx={SX.pid}>
                        <CancelIconButton onClick={() => onCancel(rows[pidIndex])}/>
                        <TerminateIconButton onClick={() => onTerminate(rows[pidIndex])} color={"error"}/>
                    </TableCell>
                )}
            </TableRow>
        ))
    }

    function getParseRow(row: any): string {
        if (typeof row === 'object') return JSON.stringify(row)
        else if (typeof row === "boolean") return Boolean(row).toString()
        else return row
    }

    function getRowCount() {
        if (!data) return 2
        const len = data.rows.length

        if (len === 0) return 1
        else if (len > 8) return 8
        else return data.rows.length + 1
    }
}
