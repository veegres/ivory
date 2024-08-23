import scroll from "../../../style/scroll.module.css";
import {Box, CircularProgress, Tooltip} from "@mui/material";
import {SxPropsMap} from "../../../type/general";
import {Fragment, ReactNode, useRef} from "react";
import {useVirtualizer} from "@tanstack/react-virtual";
import {NoBox} from "../box/NoBox";

const SX: SxPropsMap = {
    box: {position: "relative"},
    body: {overflow: "auto", width: "100%", height: "100%", padding: "0px 4px 4px 0px", fontSize: "12px"},
    loader: {
        position: "absolute", top: 0, left: 0, right: 0, bottom: 0,
        display: "flex", alignItems: "center", justifyContent: "center",
        width: "inherit", height: "inherit", zIndex: 4, bgcolor: "rgba(31,44,126,0.05)",
    },
    cell: {
        position: "absolute", top: 0, left: 0,
        borderBottom: 1, borderRight: 1, borderColor: "divider",
        display: "flex", alignItems: "center", padding: "0px 10px",
    },
    cellFixed: {
        position: "absolute", top: 0, left: 0,
        display: "flex", alignItems: "center", justifyContent: "center",
        whiteSpace: "nowrap", color: "text.disabled",
        bgcolor: "background.paper", borderColor: "divider",
    },
    noWrap: {whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis"},
    headTitle: {display: "flex", fontFamily: "monospace"},
}

type Props = {
    columns: { name: string, description?: string }[],
    rows: any[][],
    width?: string,
    height?: string,
    fetching?: boolean,
    renderHeaderCell?: () => ReactNode,
    renderRowCell?: (row: any[]) => ReactNode,
}

export function SimpleStickyTable(props: Props) {
    const {columns, rows, renderHeaderCell, renderRowCell} = props
    const {fetching = false, width = "100%", height = "391px"} = props
    const parentRef = useRef(null)

    const rowHeight = 32
    const rowHeightPx = `${rowHeight}px`
    const rowVirtualizer = useVirtualizer({
        count: rows.length,
        getScrollElement: () => parentRef.current,
        estimateSize: () => rowHeight,
        scrollPaddingStart: rowHeight,
        overscan: 3,
    })

    const columnWidthId = 50
    const columnWidth = 100
    const columnWidthPx = `${columnWidthId}px`
    const columnVirtualizer = useVirtualizer({
        horizontal: true,
        count: columns.length,
        getScrollElement: () => parentRef.current,
        estimateSize: () => columnWidth,
        scrollPaddingStart: columnWidthId,
        overscan: 3,
    })

    const bodyWidthPx = `${columnVirtualizer.getTotalSize()}px`
    const bodyHeightPx = `${rowVirtualizer.getTotalSize()}px`

    return (
        <Box sx={SX.box} width={width} height={height}>
            {renderLoader()}
            {columns.length === 0 ? renderEmpty() : renderTable()}
        </Box>
    )

    function renderLoader() {
        if (!fetching) return
        return (
            <Box sx={SX.loader}>
                <CircularProgress />
            </Box>
        )
    }

    function renderEmpty() {
        return (
            <Box sx={{padding: "10px"}}>
                <NoBox text={"NO DATA"}/>
            </Box>
        )
    }

    function renderTable() {
        return (
            <Box
                ref={parentRef}
                sx={SX.body}
                display={"grid"}
                gridTemplateColumns={`${columnWidthPx} ${bodyWidthPx}`}
                gridTemplateRows={`${rowHeightPx} ${bodyHeightPx}`}
                className={scroll.small}
            >
                <Box position={"sticky"} zIndex={3} top={0} left={0}>
                    {renderCorner()}
                </Box>
                <Box position={"sticky"} zIndex={2} top={0}>
                    {renderHead()}
                </Box>
                <Box position={"sticky"} zIndex={2} left={0}>
                    {renderIndex()}
                </Box>
                <Box position={"relative"}>
                    {renderBody()}
                </Box>
            </Box>
        )
    }

    function renderCorner() {
        return (
            <Box sx={SX.cellFixed} border={1} borderBottom={2} borderRight={2} style={{width: columnWidthPx, height: rowHeightPx}}/>
        )
    }

    function renderHead() {
        if (!columns) return

        const width = `${columnWidth}px`
        const height = `${rowHeight}px`

        return (
            <Fragment>
                {columnVirtualizer.getVirtualItems().map((virtualColumn) => {
                    const column = columns[virtualColumn.index]
                    const transform = `translateX(${virtualColumn.start}px)`
                    return (
                        <Box key={virtualColumn.key} sx={SX.cellFixed} borderRight={1} borderTop={1} borderBottom={2} style={{width, height, transform}}>
                            <Tooltip title={renderHeadTitle(column.name, column.description)} placement={"top"} arrow={true}>
                                <Box>{column.name}</Box>
                            </Tooltip>
                        </Box>
                    )}
                )}
                {renderHeaderCell && renderHeaderCell()}
            </Fragment>
        )
    }

    function renderHeadTitle(name: string, description?: string) {
        return (
            <Box sx={SX.headTitle}>
                <Box>{name}</Box>
                {description && <Box>({description})</Box>}
            </Box>
        )
    }

    function renderIndex() {
        const width = `${columnWidthId}px`
        const height = `${rowHeight}px`

        return rowVirtualizer.getVirtualItems().map((virtualRow) => {
            const transform = `translateY(${virtualRow.start}px)`
            return (
                <Box key={virtualRow.key} sx={SX.cellFixed} borderLeft={1} borderBottom={1} borderRight={2} style={{width, height, transform}}>
                    {virtualRow.index + 1}
                </Box>
            )
        })
    }

    function renderBody() {
        if (!columns) return

        return rowVirtualizer.getVirtualItems().map((virtualRow) => {
            const row = rows[virtualRow.index]

            return (
                <Fragment key={virtualRow.key}>
                    {columnVirtualizer.getVirtualItems().map((virtualColumn) => {
                        const column = row[virtualColumn.index]
                        const parseColumn = getParseColumn(column)

                        const height = `${virtualRow.size}px`
                        const width = `${virtualColumn.size}px`
                        const transform = `translateX(${virtualColumn.start}px) translateY(${virtualRow.start}px)`
                        return (
                            <Box key={virtualColumn.key} sx={SX.cell} style={{width, height, transform}}>
                                <Box sx={SX.noWrap}>{parseColumn}</Box>
                            </Box>
                        )
                    })}
                    {renderRowCell && renderRowCell(row)}
                </Fragment>
            )
        })
    }

    function getParseColumn(column: any): string {
        if (typeof column === 'object') return JSON.stringify(column)
        else if (typeof column === "boolean") return Boolean(column).toString()
        else return column
    }
}
