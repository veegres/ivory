import scroll from "../../../style/scroll.module.css";
import {Box, Table, TableCell, TableHead, TableRow, ToggleButton, Tooltip} from "@mui/material";
import {TableBody} from "./TableBody";
import {SxPropsMap} from "../../../type/general";
import {memo, ReactNode, useState} from "react";
import {WrapText} from "@mui/icons-material";

const SX: SxPropsMap = {
    body: {overflow: "auto", maxHeight: "315px"},
    table: {
        padding: "0px 3px 3px 0px",
        "th": {color: "text.disabled", lineHeight: "1.7"},
        "tr:first-of-type td, th": {borderTop: 1, borderColor: "divider"},
        "tr td, th": {
            borderRight: 1, borderColor: "divider", bgcolor: "background.paper", fontSize: "12px", maxWidth: "300px",
            "&:nth-of-type(2), &:nth-of-type(1)": {borderLeft: 1, borderColor: "divider"},
        },
    },
    preWrap: {whiteSpace: "pre", overflow: "hidden", textOverflow: "ellipsis"},
    noWrap: {whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis"},
    toggle: {border: 0, padding: "5px", width: "100%"},
    description: {fontFamily: "monospace", color: "text.secondary"},
    number: {
        width: "1%", whiteSpace: "nowrap", color: "text.disabled",
        position: "sticky", left: 0, zIndex: 2, textAlign: "center",
    },
    headTitle: {display: "flex", gap: 1},
}

type Props = {
    columns: { name: string, description?: string }[],
    rows: any[][],
    fetching?: boolean,
    renderHeaderCell?: () => ReactNode,
    renderRowCell?: (row: any[]) => ReactNode,
}

export const SimpleStickyTable = memo(SimpleStickyTableMemo)

export function SimpleStickyTableMemo(props: Props) {
    const {columns, rows, fetching = true, renderHeaderCell, renderRowCell} = props
    const [toggle, setToggle] = useState<boolean | undefined>(undefined)
    const wrap = toggle ?? columns.length === 1

    return (
        <Box sx={SX.body} className={scroll.small}>
            <Table sx={SX.table} size={"small"} stickyHeader>
                <TableHead>
                    {renderHead()}
                </TableHead>
                <TableBody isLoading={fetching} rowCount={getRowCount()} height={21}>
                    {renderBody()}
                </TableBody>
            </Table>
        </Box>
    )

    function renderHead() {
        if (fetching) return
        if (!columns) return

        return (
            <TableRow>
                <TableCell sx={{...SX.number, zIndex: 3, padding: "3px"}} onClick={() => setToggle(!toggle)}>
                    <ToggleButton sx={SX.toggle} value={"wrap"} selected={wrap}>
                        <Tooltip title={"Wrap text in the table"} placement={"top"}>
                            <WrapText sx={{fontSize: "15px"}}/>
                        </Tooltip>
                    </ToggleButton>
                </TableCell>
                {columns.map(column => (
                    <TableCell key={column.name}>
                        <Tooltip title={renderHeadTitle(column.name, column.description)} placement={"top"}>
                            <Box>{column.name}</Box>
                        </Tooltip>
                    </TableCell>
                ))}
                {renderHeaderCell && renderHeaderCell()}
            </TableRow>
        )
    }

    function renderHeadTitle(name: string, description?: string) {
        return (
            <Box sx={SX.headTitle}>
                <Box>{name}</Box>
                {description && <Box sx={SX.description}>({description})</Box>}
            </Box>
        )
    }

    function renderBody() {
        if (!columns) return

        return rows.map((row, i) => {
            return (
                <TableRow key={i}>
                    <TableCell sx={SX.number}>{i + 1}</TableCell>
                    {row.map((column, j) => {
                        const parseColumn = getParseColumn(column)
                        return (
                            <TableCell key={j}>
                                <Tooltip title={parseColumn} placement={"top"}>
                                    <Box sx={wrap ? SX.preWrap : SX.noWrap}>{parseColumn}</Box>
                                </Tooltip>
                            </TableCell>
                        )
                    })}
                    {renderRowCell && renderRowCell(row)}
                </TableRow>
            )
        })
    }

    function getParseColumn(column: any): string {
        if (typeof column === 'object') return JSON.stringify(column)
        else if (typeof column === "boolean") return Boolean(column).toString()
        else return column
    }

    function getRowCount() {
        if (!rows) return 2
        const len = rows.length

        if (len === 0) return 1
        else if (len > 9) return 9
        else return rows.length + 1
    }
}
