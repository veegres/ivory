import {renderHook} from "@testing-library/react"
import {beforeEach, describe, expect, it, vi} from "vitest"

import {useWindowChildCount, useWindowScrolled, useWindowSize} from "../../src/hook/WindowObservers"

// Mock ResizeObserver
class ResizeObserverMock {
    callback: ResizeObserverCallback
    constructor(callback: ResizeObserverCallback) {
        this.callback = callback
    }
    observe() {
        // Immediately call the callback to simulate resize
        this.callback([] as any, this as any)
    }
    disconnect() {}
    unobserve() {}
}

// Mock MutationObserver
class MutationObserverMock {
    callback: MutationCallback
    constructor(callback: MutationCallback) {
        this.callback = callback
    }
    observe() {
        // Immediately call the callback to simulate mutation
        this.callback([] as any, this as any)
    }
    disconnect() {}
    takeRecords() {
        return []
    }
}

describe("useWindowSize", () => {
    beforeEach(() => {
        (globalThis as any).ResizeObserver = ResizeObserverMock as any
    })

    it("should return undefined for both dimensions when no element provided", () => {
        const {result} = renderHook(() => useWindowSize())

        expect(result.current[0]).toBeUndefined()
        expect(result.current[1]).toBeUndefined()
    })

    it("should return element dimensions", () => {
        const element = document.createElement("div")
        Object.defineProperty(element, "clientWidth", {value: 800, configurable: true})
        Object.defineProperty(element, "clientHeight", {value: 600, configurable: true})

        const {result} = renderHook(() => useWindowSize(element))

        expect(result.current[0]).toBe(800)
        expect(result.current[1]).toBe(600)
    })

    it("should update dimensions when element changes", () => {
        const element1 = document.createElement("div")
        Object.defineProperty(element1, "clientWidth", {value: 800, configurable: true})
        Object.defineProperty(element1, "clientHeight", {value: 600, configurable: true})

        const {result, rerender} = renderHook(
            ({el}) => useWindowSize(el),
            {initialProps: {el: element1}}
        )

        expect(result.current[0]).toBe(800)
        expect(result.current[1]).toBe(600)

        const element2 = document.createElement("div")
        Object.defineProperty(element2, "clientWidth", {value: 1024, configurable: true})
        Object.defineProperty(element2, "clientHeight", {value: 768, configurable: true})

        rerender({el: element2})

        expect(result.current[0]).toBe(1024)
        expect(result.current[1]).toBe(768)
    })

    it("should handle element with zero dimensions", () => {
        const element = document.createElement("div")
        Object.defineProperty(element, "clientWidth", {value: 0, configurable: true})
        Object.defineProperty(element, "clientHeight", {value: 0, configurable: true})

        const {result} = renderHook(() => useWindowSize(element))

        expect(result.current[0]).toBe(0)
        expect(result.current[1]).toBe(0)
    })

    it("should cleanup observer on unmount", () => {
        const disconnectSpy = vi.spyOn(ResizeObserverMock.prototype, "disconnect")

        const element = document.createElement("div")
        const {unmount} = renderHook(() => useWindowSize(element))

        unmount()

        expect(disconnectSpy).toHaveBeenCalled()
    })
})

