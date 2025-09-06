import {useRef} from "react";
import {useVirtualizer} from "@tanstack/react-virtual";
import {Box, SxProps} from "@mui/material";
import {Theme} from "@mui/material/styles";
import {AutoScrolling} from "./AutoScrolling";
import {SxPropsMap} from "../../../api/management/type";
import {SxPropsFormatter} from "../../../app/utils";

const SX: SxPropsMap = {
    container: {width: "100%", overflow: "auto", contain: "strict"},
    boxAbsolute: {position: "absolute", top: 0, left: 0, width: "100%"},
    boxRelative: {width: "100%", position: "relative"},
}

type Props = {
    height: number;
    auto: boolean;
    rows: string[],
    sx?: SxProps<Theme>,
    className?: string,
    sxVirtualRow?: SxProps<Theme>,
    classNameVirtualRow?: string,
}

/**
 *  This Component uses `@tanstack/react-virtual` to render only visible elements.
 *  Dynamic means that the height of the element is unknown before render
 *  Guide https://tanstack.com/virtual/v3/docs/examples/react/dynamic
 */
export function DynamicRowVirtualizer(props: Props) {
    const {rows, height, sx, className, sxVirtualRow, classNameVirtualRow, auto} = props
    const parentRef = useRef<Element>(null)

    const virtualizer = useVirtualizer({
        count: rows.length,
        getScrollElement: () => parentRef.current,
        estimateSize: () => 25,
    })

    const items = virtualizer.getVirtualItems()
    return (
        <AutoScrolling auto={auto} length={rows.length} scroll={virtualizer.scrollToIndex}>
            <Box ref={parentRef} sx={SxPropsFormatter.merge(sx, SX.container)} className={className} style={{height: `${height}px`}}>
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
            </Box>
        </AutoScrolling>
    )
}
