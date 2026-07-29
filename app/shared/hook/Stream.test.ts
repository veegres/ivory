import {afterEach, beforeEach, describe, expect, it, rs} from "@rstest/core"
import {act, renderHook, waitFor} from "@testing-library/react"

import {EventStreamType, EventType} from "../../tools/pg_compacttable/api/job/PgCompactTableJobType"
import {useStream} from "./Stream"

class MockEventSource {
    static instances: MockEventSource[] = []

    url: string
    onopen?: () => void
    onerror?: () => void
    closed = false
    private listeners: Record<string, ((e: {data: string}) => void)[]> = {}

    constructor(url: string) {
        this.url = url
        MockEventSource.instances.push(this)
    }

    addEventListener(type: string, callback: (e: {data: string}) => void) {
        (this.listeners[type] ??= []).push(callback)
    }

    dispatch(type: string, data: string) {
        this.listeners[type]?.forEach(callback => callback({data}))
    }

    close() {
        this.closed = true
    }
}

describe("useStream", () => {
    beforeEach(() => {
        MockEventSource.instances = []
        rs.stubGlobal("EventSource", MockEventSource)
    })

    afterEach(() => {
        rs.unstubAllGlobals()
    })

    it("should not connect when disabled", () => {
        renderHook(() => useStream("/stream", {enabled: false}))

        expect(MockEventSource.instances).toHaveLength(0)
    })

    it("should not connect without a url", () => {
        renderHook(() => useStream(""))

        expect(MockEventSource.instances).toHaveLength(0)
    })

    it("should push a browser message and start loading on open", () => {
        const {result} = renderHook(() => useStream("/stream"))
        const source = MockEventSource.instances[0]

        act(() => source.onopen?.())

        expect(result.current.loading).toBe(true)
        expect(result.current.response).toEqual(["[browser] new connection was established"])
    })

    it("should push events dispatched by the server", () => {
        const {result} = renderHook(() => useStream("/stream"))
        const source = MockEventSource.instances[0]

        act(() => source.dispatch(EventType.LOG, "doing work"))

        expect(result.current.response).toContain("[log] doing work")
    })

    it("should stop loading and close the connection when the stream ends", () => {
        const {result} = renderHook(() => useStream("/stream"))
        const source = MockEventSource.instances[0]

        act(() => source.onopen?.())
        act(() => source.dispatch(EventType.STREAM, EventStreamType.END))

        expect(result.current.loading).toBe(false)
        expect(source.closed).toBe(true)
    })

    it("should push a reconnect message and stop loading on error", () => {
        const {result} = renderHook(() => useStream("/stream"))
        const source = MockEventSource.instances[0]

        act(() => source.onopen?.())
        act(() => source.onerror?.())

        expect(result.current.loading).toBe(false)
        expect(result.current.response).toContain("[browser] trying to reestablish connection")
    })

    it("should give up and close the connection after repeated consecutive errors", () => {
        const {result} = renderHook(() => useStream("/stream"))
        const source = MockEventSource.instances[0]

        act(() => source.onopen?.())
        act(() => source.onerror?.())
        act(() => source.onerror?.())
        act(() => source.onerror?.())

        expect(source.closed).toBe(true)
        expect(result.current.response).toContain("[browser] connection failed repeatedly, giving up")
    })

    it("should reset the error count after a successful reconnect", () => {
        const {result} = renderHook(() => useStream("/stream"))
        const source = MockEventSource.instances[0]

        act(() => source.onopen?.())
        act(() => source.onerror?.())
        act(() => source.onerror?.())
        act(() => source.onopen?.())
        act(() => source.onerror?.())

        expect(source.closed).toBe(false)
        expect(result.current.response).not.toContain("[browser] connection failed repeatedly, giving up")
    })

    it("should respect a custom maxConsecutiveErrors option", () => {
        const {result} = renderHook(() => useStream("/stream", {maxConsecutiveErrors: 2}))
        const source = MockEventSource.instances[0]

        act(() => source.onopen?.())
        act(() => source.onerror?.())
        act(() => source.onerror?.())

        expect(source.closed).toBe(true)
        expect(result.current.response).toContain("[browser] connection failed repeatedly, giving up")
    })

    it("should open a fresh connection when reconnect is called", async () => {
        const {result} = renderHook(() => useStream("/stream"))

        result.current.reconnect()

        await waitFor(() => expect(MockEventSource.instances).toHaveLength(2))
        expect(MockEventSource.instances[0].closed).toBe(true)
    })

    it("should close the connection and clear state on unmount", () => {
        const {unmount} = renderHook(() => useStream("/stream"))
        const source = MockEventSource.instances[0]

        unmount()

        expect(source.closed).toBe(true)
    })
})
