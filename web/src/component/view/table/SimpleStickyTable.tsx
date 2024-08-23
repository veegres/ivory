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
        position: "absolute", top: 0, left: 0, padding: "0px 10px",
        display: "flex", alignItems: "center",
        borderBottom: 1, borderRight: 1,  bgcolor: "background.paper", borderColor: "divider",
    },
    cellFixed: {
        position: "absolute", top: 0, left: 0, padding: "0px 10px",
        display: "flex", alignItems: "center", justifyContent: "center",
        color: "text.disabled", bgcolor: "background.paper", borderColor: "divider",
    },
    noWrap: {whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis", width: "100%"},
    preWrap: {whiteSpace: "pre-wrap", overflow: "hidden", textOverflow: "ellipsis", width: "100%"},
    headTitle: {display: "flex", fontFamily: "monospace", gap: "4px"},
}

type Props = {
    columns: { name: string, description?: string }[],
    rows: any[][],
    width?: string,
    height?: string,
    fetching?: boolean,
    renderRowCell?: (row: any[]) => ReactNode,
}

export function SimpleStickyTable(props: Props) {
    const {columns, rows, renderRowCell} = props
    const {fetching = false, width = "100%", height = "391px"} = props
    const parentRef = useRef<Element>(null)

    const rowHeight = 32
    const rowHeightPx = `${rowHeight}px`
    const rowVirtualizer = useVirtualizer({
        count: rows.length,
        getScrollElement: () => parentRef.current,
        estimateSize: () => rowHeight,
        scrollPaddingStart: rowHeight,
        overscan: 3,
    })

    const columnCount = renderRowCell !== undefined ? columns.length + 1 : columns.length
    const columnWidthId = 50
    // NOTE: calculate the cell width, the component should be rendered from scratch to change the size,
    //  this is how useVirtualizer work (-1 is needed to always choose 100 if width is unknown; -4 is the padding)
    const columnWidth = Math.max((Math.round((parentRef.current?.clientWidth ?? -1) - columnWidthId - 4) / columnCount), 100)
    const columnWidthPx = `${columnWidthId}px`
    const columnVirtualizer = useVirtualizer({
        horizontal: true,
        count: columnCount,
        getScrollElement: () => parentRef.current,
        estimateSize: () => columnWidth,
        scrollPaddingStart: columnWidthId,
        overscan: 3,
        enabled: parentRef.current !== null,
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
        const height = `${rowHeight}px`

        return columnVirtualizer.getVirtualItems().map((virtualColumn) => {
            const column = columns.length <= virtualColumn.index ? {name: "", description: ""} : columns[virtualColumn.index]
            const transform = `translateX(${virtualColumn.start}px)`
            const width = `${virtualColumn.size}px`
            return (
                <Box key={virtualColumn.key} sx={SX.cellFixed} borderRight={1} borderTop={1} borderBottom={2} style={{width, height, transform}}>
                    <Tooltip title={renderHeadTitle(column.name, column.description)} placement={"top"} arrow={true}>
                        <Box sx={SX.noWrap}>{column.name}</Box>
                    </Tooltip>
                </Box>
            )}
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
        const wrap = columnCount === 1 ? SX.preWrap : SX.noWrap

        return rowVirtualizer.getVirtualItems().map((virtualRow) => {
            const row = rows[virtualRow.index]
            return (
                <Fragment key={virtualRow.key}>
                    {columnVirtualizer.getVirtualItems().map((virtualColumn) => {
                        const parseColumn = getParseColumn(row, virtualColumn.index)

                        const height = `${virtualRow.size}px`
                        const width = `${virtualColumn.size}px`
                        const transform = `translateX(${virtualColumn.start}px) translateY(${virtualRow.start}px)`
                        return (
                            <Box key={virtualColumn.key} sx={SX.cell} style={{width, height, transform}}>
                                <Box sx={wrap}>{parseColumn(row)}</Box>
                            </Box>
                        )
                    })}
                </Fragment>
            )
        })
    }

    function getParseColumn(row: any[], index: number): ((row: any[]) => ReactNode) {
        if (columns.length <= index) return renderRowCell ?? (() => undefined)
        const column = row[index]
        if (typeof column === 'object') return () => JSON.stringify(column)
        else if (typeof column === "boolean") return () => Boolean(column).toString()
        else return () => column
    }
}