describe("useWindowChildCount", () => {
    beforeEach(() => {
        (globalThis as any).MutationObserver = MutationObserverMock as any
    })

    it("should return undefined when no element provided", () => {
        const {result} = renderHook(() => useWindowChildCount())

        expect(result.current[0]).toBeUndefined()
    })

    it("should return element child count", () => {
        const element = document.createElement("div")
        element.appendChild(document.createElement("span"))
        element.appendChild(document.createElement("span"))
        element.appendChild(document.createElement("span"))

        const {result} = renderHook(() => useWindowChildCount(element))

        expect(result.current[0]).toBe(3)
    })

    it("should handle element with no children", () => {
        const element = document.createElement("div")

        const {result} = renderHook(() => useWindowChildCount(element))

        expect(result.current[0]).toBe(0)
    })

    it("should update count when element changes", () => {
        const element1 = document.createElement("div")
        element1.appendChild(document.createElement("span"))
        element1.appendChild(document.createElement("span"))

        const {result, rerender} = renderHook(
            ({el}) => useWindowChildCount(el),
            {initialProps: {el: element1}}
        )

        expect(result.current[0]).toBe(2)

        const element2 = document.createElement("div")
        element2.appendChild(document.createElement("span"))
        element2.appendChild(document.createElement("span"))
        element2.appendChild(document.createElement("span"))
        element2.appendChild(document.createElement("span"))

        rerender({el: element2})

        expect(result.current[0]).toBe(4)
    })

    it("should cleanup observer on unmount", () => {
        const disconnectSpy = vi.spyOn(MutationObserverMock.prototype, "disconnect")

        const element = document.createElement("div")
        const {unmount} = renderHook(() => useWindowChildCount(element))

        unmount()

        expect(disconnectSpy).toHaveBeenCalled()
    })
})

describe("useWindowScrolled", () => {
    beforeEach(() => {
        (globalThis as any).ResizeObserver = ResizeObserverMock as any;
        (globalThis as any).MutationObserver = MutationObserverMock as any
    })

    it("should return false when no element provided", () => {
        const {result} = renderHook(() => useWindowScrolled())

        expect(result.current[0]).toBe(false)
    })

    it("should return false when element has no horizontal scroll", () => {
        const element = document.createElement("div")
        Object.defineProperty(element, "scrollWidth", {value: 800, configurable: true})
        Object.defineProperty(element, "clientWidth", {value: 800, configurable: true})
        Object.defineProperty(element, "clientHeight", {value: 600, configurable: true})

        const {result} = renderHook(() => useWindowScrolled(element))

        expect(result.current[0]).toBe(false)
    })

    it("should return true when element has horizontal scroll", () => {
        const element = document.createElement("div")
        Object.defineProperty(element, "scrollWidth", {value: 1000, configurable: true})
        Object.defineProperty(element, "clientWidth", {value: 800, configurable: true})
        Object.defineProperty(element, "clientHeight", {value: 600, configurable: true})

        const {result} = renderHook(() => useWindowScrolled(element))

        expect(result.current[0]).toBe(true)
    })

    it("should update when element dimensions change", () => {
        const element1 = document.createElement("div")
        Object.defineProperty(element1, "scrollWidth", {value: 800, configurable: true})
        Object.defineProperty(element1, "clientWidth", {value: 800, configurable: true})
        Object.defineProperty(element1, "clientHeight", {value: 600, configurable: true})

        const {result, rerender} = renderHook(
            ({el}) => useWindowScrolled(el),
            {initialProps: {el: element1}}
        )

        expect(result.current[0]).toBe(false)

        // Create a new element with scroll
        const element2 = document.createElement("div")
        Object.defineProperty(element2, "scrollWidth", {value: 1200, configurable: true})
        Object.defineProperty(element2, "clientWidth", {value: 800, configurable: true})
        Object.defineProperty(element2, "clientHeight", {value: 600, configurable: true})

        rerender({el: element2})

        expect(result.current[0]).toBe(true)
    })

    it("should handle edge case where scrollWidth equals clientWidth", () => {
        const element = document.createElement("div")
        Object.defineProperty(element, "scrollWidth", {value: 800, configurable: true})
        Object.defineProperty(element, "clientWidth", {value: 800, configurable: true})
        Object.defineProperty(element, "clientHeight", {value: 600, configurable: true})

        const {result} = renderHook(() => useWindowScrolled(element))

        expect(result.current[0]).toBe(false)
    })

    it("should handle edge case where scrollWidth is 1px more than clientWidth", () => {
        const element = document.createElement("div")
        Object.defineProperty(element, "scrollWidth", {value: 801, configurable: true})
        Object.defineProperty(element, "clientWidth", {value: 800, configurable: true})
        Object.defineProperty(element, "clientHeight", {value: 600, configurable: true})

        const {result} = renderHook(() => useWindowScrolled(element))

        expect(result.current[0]).toBe(true)
    })
})
