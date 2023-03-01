import {Box, Skeleton, Table, TableCell, TableHead, TableRow, Tooltip} from "@mui/material";
import {QueryRunResponse, SxPropsMap} from "../../../app/types";
import {ErrorAlert} from "../../view/ErrorAlert";
import {TableBody} from "../../view/TableBody";
import {useTheme} from "../../../provider/ThemeProvider";
import scroll from "../../../style/scroll.module.css"

const SX: SxPropsMap = {
    table: {"tr:last-child td": {border: 0}, marginBottom: "5px"},
    box: {overflow: "auto"},
    cell: {maxWidth: "400px", whiteSpace: 'nowrap', overflow: "hidden", textOverflow: "ellipsis"},
    type: {fontSize: "12px", fontFamily: "monospace"},
    name: {fontWeight: "bold"},
}

type Props = {
    loading: boolean,
    data?: QueryRunResponse,
    error: unknown,
}

export function QueryItemRun(props: Props) {
    const {loading, data, error} = props
    const {info} = useTheme()

    if (error) return <ErrorAlert error={error}/>

    return (
        <Box sx={SX.box} className={scroll.tiny}>
            <Table size={"small"} sx={SX.table}>
                <TableHead>
                    <TableRow>
                        {renderHead()}
                    </TableRow>
                </TableHead>
                <TableBody isLoading={loading}>
                    {renderBody()}
                </TableBody>
            </Table>
        </Box>
    )

    function renderHead() {
        if (loading) return <TableCell><Skeleton width={'100%'}/></TableCell>
        if (!data) return

        return data.fields.map(field => (
            <TableCell key={field.name} sx={SX.cell}>
                <Box component={"span"}>{field.name}</Box>
                {" "}
                <Box sx={{...SX.type, color: info?.palette.text.disabled}} component={"span"}>({field.dataType})</Box>
            </TableCell>
        ))
    }

    function renderBody() {
        if (!data) return

        return data.rows.map((rows, i) => (
            <TableRow key={i}>
                {rows.map((row, j) => {
                    const parsedRow = parseRow(row)
                    return (
                        <TableCell key={j} sx={{...SX.cell, color: info?.palette.text.secondary}}>
                            <Tooltip title={parsedRow} placement={"top-start"}>
                                <span>{parsedRow}</span>
                            </Tooltip>
                        </TableCell>
                    )
                })}
            </TableRow>
        ))
    }

    function parseRow(row: any): string {
        if (typeof row === 'object') return JSON.stringify(row)
        else return row
    }
}
