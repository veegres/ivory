import {renderHook, waitFor} from "@testing-library/react"
import {afterEach, beforeEach, describe, expect, it, vi} from "vitest"

import {useLocalStorageState} from "../../src/hook/LocalStorage"

describe("useLocalStorageState", () => {
    beforeEach(() => {
        localStorage.clear()
        vi.clearAllMocks()
    })

    afterEach(() => {
        localStorage.clear()
    })

    describe("initialization", () => {
        it("should return initial value when localStorage is empty", () => {
            const {result} = renderHook(() => useLocalStorageState("test-key", "initial"))

            expect(result.current[0]).toBe("initial")
        })

        it("should return value from localStorage when it exists", () => {
            localStorage.setItem("test-key", "stored")

            const {result} = renderHook(() => useLocalStorageState("test-key", "initial"))

            expect(result.current[0]).toBe("stored")
        })

        it("should handle initial value from function", () => {
            const initialFn = () => "computed"
            const {result} = renderHook(() => useLocalStorageState("test-key", initialFn))

            expect(result.current[0]).toBe("computed")
        })

        it("should parse JSON for object values", () => {
            const storedObject = {name: "test", count: 42}
            localStorage.setItem("test-key", JSON.stringify(storedObject))

            const {result} = renderHook(() => useLocalStorageState("test-key", {}))

            expect(result.current[0]).toEqual(storedObject)
        })

        it("should parse JSON for array values", () => {
            const storedArray = [1, 2, 3]
            localStorage.setItem("test-key", JSON.stringify(storedArray))

            const {result} = renderHook(() => useLocalStorageState("test-key", []))

            expect(result.current[0]).toEqual(storedArray)
        })
    })

    describe("setValue", () => {
        it("should update state and localStorage with string value", async () => {
            const {result} = renderHook(() => useLocalStorageState("test-key", "initial"))

            const [, setValue] = result.current
            setValue("updated")

            await waitFor(() => {
                expect(result.current[0]).toBe("updated")
                expect(localStorage.getItem("test-key")).toBe("updated")
            })
        })

        it("should update state and localStorage with number value", async () => {
            const {result} = renderHook(() => useLocalStorageState("test-key", 0))

            const [, setValue] = result.current
            setValue(42)

            await waitFor(() => {
                expect(result.current[0]).toBe(42)
                expect(localStorage.getItem("test-key")).toBe("42")
            })
        })

        it("should serialize object to JSON in localStorage", async () => {
            const {result} = renderHook(() =>
                useLocalStorageState("test-key", {count: 0})
            )

            const [, setValue] = result.current
            setValue({count: 5})

            await waitFor(() => {
                expect(result.current[0]).toEqual({count: 5})
                expect(localStorage.getItem("test-key")).toBe(JSON.stringify({count: 5}))
            })
        })

        it("should handle functional updates", async () => {
            const {result} = renderHook(() => useLocalStorageState("test-key", 0))

            const [, setValue] = result.current
            setValue((prev) => prev + 1)

            await waitFor(() => {
                expect(result.current[0]).toBe(1)
            })

            setValue((prev) => prev + 1)

            await waitFor(() => {
                expect(result.current[0]).toBe(2)
            })
        })

        it("should not update localStorage with undefined value", () => {
            const {result} = renderHook(() => useLocalStorageState("test-key", "initial"))

            const [, setValue] = result.current
            setValue(undefined as any)

            // localStorage should not be updated
            expect(localStorage.getItem("test-key")).toBe("initial")
        })

        it("should not update localStorage with null value", () => {
            const {result} = renderHook(() => useLocalStorageState("test-key", "initial"))

            const [, setValue] = result.current
            setValue(null as any)

            // localStorage should not be updated
            expect(localStorage.getItem("test-key")).toBe("initial")
        })
    })

    describe("storage sync", () => {
        it("should not sync by default", () => {
            const {result} = renderHook(() => useLocalStorageState("test-key", "initial"))

            // Simulate storage event from another tab
            const storageEvent = new StorageEvent("storage", {
                key: "test-key",
                newValue: "from-other-tab",
            })
            window.dispatchEvent(storageEvent)

            // Value should remain unchanged since sync is false
            expect(result.current[0]).toBe("initial")
        })

        it("should sync when sync=true and storage event is fired", async () => {
            const {result} = renderHook(() =>
                useLocalStorageState("test-key", "initial", true)
            )

            // Simulate storage event from another tab
            const storageEvent = new StorageEvent("storage", {
                key: "test-key",
                newValue: "from-other-tab",
            })
            window.dispatchEvent(storageEvent)

            await waitFor(() => {
                expect(result.current[0]).toBe("from-other-tab")
            })
        })

        it("should sync object values", async () => {
            const {result} = renderHook(() =>
                useLocalStorageState("test-key", {count: 0}, true)
            )

            const storageEvent = new StorageEvent("storage", {
                key: "test-key",
                newValue: JSON.stringify({count: 10}),
            })
            window.dispatchEvent(storageEvent)

            await waitFor(() => {
                expect(result.current[0]).toEqual({count: 10})
            })
        })

        it("should reset to initial value when storage is cleared", async () => {
            const {result} = renderHook(() =>
                useLocalStorageState("test-key", "initial", true)
            )

            const storageEvent = new StorageEvent("storage", {
                key: "test-key",
                newValue: null,
            })
            window.dispatchEvent(storageEvent)

            await waitFor(() => {
                expect(result.current[0]).toBe("initial")
            })
        })

        it("should ignore storage events for different keys", async () => {
            const {result} = renderHook(() =>
                useLocalStorageState("test-key", "initial", true)
            )

            const storageEvent = new StorageEvent("storage", {
                key: "other-key",
                newValue: "other-value",
            })
            window.dispatchEvent(storageEvent)

            // Value should remain unchanged
            expect(result.current[0]).toBe("initial")
        })
    })

    describe("edge cases", () => {
        it("should handle boolean values", async () => {
            const {result} = renderHook(() => useLocalStorageState("test-key", false))

            const [, setValue] = result.current
            setValue(true)

            await waitFor(() => {
                expect(result.current[0]).toBe(true)
                expect(localStorage.getItem("test-key")).toBe("true")
            })
        })

        it("should handle empty string", async () => {
            const {result} = renderHook(() => useLocalStorageState("test-key", "initial"))

            const [, setValue] = result.current
            setValue("")

            await waitFor(() => {
                expect(result.current[0]).toBe("")
                expect(localStorage.getItem("test-key")).toBe("")
            })
        })

        it("should handle zero value", async () => {
            const {result} = renderHook(() => useLocalStorageState("test-key", 100))

            const [, setValue] = result.current
            setValue(0)

            await waitFor(() => {
                expect(result.current[0]).toBe(0)
                expect(localStorage.getItem("test-key")).toBe("0")
            })
        })

        it("should handle nested objects", async () => {
            const initialValue = {user: {name: "John", settings: {theme: "dark"}}}
            const {result} = renderHook(() => useLocalStorageState("test-key", initialValue))

            const [, setValue] = result.current
            const updatedValue = {user: {name: "Jane", settings: {theme: "light"}}}
            setValue(updatedValue)

            await waitFor(() => {
                expect(result.current[0]).toEqual(updatedValue)
                expect(JSON.parse(localStorage.getItem("test-key") || "")).toEqual(updatedValue)
            })
        })
    })

    describe("multiple instances", () => {
        it("should sync state between multiple hooks with same key", async () => {
            const {result: result1} = renderHook(() =>
                useLocalStorageState("shared-key", "initial")
            )
            const {result: result2} = renderHook(() =>
                useLocalStorageState("shared-key", "initial")
            )

            expect(result1.current[0]).toBe("initial")
            expect(result2.current[0]).toBe("initial")

            // Update from first hook
            const [, setValue1] = result1.current
            setValue1("updated")

            await waitFor(() => {
                expect(result1.current[0]).toBe("updated")
            })

            // Second hook should see the localStorage change on next render
            const {result: result3} = renderHook(() =>
                useLocalStorageState("shared-key", "initial")
            )

            expect(result3.current[0]).toBe("updated")
        })
    })
})
