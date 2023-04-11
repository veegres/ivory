import {Box, Table, TableCell, TableHead, TableRow, Tooltip} from "@mui/material";
import {SxPropsMap} from "../../../type/common";
import {ErrorAlert} from "../../view/box/ErrorAlert";
import {TableBody} from "../../view/table/TableBody";
import scroll from "../../../style/scroll.module.css"
import {QueryRunResponse} from "../../../type/query";
import {CancelIconButton, TerminateIconButton} from "../../view/button/IconButtons";

const SX: SxPropsMap = {
    // we need this box margin and table padding to see table border when scroll appears
    box: {overflow: "auto", maxHeight: "300px", margin: "5px 0px 0px 5px"},
    table: {"tr td, th": {border: "1px solid", borderColor: "divider"}, padding: "0px 5px 5px 0px"},
    cell: {maxWidth: "600px", whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis"},
    type: {fontSize: "12px", fontFamily: "monospace", color: "text.disabled"},
    name: {color: "text.secondary"},
    number: {
        width: "1%", whiteSpace: "nowrap", color: "text.secondary", position: "sticky", left: 0,
        bgcolor: "background.paper", zIndex: 2, textAlign: "center"
    },
    no: {display: "flex", alignItems: "center", justifyContent: "center", textTransform: "uppercase", padding: "10px"},
    pid: {display: "flex", padding: "0 5px"},
}

type Props = {
    loading: boolean,
    data?: QueryRunResponse,
    onTerminate: (pid: number) => void,
    onCancel: (pid: number) => void,
    error: unknown,
}

export function QueryItemRun(props: Props) {
    const {loading, data, error, onTerminate, onCancel} = props

    if (error) return <ErrorAlert error={error}/>
    if (!loading && (!data || !data.fields.length || !data.rows.length)) {
        return <Box sx={SX.no}>Response is empty</Box>
    }
    const pidIndex = data?.fields.findIndex(field => field.name === "pid") ?? -1

    return (
        <Box sx={SX.box} className={scroll.small}>
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

    function renderHead() {
        if (loading) return
        if (!data) return

        return (
            <TableRow>
                <TableCell sx={{...SX.number, zIndex: 3}}/>
                {data.fields.map(field => (
                    <TableCell key={field.name} sx={SX.cell}>
                        <Box sx={SX.name} component={"span"}>{field.name}</Box>
                        {" "}
                        <Box sx={SX.type} component={"span"}>({field.dataType})</Box>
                    </TableCell>
                ))}
                {pidIndex !== -1 && <TableCell/>}
            </TableRow>
        )
    }

    function renderBody() {
        if (!data) return

        return data.rows.map((rows, i) => (
            <TableRow key={i}>
                <TableCell sx={SX.number}>{i+1}</TableCell>
                {rows.map((row, j) => {
                    const parsedRow = getParseRow(row)
                    return (
                        <TableCell key={j} sx={SX.cell}>
                            <Tooltip title={parsedRow} placement={"top-start"}>
                                <span>{parsedRow}</span>
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
        else if (typeof  row === "boolean") return Boolean(row).toString()
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
