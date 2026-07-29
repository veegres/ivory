import {afterEach, beforeEach, describe, expect, it, rs} from "@rstest/core"
import {act, renderHook, screen} from "@testing-library/react"

import {SnackbarProvide} from "../provider/SnackbarProvider"
import {useQueryParamErrorHandler} from "./QueryParamErrorHandler"

describe("useQueryParamErrorHandler", () => {
    beforeEach(() => {
        rs.useFakeTimers()
    })

    afterEach(() => {
        rs.useRealTimers()
        window.history.pushState({}, "", "/")
    })

    it("does nothing when there is no error query param", async () => {
        window.history.pushState({}, "", "/some/path")
        renderHook(() => useQueryParamErrorHandler(), {wrapper: SnackbarProvide})

        await act(async () => {
            await rs.advanceTimersByTimeAsync(200)
        })

        expect(screen.queryByRole("alert")).not.toBeInTheDocument()
    })

    it("strips the error param from the url", () => {
        window.history.pushState({}, "", "/some/path?error=oops")
        renderHook(() => useQueryParamErrorHandler(), {wrapper: SnackbarProvide})

        expect(window.location.search).toBe("")
        expect(window.location.pathname).toBe("/some/path")
    })

    it("shows a snackbar with the error message after a short delay", async () => {
        window.history.pushState({}, "", "/some/path?error=oops")
        renderHook(() => useQueryParamErrorHandler(), {wrapper: SnackbarProvide})

        await act(async () => {
            await rs.advanceTimersByTimeAsync(100)
        })

        expect(screen.getByText("EXTERNAL - ERROR, oops")).toBeInTheDocument()
    })
})
