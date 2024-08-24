import scroll from "../../../style/scroll.module.css";
import {Box, CircularProgress, Tooltip} from "@mui/material";
import {SxPropsMap} from "../../../type/general";
import React, {Fragment, ReactNode, useCallback, useEffect, useRef, useState} from "react";
import {useVirtualizer} from "@tanstack/react-virtual";
import {NoBox} from "../box/NoBox";

const SX: SxPropsMap = {
    box: {position: "relative"},
    body: {overflow: "auto", width: "100%", height: "100%", fontSize: "12px"},
    loader: {
        position: "absolute", top: 0, left: 0, right: 0, bottom: 0, zIndex: 5,
        display: "flex", alignItems: "center", justifyContent: "center",
        width: "inherit", height: "inherit", bgcolor: "rgba(31,44,126,0.05)",
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
    columnSeparator: {
        position: "absolute", top: 0, left: 0, zIndex: 4,
        cursor: "col-resize", touchAction: "none",
        display: "flex", alignItems: "center", justifyContent: "center",
        "&:hover:after": {content: `""`, width: "30%", height: "60%", bgcolor: "gray"},
    },
    noWrap: {whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis", width: "100%"},
    preWrap: {whiteSpace: "pre-wrap", overflow: "hidden", textOverflow: "ellipsis", width: "100%"},
    headTitle: {display: "flex", fontFamily: "monospace", gap: "4px"},
    empty: {padding: "10px"},
}

type Props = {
    columns: { name: string, description?: string }[],
    rows: any[][],
    width?: string,
    height?: string,
    fetching?: boolean,
    renderRowCell?: (row: any[]) => ReactNode,
}

export function VirtualizedTable(props: Props) {
    const {columns, rows, renderRowCell} = props
    const {fetching = false, width = "100%", height = "391px"} = props
    const parentRef = useRef<Element>(null)
    const [columnDrag, setColumnDrag] = useState({index: -1, size: 0, ps: 0})

    const rowHeight = 32
    const rowHeightPx = `${rowHeight}px`
    const rowVirtualizer = useVirtualizer({
        count: rows.length,
        getScrollElement: () => parentRef.current,
        estimateSize: () => rowHeight,
        scrollPaddingStart: rowHeight,
        overscan: 3,
    })

    const scrollOffset = 6
    const columnCount = renderRowCell !== undefined ? columns.length + 1 : columns.length
    const columnWidthId = 50
    // NOTE: calculate the cell width, the component should be rendered from scratch to change the size,
    //  this is how useVirtualizer work (-1 is needed to always choose 100 if width is unknown)
    const columnWidth = Math.max((Math.round((parentRef.current?.clientWidth ?? -1) - columnWidthId - scrollOffset) / columnCount), 100)
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

    const callbackMouseMove = useCallback(handleMouseMove, [columnDrag.index, columnDrag.ps, columnDrag.size, columnVirtualizer])
    const callbackMouseUp = useCallback(handleMouseUp, [])
    useEffect(handleEffectDragging, [callbackMouseMove, callbackMouseUp]);

    const bodyWidthPx = `${columnVirtualizer.getTotalSize()}px`
    const bodyHeightPx = `${rowVirtualizer.getTotalSize()}px`

    return (
        <Box sx={SX.box} width={width} height={height}>
            {renderLoader()}
            {columns.length ? renderTable() : renderEmpty()}
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
        if (fetching) return
        return (
            <Box sx={SX.empty}>
                <NoBox text={"NO DATA"}/>
            </Box>
        )
    }

    function renderTable() {
        return (
            <Box
                ref={parentRef}
                sx={SX.body}
                padding={`0px ${scrollOffset}px ${scrollOffset}px 0px`}
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
        const heightPx = `${rowHeight}px`
        const dragWidth = scrollOffset * 2
        const dragWidthPx = `${dragWidth}px`
        const dragOffset = scrollOffset + 1

        return columnVirtualizer.getVirtualItems().map((virtualColumn) => {
            const column = columns.length <= virtualColumn.index ? {name: "", description: ""} : columns[virtualColumn.index]
            const width = `${virtualColumn.size}px`
            const transform = `translateX(${virtualColumn.start}px)`
            const dragTransform = `translateX(${virtualColumn.start + virtualColumn.size - dragOffset}px)`
            return (
                <Fragment key={virtualColumn.key}>
                    <Box sx={SX.cellFixed} borderRight={1} borderTop={1} borderBottom={2} style={{width, height: heightPx, transform}}>
                        <Tooltip title={renderHeadTitle(column.name, column.description)} placement={"top"} arrow={true}>
                            <Box sx={SX.noWrap}>{column.name}</Box>
                        </Tooltip>
                    </Box>
                    <Box
                        sx={SX.columnSeparator}
                        style={{height: heightPx, width: dragWidthPx, transform: dragTransform}}
                        onMouseDown={(e) => handleMouseDown(e, virtualColumn.index, virtualColumn.size)}
                    />
                </Fragment>

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

    function handleEffectDragging() {
        document.addEventListener("mousemove", callbackMouseMove)
        document.addEventListener("mouseup", callbackMouseUp)

        return () => {
            document.removeEventListener("mousemove", callbackMouseMove)
            document.addEventListener("mouseup", callbackMouseUp)
        }
    }

    function handleMouseMove(e: MouseEvent) {
        if (columnDrag.index != -1) {
            const offset = e.pageX - columnDrag.ps
            const width = columnDrag.size + offset
            if (width > 50) columnVirtualizer.resizeItem(columnDrag.index, width)
        }
    }

    function handleMouseUp() {
        document.body.style.removeProperty("cursor")
        setColumnDrag({index: -1, size: 0, ps: 0})
    }

    function handleMouseDown(e: React.MouseEvent, index: number, size: number) {
        e.preventDefault()
        e.stopPropagation()
        document.body.style.cursor = "col-resize"
        setColumnDrag({index, size, ps: e.pageX})
    }
}
