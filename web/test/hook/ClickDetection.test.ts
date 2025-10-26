import {fireEvent, renderHook} from "@testing-library/react"
import {createRef} from "react"
import {describe, expect, it} from "vitest"

import {useClickDetection} from "../../src/hook/ClickDetection"

describe("useClickDetection", () => {
    it("should return undefined initially", () => {
        const ref = createRef<HTMLDivElement>()
        const {result} = renderHook(() => useClickDetection(ref))

        expect(result.current).toBeUndefined()
    })

    it("should detect click inside element", () => {
        const element = document.createElement("div")
        document.body.appendChild(element)

        const ref = {current: element}
        const {result} = renderHook(() => useClickDetection(ref))

        // Click inside the element
        fireEvent.mouseDown(element)

        expect(result.current).toBe("inside")

        document.body.removeChild(element)
    })

    it("should detect click outside element", () => {
        const element = document.createElement("div")
        document.body.appendChild(element)

        const ref = {current: element}
        const {result} = renderHook(() => useClickDetection(ref))

        // Click outside the element (on document body)
        fireEvent.mouseDown(document.body)

        expect(result.current).toBe("outside")

        document.body.removeChild(element)
    })

    it("should detect clicks on child elements as inside", () => {
        const parent = document.createElement("div")
        const child = document.createElement("button")
        parent.appendChild(child)
        document.body.appendChild(parent)

        const ref = {current: parent}
        const {result} = renderHook(() => useClickDetection(ref))

        // Click on child element
        fireEvent.mouseDown(child)

        expect(result.current).toBe("inside")

        document.body.removeChild(parent)
    })

    it("should update click state on multiple clicks", () => {
        const element = document.createElement("div")
        document.body.appendChild(element)

        const ref = {current: element}
        const {result} = renderHook(() => useClickDetection(ref))

        // First click inside
        fireEvent.mouseDown(element)
        expect(result.current).toBe("inside")

        // Then click outside
        fireEvent.mouseDown(document.body)
        expect(result.current).toBe("outside")

        // Click inside again
        fireEvent.mouseDown(element)
        expect(result.current).toBe("inside")

        document.body.removeChild(element)
    })

    it("should handle null ref", () => {
        const ref = {current: null}
        const {result} = renderHook(() => useClickDetection(ref))

        // Click anywhere
        fireEvent.mouseDown(document.body)

        expect(result.current).toBe("outside")
    })

    it("should cleanup event listener on unmount", () => {
        const element = document.createElement("div")
        document.body.appendChild(element)

        const ref = {current: element}
        const {result, unmount} = renderHook(() => useClickDetection(ref))

        // Click inside
        fireEvent.mouseDown(element)
        expect(result.current).toBe("inside")

        // Unmount the hook
        unmount()

        // Click again - should not update state
        const previousValue = result.current
        fireEvent.mouseDown(document.body)

        // State should remain the same after unmount
        expect(result.current).toBe(previousValue)

        document.body.removeChild(element)
    })

    it("should update event listener when ref changes", () => {
        const element1 = document.createElement("div")
        const element2 = document.createElement("div")
        document.body.appendChild(element1)
        document.body.appendChild(element2)

        const ref1 = {current: element1}
        const {result, rerender} = renderHook(
            ({ref}) => useClickDetection(ref),
            {initialProps: {ref: ref1}}
        )

        // Click on first element
        fireEvent.mouseDown(element1)
        expect(result.current).toBe("inside")

        // Change ref to second element
        const ref2 = {current: element2}
        rerender({ref: ref2})

        // Click on first element should now be outside
        fireEvent.mouseDown(element1)
        expect(result.current).toBe("outside")

        // Click on second element should be inside
        fireEvent.mouseDown(element2)
        expect(result.current).toBe("inside")

        document.body.removeChild(element1)
        document.body.removeChild(element2)
    })
})
