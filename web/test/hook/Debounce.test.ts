import {renderHook} from "@testing-library/react"
import {afterEach, beforeEach, describe, expect, it, vi} from "vitest"

import {useDebounce, useDebounceFunction} from "../../src/hook/Debounce"

describe("useDebounce", () => {
    afterEach(() => {
        vi.clearAllTimers()
    })

    it("should return initial value immediately", () => {
        const {result} = renderHook(() => useDebounce("initial", 500))

        expect(result.current).toBe("initial")
    })

    it("should return initial value when no changes occur", () => {
        const {result} = renderHook(() => useDebounce("test-value", 500))

        expect(result.current).toBe("test-value")
    })

    it("should work with object values", () => {
        const initialValue = {count: 0}
        const {result} = renderHook(() => useDebounce(initialValue, 500))

        expect(result.current).toEqual({count: 0})
    })

    it("should work with number values", () => {
        const {result} = renderHook(() => useDebounce(42, 500))

        expect(result.current).toBe(42)
    })

    it("should work with boolean values", () => {
        const {result} = renderHook(() => useDebounce(true, 500))

        expect(result.current).toBe(true)
    })

    it("should accept custom delay", () => {
        const {result} = renderHook(() => useDebounce("value", 1000))

        expect(result.current).toBe("value")
    })
})

describe("useDebounceFunction", () => {
    beforeEach(() => {
        vi.useFakeTimers()
    })

    afterEach(() => {
        vi.useRealTimers()
    })

    it("should call onChange after delay", () => {
        const onChange = vi.fn()

        renderHook(() => useDebounceFunction("test", onChange, 500))

        // Should not call immediately
        expect(onChange).not.toHaveBeenCalled()

        // Should call after delay
        vi.advanceTimersByTime(500)
        expect(onChange).toHaveBeenCalledWith("test")
        expect(onChange).toHaveBeenCalledTimes(1)
    })

    it("should use default delay of 500ms", () => {
        const onChange = vi.fn()

        renderHook(() => useDebounceFunction("test", onChange))

        vi.advanceTimersByTime(499)
        expect(onChange).not.toHaveBeenCalled()

        vi.advanceTimersByTime(1)
        expect(onChange).toHaveBeenCalledWith("test")
    })

    it("should cancel previous timeout on value change", () => {
        const onChange = vi.fn()
        const {rerender} = renderHook(
            ({value}) => useDebounceFunction(value, onChange, 500),
            {initialProps: {value: "initial"}}
        )

        vi.advanceTimersByTime(200)
        expect(onChange).not.toHaveBeenCalled()

        // Change value before timeout completes
        rerender({value: "updated"})

        // Complete original timeout duration
        vi.advanceTimersByTime(300)
        expect(onChange).not.toHaveBeenCalled()

        // Complete new timeout
        vi.advanceTimersByTime(200)
        expect(onChange).toHaveBeenCalledWith("updated")
        expect(onChange).toHaveBeenCalledTimes(1)
    })

    it("should handle multiple rapid changes", () => {
        const onChange = vi.fn()
        const {rerender} = renderHook(
            ({value}) => useDebounceFunction(value, onChange, 500),
            {initialProps: {value: "first"}}
        )

        vi.advanceTimersByTime(100)
        rerender({value: "second"})

        vi.advanceTimersByTime(100)
        rerender({value: "third"})

        vi.advanceTimersByTime(100)
        rerender({value: "fourth"})

        // Should not have called onChange yet
        expect(onChange).not.toHaveBeenCalled()

        // Complete timeout from last change
        vi.advanceTimersByTime(500)

        // Should only call with the last value
        expect(onChange).toHaveBeenCalledWith("fourth")
        expect(onChange).toHaveBeenCalledTimes(1)
    })

    it("should work with different value types", () => {
        const onChange = vi.fn()

        const {rerender} = renderHook(
            ({value}) => useDebounceFunction(value, onChange, 500),
            {initialProps: {value: {id: 1}}}
        )

        vi.advanceTimersByTime(500)
        expect(onChange).toHaveBeenCalledWith({id: 1})

        onChange.mockClear()
        rerender({value: {id: 2}})
        vi.advanceTimersByTime(500)
        expect(onChange).toHaveBeenCalledWith({id: 2})
    })

    it("should handle onChange function changes", () => {
        const onChange1 = vi.fn()
        const onChange2 = vi.fn()

        const {rerender} = renderHook(
            ({callback}) => useDebounceFunction("test", callback, 500),
            {initialProps: {callback: onChange1}}
        )

        vi.advanceTimersByTime(500)
        expect(onChange1).toHaveBeenCalledWith("test")
        expect(onChange2).not.toHaveBeenCalled()

        onChange1.mockClear()
        rerender({callback: onChange2})

        vi.advanceTimersByTime(500)
        expect(onChange1).not.toHaveBeenCalled()
        expect(onChange2).toHaveBeenCalledWith("test")
    })
})
