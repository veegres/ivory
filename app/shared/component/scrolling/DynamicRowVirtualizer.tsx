import {Box, SxProps} from "@mui/material"
import {Theme} from "@mui/material/styles"
import {useVirtualizer} from "@tanstack/react-virtual"
import {ReactNode, useRef} from "react"

import {SxPropsMap} from "../../helper/type"
import {printLogs, SxPropsFormatter} from "../../helper/utils"
import {AutoScrolling} from "./AutoScrolling"

const SX: SxPropsMap = {
    container: {width: "100%", overflow: "auto", contain: "strict"},
    boxAbsolute: {position: "absolute", top: 0, left: 0, width: "100%"},
    boxRelative: {width: "100%", position: "relative"},
}

type Props = {
    height: number,
    auto: boolean,
    rows: string[],
    sx?: SxProps<Theme>,
    className?: string,
    sxVirtualRow?: SxProps<Theme>,
    classNameVirtualRow?: string,
    reconnect?: () => void,
    empty?: ReactNode,
}

/**
 *  This Component uses `@tanstack/react-virtual` to render only visible elements.
 *  Dynamic means that the height of the element is unknown before render
 *  Guide https://tanstack.com/virtual/v3/docs/examples/react/dynamic
 */
export function DynamicRowVirtualizer(props: Props) {
    const {rows, height, sx, className, sxVirtualRow, classNameVirtualRow, auto, reconnect, empty} = props
    const parentRef = useRef<Element>(null)

    // eslint-disable-next-line react-hooks/incompatible-library
    const virtualizer = useVirtualizer({
        count: rows.length,
        getScrollElement: () => parentRef.current,
        estimateSize: () => 25,
    })

    const items = virtualizer.getVirtualItems()
    return (
        // NOTE: AutoScrolling (and its buttons, e.g. print/reconnect) stays mounted regardless
        //  of whether rows is currently empty, so it does not flicker/remount whenever a stream
        //  reconnect briefly clears the rows before new ones arrive. print is always passed
        //  (even for an empty array) rather than toggled to undefined, for the same reason.
        <AutoScrolling
            auto={auto}
            length={rows.length}
            scroll={virtualizer.scrollToIndex}
            print={() => printLogs(`LOGS: ${new Date()}`, rows)}
            reconnect={reconnect}
        >
            <Box ref={parentRef} sx={SxPropsFormatter.merge(sx, SX.container)} className={className} style={{height: `${height}px`}}>
                {rows.length === 0 && empty ? empty : (
                    <Box sx={SX.boxRelative} style={{height: virtualizer.getTotalSize()}}>
                        <Box sx={SX.boxAbsolute} style={{transform: `translateY(${items[0]?.start ?? 0}px)`}}>
                            {items.map((virtualRow) => (
                                <Box
                                    ref={virtualizer.measureElement}
                                    key={virtualRow.key}
                                    data-index={virtualRow.index}
                                >
                                    <Box sx={sxVirtualRow} className={classNameVirtualRow}>
                                        {rows[virtualRow.index]}
                                    </Box>
                                </Box>
                            ))}
                        </Box>
                    </Box>
                )}
            </Box>
        </AutoScrolling>
    )
}
