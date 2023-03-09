import {Box, Table, TableCell, TableHead, TableRow, Tooltip} from "@mui/material";
import {QueryRunResponse, SxPropsMap} from "../../../app/types";
import {ErrorAlert} from "../../view/ErrorAlert";
import {TableBody} from "../../view/TableBody";
import scroll from "../../../style/scroll.module.css"

const SX: SxPropsMap = {
    table: {"tr:last-child td": {borderBottom: 0}, marginBottom: "5px"},
    box: {overflow: "auto", maxHeight: "300px"},
    cell: {maxWidth: "400px", whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis"},
    type: {fontSize: "12px", fontFamily: "monospace", color: "text.disabled"},
    name: {color: "text.secondary"},
    number: {width: "1%", whiteSpace: "nowrap", color: "text.secondary"},
    no: {display: "flex", alignItems: "center", justifyContent: "center", textTransform: "uppercase"},
}

type Props = {
    loading: boolean,
    data?: QueryRunResponse,
    error: unknown,
}

export function QueryItemRun(props: Props) {
    const {loading, data, error} = props

    if (error) return <ErrorAlert error={error}/>
    if (!loading && (!data?.fields.length || !data.rows.length)) {
        return <Box sx={SX.no}>Response is empty</Box>
    }

    return (
        <Box sx={SX.box} className={scroll.tiny}>
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
                <TableCell/>
                {data.fields.map(field => (
                    <TableCell key={field.name} sx={SX.cell}>
                        <Box sx={SX.name} component={"span"}>{field.name}</Box>
                        {" "}
                        <Box sx={SX.type} component={"span"}>({field.dataType})</Box>
                    </TableCell>
                ))}
            </TableRow>
        )
    }

    function renderBody() {
        if (!data) return

        return data.rows.map((rows, i) => (
            <TableRow key={i}>
                <TableCell sx={SX.number} padding={"checkbox"}>{i+1}</TableCell>
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
        else if (len > 9) return 9
        else return data.rows.length + 1
    }
}
