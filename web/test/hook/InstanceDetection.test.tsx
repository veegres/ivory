import {QueryClient, QueryClientProvider} from "@tanstack/react-query"
import {renderHook, waitFor} from "@testing-library/react"
import {beforeEach, describe, expect, it, vi} from "vitest"

import type {ActiveCluster} from "../../src/api/cluster/type"
import {useInstanceDetection} from "../../src/hook/InstanceDetection"
import * as StoreProvider from "../../src/provider/StoreProvider"
import {createMockCluster, createMockInstanceWeb} from "../test-helpers"

// Mock the main module to avoid DOM dependencies
vi.mock("../../src/app/main", () => ({
    MainQueryClient: {
        clear: vi.fn(),
    },
}))

// Mock the router hook
vi.mock("../../src/api/instance/hook", () => ({
    useRouterInstanceOverview: vi.fn(() => ({
        data: {},
        error: null,
        isFetching: false,
        errorUpdateCount: 0,
        refetch: vi.fn(),
    })),
}))

// Import after mocking
import {useRouterInstanceOverview} from "../../src/api/instance/hook"

describe("useInstanceDetection", () => {
    let queryClient: QueryClient
    let mockSetCluster: ReturnType<typeof vi.fn>
    let mockSetWarnings: ReturnType<typeof vi.fn>

    beforeEach(() => {
        queryClient = new QueryClient({
            defaultOptions: {
                queries: {retry: false},
            },
        })

        mockSetCluster = vi.fn()
        mockSetWarnings = vi.fn()

        // Mock useStore and useStoreAction
        vi.spyOn(StoreProvider, "useStore").mockReturnValue(undefined)
        vi.spyOn(StoreProvider, "useStoreAction", "get").mockReturnValue({
            setCluster: mockSetCluster,
            setWarnings: mockSetWarnings,
        } as any)

        // Reset the router hook mock
        vi.mocked(useRouterInstanceOverview).mockReturnValue({
            data: {},
            error: null,
            isFetching: false,
            errorUpdateCount: 0,
            refetch: vi.fn().mockResolvedValue({}),
        } as any)
    })

    function wrapper({children}: {children: React.ReactNode}) {
        return <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    }

    describe("Initial state", () => {
        it("should return detection object with auto detection when no active cluster", () => {
            const cluster = createMockCluster({name: "test-cluster"})
            const instances = [{host: "localhost", port: 8008}]

            const {result} = renderHook(() => useInstanceDetection(cluster, instances), {wrapper})

            expect(result.current.detection).toBe("auto")
            expect(result.current.active).toBe(false)
            expect(result.current.fetching).toBe(false)
        })

        it("should return default instance from first sidecar when no data", () => {
            const cluster = createMockCluster({name: "test-cluster"})
            const instances = [{host: "localhost", port: 8008}]

            const {result} = renderHook(() => useInstanceDetection(cluster, instances), {wrapper})

            expect(result.current.defaultInstance.sidecar).toEqual({host: "localhost", port: 8008})
        })

        it("should create combined instance map from config instances when no data", () => {
            const cluster = createMockCluster({name: "test-cluster"})
            const instances = [{host: "localhost", port: 8008}]

            const {result} = renderHook(() => useInstanceDetection(cluster, instances), {wrapper})

            // The hook combines config instances with query data, creating initial instances
            expect(result.current.combinedInstanceMap["localhost:8008"]).toBeDefined()
            expect(result.current.combinedInstanceMap["localhost:8008"].inInstances).toBe(true)
        })

        it("should generate warning colors for instances with no data", () => {
            const cluster = createMockCluster({name: "test-cluster"})
            const instances = [{host: "localhost", port: 8008}]

            const {result} = renderHook(() => useInstanceDetection(cluster, instances), {wrapper})

            // Instances without data get warning color
            expect(result.current.colors["localhost:8008"]).toBe("warning")
        })

        it("should set warning when instances are not in cluster", () => {
            const cluster = createMockCluster({name: "test-cluster"})
            const instances = [{host: "localhost", port: 8008}]

            const {result} = renderHook(() => useInstanceDetection(cluster, instances), {wrapper})

            // Warning is true when instances exist in config but not in cluster data
            expect(result.current.warning).toBe(true)
        })
    })

    describe("Active cluster detection", () => {
        it("should return active as true when cluster is active", () => {
            const cluster = createMockCluster({name: "test-cluster"})
            const instances = [{host: "localhost", port: 8008}]

            const activeCluster: ActiveCluster = {
                cluster,
                defaultInstance: createMockInstanceWeb(),
                combinedInstanceMap: {},
                warning: false,
                detection: "auto",
            }

            vi.spyOn(StoreProvider, "useStore").mockReturnValue(activeCluster)

            const {result} = renderHook(() => useInstanceDetection(cluster, instances), {wrapper})

            expect(result.current.active).toBe(true)
        })

        it("should use activeCluster detection when cluster is active", () => {
            const cluster = createMockCluster({name: "test-cluster"})
            const instances = [{host: "localhost", port: 8008}]

            const mockInstance = createMockInstanceWeb({
                sidecar: {host: "localhost", port: 8008},
            })

            const activeCluster: ActiveCluster = {
                cluster,
                defaultInstance: mockInstance,
                combinedInstanceMap: {"localhost:8008": mockInstance},
                warning: false,
                detection: "manual",
            }

            vi.spyOn(StoreProvider, "useStore").mockReturnValue(activeCluster)

            vi.mocked(useRouterInstanceOverview).mockReturnValue({
                data: {"localhost:8008": mockInstance},
                error: null,
                isFetching: false,
                errorUpdateCount: 0,
                refetch: vi.fn().mockResolvedValue({}),
            } as any)

            const {result} = renderHook(() => useInstanceDetection(cluster, instances), {wrapper})

            expect(result.current.detection).toBe("manual")
        })
    })

    describe("Data fetching", () => {
        it("should show fetching state when query is loading", () => {
            const cluster = createMockCluster({name: "test-cluster"})
            const instances = [{host: "localhost", port: 8008}]

            vi.mocked(useRouterInstanceOverview).mockReturnValue({
                data: undefined,
                error: null,
                isFetching: true,
                errorUpdateCount: 0,
                refetch: vi.fn().mockResolvedValue({}),
            } as any)

            const {result} = renderHook(() => useInstanceDetection(cluster, instances), {wrapper})

            expect(result.current.fetching).toBe(true)
        })

        it("should use data from query when available", () => {
            const cluster = createMockCluster({name: "test-cluster"})
            const instances = [{host: "localhost", port: 8008}]

            const mockInstance = createMockInstanceWeb({
                sidecar: {host: "localhost", port: 8008},
                role: "leader",
            })

            vi.mocked(useRouterInstanceOverview).mockReturnValue({
                data: {"localhost:8008": mockInstance},
                error: null,
                isFetching: false,
                errorUpdateCount: 0,
                refetch: vi.fn().mockResolvedValue({}),
            } as any)

            const {result} = renderHook(() => useInstanceDetection(cluster, instances), {wrapper})

            expect(result.current.combinedInstanceMap["localhost:8008"]).toBeDefined()
        })
    })

    describe("Leader selection in auto mode", () => {
        it("should select leader instance as default when available in auto mode", () => {
            const cluster = createMockCluster({name: "test-cluster"})
            const instances = [
                {host: "localhost", port: 8008},
                {host: "localhost", port: 8009},
            ]

            const mockLeader = createMockInstanceWeb({
                sidecar: {host: "localhost", port: 8009},
                role: "leader",
                leader: true,
            })

            const mockReplica = createMockInstanceWeb({
                sidecar: {host: "localhost", port: 8008},
                role: "replica",
                leader: false,
            })

            vi.mocked(useRouterInstanceOverview).mockReturnValue({
                data: {
                    "localhost:8008": mockReplica,
                    "localhost:8009": mockLeader,
                },
                error: null,
                isFetching: false,
                errorUpdateCount: 0,
                refetch: vi.fn().mockResolvedValue({}),
            } as any)

            const {result} = renderHook(() => useInstanceDetection(cluster, instances), {wrapper})

            expect(result.current.defaultInstance.sidecar.port).toBe(8009)
            expect(result.current.defaultInstance.leader).toBe(true)
        })

        it("should fallback to first instance when no leader in auto mode", () => {
            const cluster = createMockCluster({name: "test-cluster"})
            const instances = [
                {host: "localhost", port: 8008},
                {host: "localhost", port: 8009},
            ]

            const mockReplica1 = createMockInstanceWeb({
                sidecar: {host: "localhost", port: 8008},
                role: "replica",
                leader: false,
            })

            const mockReplica2 = createMockInstanceWeb({
                sidecar: {host: "localhost", port: 8009},
                role: "replica",
                leader: false,
            })

            vi.mocked(useRouterInstanceOverview).mockReturnValue({
                data: {
                    "localhost:8008": mockReplica1,
                    "localhost:8009": mockReplica2,
                },
                error: null,
                isFetching: false,
                errorUpdateCount: 0,
                refetch: vi.fn().mockResolvedValue({}),
            } as any)

            const {result} = renderHook(() => useInstanceDetection(cluster, instances), {wrapper})

            expect(result.current.defaultInstance.sidecar.port).toBe(8008)
        })
    })

    describe("Manual detection mode", () => {
        it("should use activeCluster default instance in manual mode when cluster is active", () => {
            const cluster = createMockCluster({name: "test-cluster"})
            const instances = [
                {host: "localhost", port: 8008},
                {host: "localhost", port: 8009},
            ]

            const mockInstance = createMockInstanceWeb({
                sidecar: {host: "localhost", port: 8009},
            })

            const activeCluster: ActiveCluster = {
                cluster,
                defaultInstance: mockInstance,
                combinedInstanceMap: {"localhost:8009": mockInstance},
                warning: false,
                detection: "manual",
            }

            vi.spyOn(StoreProvider, "useStore").mockReturnValue(activeCluster)

            vi.mocked(useRouterInstanceOverview).mockReturnValue({
                data: {"localhost:8009": mockInstance},
                error: null,
                isFetching: false,
                errorUpdateCount: 0,
                refetch: vi.fn().mockResolvedValue({}),
            } as any)

            const {result} = renderHook(() => useInstanceDetection(cluster, instances), {wrapper})

            expect(result.current.detection).toBe("manual")
            expect(result.current.defaultInstance.sidecar.port).toBe(8009)
        })
    })

    describe("Warning state", () => {
        it("should set warning in store when warning is true", async () => {
            const cluster = createMockCluster({name: "test-cluster"})
            const instances = [{host: "localhost", port: 8008}]

            const mockInstance = createMockInstanceWeb({
                sidecar: {host: "localhost", port: 8008},
                inCluster: false,
            })

            vi.mocked(useRouterInstanceOverview).mockReturnValue({
                data: {"localhost:8008": mockInstance},
                error: null,
                isFetching: false,
                errorUpdateCount: 0,
                refetch: vi.fn().mockResolvedValue({}),
            } as any)

            renderHook(() => useInstanceDetection(cluster, instances), {wrapper})

            await waitFor(() => {
                expect(mockSetWarnings).toHaveBeenCalledWith("test-cluster", true)
            })
        })
    })

    describe("Refetch functionality", () => {
        it("should provide refetch function", () => {
            const cluster = createMockCluster({name: "test-cluster"})
            const instances = [{host: "localhost", port: 8008}]

            const {result} = renderHook(() => useInstanceDetection(cluster, instances), {wrapper})

            expect(typeof result.current.refetch).toBe("function")
        })

        it("should call refetch when refetch function is invoked", async () => {
            const cluster = createMockCluster({name: "test-cluster"})
            const instances = [{host: "localhost", port: 8008}]

            const mockRefetch = vi.fn().mockResolvedValue({})

            vi.mocked(useRouterInstanceOverview).mockReturnValue({
                data: {},
                error: null,
                isFetching: false,
                errorUpdateCount: 0,
                refetch: mockRefetch,
            } as any)

            const {result} = renderHook(() => useInstanceDetection(cluster, instances), {wrapper})

            await result.current.refetch()

            expect(mockRefetch).toHaveBeenCalled()
        })
    })

    describe("Error handling", () => {
        it("should still create instances from config on error", () => {
            const cluster = createMockCluster({name: "test-cluster"})
            const instances = [{host: "localhost", port: 8008}]

            vi.mocked(useRouterInstanceOverview).mockReturnValue({
                data: undefined,
                error: new Error("Network error"),
                isFetching: false,
                errorUpdateCount: 1,
                refetch: vi.fn().mockResolvedValue({}),
            } as any)

            const {result} = renderHook(() => useInstanceDetection(cluster, instances), {wrapper})

            // Even on error, the hook creates instances from config
            expect(result.current.combinedInstanceMap["localhost:8008"]).toBeDefined()
            expect(result.current.combinedInstanceMap["localhost:8008"].inInstances).toBe(true)
            expect(result.current.combinedInstanceMap["localhost:8008"].inCluster).toBe(false)
        })
    })

    describe("Color generation", () => {
        it("should generate colors for instances", () => {
            const cluster = createMockCluster({name: "test-cluster"})
            const instances = [
                {host: "localhost", port: 8008},
                {host: "localhost", port: 8009},
            ]

            const mockInstance1 = createMockInstanceWeb({
                sidecar: {host: "localhost", port: 8008},
            })

            const mockInstance2 = createMockInstanceWeb({
                sidecar: {host: "localhost", port: 8009},
            })

            vi.mocked(useRouterInstanceOverview).mockReturnValue({
                data: {
                    "localhost:8008": mockInstance1,
                    "localhost:8009": mockInstance2,
                },
                error: null,
                isFetching: false,
                errorUpdateCount: 0,
                refetch: vi.fn().mockResolvedValue({}),
            } as any)

            const {result} = renderHook(() => useInstanceDetection(cluster, instances), {wrapper})

            expect(Object.keys(result.current.colors).length).toBeGreaterThan(0)
        })
    })

    describe("Store updates when cluster is active", () => {
        it("should update store when cluster is active", async () => {
            const cluster = createMockCluster({name: "test-cluster"})
            const instances = [{host: "localhost", port: 8008}]

            const mockInstance = createMockInstanceWeb({
                sidecar: {host: "localhost", port: 8008},
            })

            const activeCluster: ActiveCluster = {
                cluster,
                defaultInstance: mockInstance,
                combinedInstanceMap: {},
                warning: false,
                detection: "auto",
            }

            vi.spyOn(StoreProvider, "useStore").mockReturnValue(activeCluster)

            vi.mocked(useRouterInstanceOverview).mockReturnValue({
                data: {"localhost:8008": mockInstance},
                error: null,
                isFetching: false,
                errorUpdateCount: 0,
                refetch: vi.fn().mockResolvedValue({}),
            } as any)

            renderHook(() => useInstanceDetection(cluster, instances), {wrapper})

            await waitFor(() => {
                expect(mockSetCluster).toHaveBeenCalled()
            })
        })
    })
})
