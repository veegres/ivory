import {QueryClient, QueryClientProvider} from "@tanstack/react-query"
import {render, screen} from "@testing-library/react"
import {AxiosError, AxiosHeaders, HttpStatusCode} from "axios"
import {describe, expect, it, vi} from "vitest"

import {api} from "../../features/Api"
import {ManagementApi} from "../../features/management/api/ManagementRouter"
import {AuthProvider} from "./AuthProvider"

type RejectedHandler = (error: AxiosError) => Promise<never>

function getRejectedHandler(): RejectedHandler {
    const handlers = (api.interceptors.response as unknown as {handlers: {rejected: RejectedHandler}[]}).handlers
    return handlers[handlers.length - 1].rejected
}

function makeError(status?: number): AxiosError {
    return new AxiosError(
        "boom",
        undefined,
        undefined,
        undefined,
        status ? {status, statusText: "", data: {}, headers: new AxiosHeaders(), config: {headers: new AxiosHeaders()}} : undefined,
    )
}

function renderWithClient(queryClient: QueryClient) {
    return render(
        <QueryClientProvider client={queryClient}>
            <AuthProvider><div>child</div></AuthProvider>
        </QueryClientProvider>
    )
}

describe("AuthProvider", () => {
    it("renders children", () => {
        renderWithClient(new QueryClient())

        expect(screen.getByText("child")).toBeInTheDocument()
    })

    it("refetches auth info on a 401 response", async () => {
        const queryClient = new QueryClient()
        const refetchSpy = vi.spyOn(queryClient, "refetchQueries").mockResolvedValue()
        renderWithClient(queryClient)

        await expect(getRejectedHandler()(makeError(HttpStatusCode.Unauthorized))).rejects.toBeTruthy()

        expect(refetchSpy).toHaveBeenCalledWith({queryKey: ManagementApi.info.key()})
    })

    it("refetches auth info on a 403 response", async () => {
        const queryClient = new QueryClient()
        const refetchSpy = vi.spyOn(queryClient, "refetchQueries").mockResolvedValue()
        renderWithClient(queryClient)

        await expect(getRejectedHandler()(makeError(HttpStatusCode.Forbidden))).rejects.toBeTruthy()

        expect(refetchSpy).toHaveBeenCalledWith({queryKey: ManagementApi.info.key()})
    })

    it("does not refetch on other errors", async () => {
        const queryClient = new QueryClient()
        const refetchSpy = vi.spyOn(queryClient, "refetchQueries").mockResolvedValue()
        renderWithClient(queryClient)

        await expect(getRejectedHandler()(makeError(HttpStatusCode.InternalServerError))).rejects.toBeTruthy()

        expect(refetchSpy).not.toHaveBeenCalled()
    })

    it("ejects the interceptor on unmount", () => {
        const {unmount} = renderWithClient(new QueryClient())
        const handlers = (api.interceptors.response as unknown as {handlers: unknown[]}).handlers
        const registeredCount = handlers.filter(Boolean).length

        unmount()

        const remainingCount = handlers.filter(Boolean).length
        expect(remainingCount).toBe(registeredCount - 1)
    })
})
