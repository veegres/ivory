import React, {useCallback, useEffect, useState} from "react"

const init = {index: -1, size: 0, ps: 0}

export function useDragger(minSize: number, onMove: (index: number, size: number) => void) {
    const [columnDrag, setColumnDrag] = useState(init)

    const callbackMouseMove = useCallback(handleMouseMove, [columnDrag.index, columnDrag.ps, columnDrag.size, minSize, onMove])
    const callbackMouseUp = useCallback(handleMouseUp, [])
    useEffect(handleEffectDragging, [callbackMouseMove, callbackMouseUp])

    return {
        onMouseDown: handleMouseDown,
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
            if (width > minSize) onMove(columnDrag.index, width)
        }
    }

    function handleMouseUp() {
        document.body.style.removeProperty("cursor")
        setColumnDrag(init)
    }

    function handleMouseDown(e: React.MouseEvent, index: number, size: number) {
        e.preventDefault()
        e.stopPropagation()
        document.body.style.cursor = "col-resize"
        setColumnDrag({index, size, ps: e.pageX})
    }
}
