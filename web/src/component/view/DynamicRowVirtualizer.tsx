import React, { useRef } from "react";
import {useVirtualizer} from "@tanstack/react-virtual";
import {Box} from "@mui/material";
import {SxProps} from "@mui/system";
import {Theme} from "@mui/material/styles";

const SX = {
    container: { width: "100%", overflow: 'auto' },
    boxSize: { width: "100%", position: "relative" },
    virtualRow: { position: "absolute", top: 0, left: 0, width: "100%" }
}

type Props = {
    height: number;
    rows: string[],
    sx?: SxProps<Theme>,
    className?: string,
    sxVirtualRow?: SxProps<Theme>,
    classNameVirtualRow?: string,
}

export function DynamicRowVirtualizer({ rows, height, sx, className, sxVirtualRow, classNameVirtualRow }: Props) {
    const parentRef = useRef<Element>(null)

    const rowVirtualizer = useVirtualizer({
        count: rows.length,
        getScrollElement: () => parentRef.current,
        estimateSize: () => 125,
    })

    return (
        <Box ref={parentRef} sx={{...sx, ...SX.container}} className={className} style={{ height: `${height}px` }}>
            <Box sx={SX.boxSize} style={{height: rowVirtualizer.getTotalSize()}}>
                {rowVirtualizer.getVirtualItems().map((virtualRow) => (
                    <Box
                        ref={parentRef}
                        key={virtualRow.index}
                        sx={SX.virtualRow}
                        style={{transform: `translateY(${virtualRow.start}px)`}}
                    >
                        <Box sx={sxVirtualRow} className={classNameVirtualRow} style={{ height: rows[virtualRow.index] }}>
                            {rows[virtualRow.index]}
                        </Box>
                    </Box>
                ))}
            </Box>
        </Box>
    )
}
