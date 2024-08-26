import scroll from "../../../style/scroll.module.css";
import {Box, CircularProgress, Tooltip} from "@mui/material";
import {SxPropsMap} from "../../../type/general";
import {Fragment, ReactNode, useMemo, useState} from "react";
import {useVirtualizer} from "@tanstack/react-virtual";
import {NoBox} from "../box/NoBox";
import {useDragger} from "../../../hook/Dragger";
import {MenuButton} from "../button/MenuButton";

const SX: SxPropsMap = {
    box: {position: "relative"},
    body: {overflow: "auto", width: "100%", height: "100%", fontSize: "12px"},
    loader: {
        position: "absolute", top: 0, left: 0, right: 0, bottom: 0, zIndex: 5,
        display: "flex", alignItems: "center", justifyContent: "center",
        width: "inherit", height: "inherit", bgcolor: "rgba(31,44,126,0.05)",
    },
    cell: {
        position: "absolute", top: 0, left: 0, zIndex: 0, padding: "5px 10px",
        bgcolor: "background.paper", borderBottom: 1, borderRight: 1, borderColor: "divider",
        "&:hover": {
            overflow: "auto", zIndex: 1,
            width: "auto!important", height: "auto!important",
            maxHeight: "150px", maxWidth: "400px",
        },
        "&:hover div": {
            whiteSpace: "pre", overflow: "unset", textOverflow: "unset",
        }
    },
    noWrap: {whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis", width: "100%"},
    preWrap: {whiteSpace: "pre-wrap", textWrap: "nowrap", overflow: "hidden", textOverflow: "ellipsis", width: "100%"},
    cellFixed: {
        position: "absolute", top: 0, left: 0, padding: "5px 10px",
        display: "flex", alignItems: "center", justifyContent: "center",
        color: "text.disabled", bgcolor: "background.paper", borderColor: "divider",
    },
    columnSeparator: {
        position: "absolute", top: 0, left: 0, zIndex: 4,
        cursor: "col-resize", touchAction: "none",
        display: "flex", alignItems: "center", justifyContent: "center",
        "&:hover:after": {content: `""`, width: "30%", height: "60%", bgcolor: "gray"},
    },
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
    const [ref, setRef] = useState<Element | null>(null);

    const columnCount = columns.length
    const scrollSize = 6

    const cellWidthSticky = 50
    const cellHeightSticky = 32
    const cellWidthStickyPx = `${cellWidthSticky}px`
    const cellHeightStickyPx = `${cellHeightSticky}px`

    const cellWidth = 100
    const cellHeight = 29
    // NOTE: calculate the cell width, the component should be rendered from scratch to change the size,
    //  this is how useVirtualizer work (-1 is needed to always choose 100 if width is unknown)
    const cellCalcWidth = useMemo(() => Math.max(Math.floor(((ref?.clientWidth ?? -1) - cellWidthSticky - (scrollSize * 2)) / columnCount), cellWidth), [columnCount, ref])

    const rowVirtualizer = useVirtualizer({
        count: rows.length,
        getScrollElement: () => ref,
        estimateSize: () => cellHeight,
        scrollPaddingStart: cellHeightSticky,
        overscan: 3,
        enabled: ref !== null,
    })
    const columnVirtualizer = useVirtualizer({
        horizontal: true,
        count: columnCount,
        getScrollElement: () => ref,
        estimateSize: () => cellCalcWidth,
        scrollPaddingStart: cellWidthSticky,
        overscan: 3,
        // NOTE: this helps to use cellWidth, because in first render when we don't have ref, it
        //  renders incorrect size
        enabled: ref !== null,
    })
    const columnDragger = useDragger(columnVirtualizer.resizeItem)

    const bodyWidthPx = `${columnVirtualizer.getTotalSize()}px`
    const bodyHeightPx = `${rowVirtualizer.getTotalSize()}px`

    return (
        <Box sx={SX.box} width={width} height={height}>
            {renderLoader()}
            {columns.length > 0 ? renderTable() : renderEmpty()}
        </Box>
    )

    function renderLoader() {
        if (!fetching) return
        return (
            <Box sx={SX.loader}>
                <CircularProgress/>
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
                ref={setRef}
                sx={SX.body}
                padding={`0px ${scrollSize}px ${scrollSize}px 0px`}
                display={"grid"}
                gridTemplateColumns={`${cellWidthStickyPx} ${bodyWidthPx} auto`}
                gridTemplateRows={`${cellHeightStickyPx} ${bodyHeightPx}`}
                className={scroll.small}
            >
                <Box position={"sticky"} zIndex={3} top={0} left={0}>
                    {renderCorner(1, 2)}
                </Box>
                <Box position={"sticky"} zIndex={2} top={0}>
                    {renderHead()}
                </Box>
                <Box position={"relative"} zIndex={3}>
                    {renderRowCell && renderCorner(2, 1)}
                </Box>
                <Box position={"sticky"} zIndex={2} left={0}>
                    {renderIndex()}
                </Box>
                <Box position={"relative"}>
                    {renderBody()}
                </Box>
                <Box position={"relative"} zIndex={3}>
                    {renderActions()}
                </Box>
            </Box>
        )
    }

    function renderCorner(left: number, right: number) {
        return (
            <Box sx={SX.cellFixed} borderTop={1} borderBottom={2} borderLeft={left} borderRight={right}
                 style={{width: cellWidthStickyPx, height: cellHeightStickyPx}}/>
        )
    }

    function renderHead() {
        const dragWidth = scrollSize * 2
        const dragWidthPx = `${dragWidth}px`
        const dragOffset = scrollSize + 1

        return columnVirtualizer.getVirtualItems().map((virtualColumn) => {
            const column = columns[virtualColumn.index]
            const width = `${virtualColumn.size}px`
            const transform = `translateX(${virtualColumn.start}px)`
            const dragTransform = `translateX(${virtualColumn.start + virtualColumn.size - dragOffset}px)`
            return (
                <Fragment key={virtualColumn.key}>
                    <Box sx={SX.cellFixed} borderRight={1} borderTop={1} borderBottom={2}
                         style={{width, height: cellHeightStickyPx, transform}}>
                        <Tooltip title={renderHeadTitle(column.name, column.description)} placement={"top"}>
                            <Box sx={SX.noWrap}>{column.name}</Box>
                        </Tooltip>
                    </Box>
                    <Box
                        sx={SX.columnSeparator}
                        style={{width: dragWidthPx, height: cellHeightStickyPx, transform: dragTransform}}
                        onMouseDown={(e) => columnDragger.onMouseDown(e, virtualColumn.index, virtualColumn.size)}
                    />
                </Fragment>
            )
        })
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
        const width = `${cellWidthSticky}px`
        const height = `${cellHeight}px`

        return rowVirtualizer.getVirtualItems().map((virtualRow) => {
            const transform = `translateY(${virtualRow.start}px)`
            return (
                <Box key={virtualRow.key} sx={SX.cellFixed} borderLeft={1} borderBottom={1} borderRight={2}
                     style={{width, height, transform}}>
                    {virtualRow.index + 1}
                </Box>
            )
        })
    }

    function renderActions() {
        if (!renderRowCell) return
        const width = `${cellWidthSticky}px`
        const height = `${cellHeight}px`

        return rowVirtualizer.getVirtualItems().map((virtualRow) => {
            const row = rows[virtualRow.index]
            const transform = `translateY(${virtualRow.start}px)`
            return (
                <Box key={virtualRow.key} sx={SX.cellFixed} borderLeft={2} borderBottom={1} borderRight={1}
                     style={{width, height, transform}}>
                    <MenuButton size={"small"}>{renderRowCell(row)}</MenuButton>
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
                        const value = getParsedCell(row[virtualColumn.index])
                        const height = `${virtualRow.size}px`
                        const width = `${virtualColumn.size}px`
                        const transform = `translateX(${virtualColumn.start}px) translateY(${virtualRow.start}px)`
                        return (
                            <Box
                                key={virtualColumn.key}
                                sx={SX.cell}
                                className={scroll.tiny}
                                style={{minWidth: width, minHeight: height, width, height, transform}}
                            >
                                <Box sx={wrap}>{value}</Box>
                            </Box>
                        )
                    })}
                </Fragment>
            )
        })
    }

    function getParsedCell(cell: any): ReactNode {
        if (typeof cell === 'object') return JSON.stringify(cell)
        else if (typeof cell === "boolean") return Boolean(cell).toString()
        else return cell
    }
}
