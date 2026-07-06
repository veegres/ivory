import {act, fireEvent, renderHook} from "@testing-library/react"
import {describe, expect, it, vi} from "vitest"

import {useDragger} from "./Dragger"

function mouseDownEvent(pageX: number) {
    return {preventDefault: vi.fn(), stopPropagation: vi.fn(), pageX} as unknown as React.MouseEvent
}

describe("useDragger", () => {
    it("should not call onMove without a preceding mouse down", () => {
        const onMove = vi.fn()
        renderHook(() => useDragger(10, onMove))

        fireEvent.mouseMove(document, {clientX: 200})

        expect(onMove).not.toHaveBeenCalled()
    })

    it("should call onMove with the new width once dragging exceeds minSize", () => {
        const onMove = vi.fn()
        const {result} = renderHook(() => useDragger(10, onMove))

        act(() => result.current.onMouseDown(mouseDownEvent(100), 2, 50))
        fireEvent.mouseMove(document, {clientX: 130})

        expect(onMove).toHaveBeenCalledWith(2, 80)
    })

    it("should not call onMove when resulting width does not exceed minSize", () => {
        const onMove = vi.fn()
        const {result} = renderHook(() => useDragger(100, onMove))

        act(() => result.current.onMouseDown(mouseDownEvent(100), 0, 50))
        fireEvent.mouseMove(document, {clientX: 105})

        expect(onMove).not.toHaveBeenCalled()
    })

    it("should stop tracking movement after mouse up", () => {
        const onMove = vi.fn()
        const {result} = renderHook(() => useDragger(10, onMove))

        act(() => result.current.onMouseDown(mouseDownEvent(100), 0, 50))
        fireEvent.mouseUp(document)
        fireEvent.mouseMove(document, {clientX: 200})

        expect(onMove).not.toHaveBeenCalled()
    })

    it("should prevent default and stop propagation on mouse down", () => {
        const onMove = vi.fn()
        const {result} = renderHook(() => useDragger(10, onMove))
        const event = mouseDownEvent(100)

        act(() => result.current.onMouseDown(event, 0, 50))

        expect(event.preventDefault).toHaveBeenCalled()
        expect(event.stopPropagation).toHaveBeenCalled()
    })

    it("should stop listening for movement after unmount", () => {
        const onMove = vi.fn()
        const {result, unmount} = renderHook(() => useDragger(10, onMove))

        act(() => result.current.onMouseDown(mouseDownEvent(100), 0, 50))
        unmount()
        fireEvent.mouseMove(document, {clientX: 200})

        expect(onMove).not.toHaveBeenCalled()
    })
})
