import {beforeEach, describe, expect, it, vi} from "vitest"

import type {ActiveCluster} from "../../src/api/cluster/type"
import {InstanceTabType} from "../../src/api/instance/type"
import {QueryType} from "../../src/api/query/type"
import {useStore, useStoreAction} from "../../src/provider/StoreProvider"
import {createMockCluster, createMockInstance} from "../test-helpers"

// Mock MainQueryClient from AppProvider
vi.mock("../../src/provider/AppProvider", async (importOriginal) => {
    const actual = await importOriginal<typeof import("../../src/provider/AppProvider")>()
    return {
        ...actual,
        MainQueryClient: {
            clear: vi.fn(),
            setDefaultOptions: vi.fn(),
        },
    }
})

describe("StoreProvider", () => {
    beforeEach(() => {
        // Reset store to initial state before each test
        useStore.setState(useStore.getInitialState())
    })

    describe("Initial state", () => {
        it("should have correct initial values", () => {
            const state = useStore.getState()

            expect(state.activeClusterTab).toBe(0)
            expect(state.activeCluster).toBeUndefined()
            expect(state.activeInstance).toEqual({})
            expect(state.activeTags).toEqual(["ALL"])
            expect(state.warnings).toEqual({})
            expect(state.settings).toBe(false)
            expect(state.instance.body).toBe(InstanceTabType.CHART)
            expect(state.instance.queryTab).toBe(QueryType.CONSOLE)
            expect(state.instance.queryConsole).toBe("")
            expect(state.instance.dbName).toBeUndefined()
        })
    })

    describe("setCluster", () => {
        it("should set active cluster", () => {
            const cluster: ActiveCluster = {
                cluster: createMockCluster({name: "test-cluster"}),
                warning: false,
            }

            useStoreAction.setCluster(cluster)

            const state = useStore.getState()
            expect(state.activeCluster).toEqual(cluster)
        })

        it("should clear active cluster when undefined", () => {
            const cluster: ActiveCluster = {
                cluster: createMockCluster({name: "test-cluster"}),
                warning: false,
            }

            useStoreAction.setCluster(cluster)
            useStoreAction.setCluster(undefined)

            const state = useStore.getState()
            expect(state.activeCluster).toBeUndefined()
        })
    })


    describe("setClusterDetection", () => {
        it("should update detection instance", () => {
            const cluster: ActiveCluster = {
                cluster: createMockCluster({name: "test-cluster"}),
                warning: false,
            }

            useStoreAction.setCluster(cluster)

            const detectInstance = createMockInstance({
                sidecar: {host: "localhost", port: 8009},
            })
            useStoreAction.setClusterDetection(detectInstance)

            const state = useStore.getState()
            expect(state.activeCluster?.detectBy).toEqual(detectInstance)
        })

        it("should not update if no active cluster", () => {
            const instance = createMockInstance()

            const stateBefore = useStore.getState()
            useStoreAction.setClusterDetection(instance)
            const stateAfter = useStore.getState()

            expect(stateAfter).toEqual(stateBefore)
        })
    })

    describe("setClusterTab", () => {
        it("should set active cluster tab", () => {
            useStoreAction.setClusterTab(2)

            const state = useStore.getState()
            expect(state.activeClusterTab).toBe(2)
        })
    })

    describe("setWarnings", () => {
        it("should add a warning", () => {
            useStoreAction.setWarnings("test-warning", true)

            const state = useStore.getState()
            expect(state.warnings["test-warning"]).toBe(true)
        })

        it("should update existing warning", () => {
            useStoreAction.setWarnings("test-warning", true)
            useStoreAction.setWarnings("test-warning", false)

            const state = useStore.getState()
            expect(state.warnings["test-warning"]).toBe(false)
        })

        it("should handle multiple warnings", () => {
            useStoreAction.setWarnings("warning1", true)
            useStoreAction.setWarnings("warning2", false)

            const state = useStore.getState()
            expect(state.warnings["warning1"]).toBe(true)
            expect(state.warnings["warning2"]).toBe(false)
        })
    })

    describe("setInstance", () => {
        it("should set active instance for cluster", () => {
            const cluster: ActiveCluster = {
                cluster: createMockCluster({name: "test-cluster"}),
                warning: false,
            }

            useStoreAction.setCluster(cluster)

            const instance = createMockInstance({
                sidecar: {host: "localhost", port: 8009},
            })

            useStoreAction.setInstance(instance)

            const state = useStore.getState()
            expect(state.activeInstance["test-cluster"]).toEqual(instance)
        })

        it("should remove active instance when undefined", () => {
            const cluster: ActiveCluster = {
                cluster: createMockCluster({name: "test-cluster"}),
                warning: false,
            }

            useStoreAction.setCluster(cluster)

            const instance = createMockInstance({
                sidecar: {host: "localhost", port: 8009},
            })

            useStoreAction.setInstance(instance)
            useStoreAction.setInstance(undefined)

            const state = useStore.getState()
            expect(state.activeInstance["test-cluster"]).toBeUndefined()
        })

        it("should not update if no active cluster", () => {
            const instance = createMockInstance()

            const stateBefore = useStore.getState()
            useStoreAction.setInstance(instance)
            const stateAfter = useStore.getState()

            expect(stateAfter).toEqual(stateBefore)
        })
    })

    describe("setTags", () => {
        it("should set active tags", () => {
            useStoreAction.setTags(["tag1", "tag2"])

            const state = useStore.getState()
            expect(state.activeTags).toEqual(["tag1", "tag2"])
        })

        it("should replace existing tags", () => {
            useStoreAction.setTags(["tag1"])
            useStoreAction.setTags(["tag2", "tag3"])

            const state = useStore.getState()
            expect(state.activeTags).toEqual(["tag2", "tag3"])
        })
    })

    describe("toggleSettingsDialog", () => {
        it("should toggle settings from false to true", () => {
            useStoreAction.toggleSettingsDialog()

            const state = useStore.getState()
            expect(state.settings).toBe(true)
        })

        it("should toggle settings from true to false", () => {
            useStoreAction.toggleSettingsDialog()
            useStoreAction.toggleSettingsDialog()

            const state = useStore.getState()
            expect(state.settings).toBe(false)
        })
    })

    describe("setConsoleQuery", () => {
        it("should set console query", () => {
            useStoreAction.setConsoleQuery("SELECT * FROM users")

            const state = useStore.getState()
            expect(state.instance.queryConsole).toBe("SELECT * FROM users")
        })
    })

    describe("setInstanceBody", () => {
        it("should set instance body tab", () => {
            useStoreAction.setInstanceBody(InstanceTabType.QUERY)

            const state = useStore.getState()
            expect(state.instance.body).toBe(InstanceTabType.QUERY)
        })
    })

    describe("setQueryTab", () => {
        it("should set query tab", () => {
            useStoreAction.setQueryTab(QueryType.ACTIVITY)

            const state = useStore.getState()
            expect(state.instance.queryTab).toBe(QueryType.ACTIVITY)
        })
    })

    describe("setDbName", () => {
        it("should set database name", () => {
            useStoreAction.setDbName("testdb")

            const state = useStore.getState()
            expect(state.instance.dbName).toBe("testdb")
        })

        it("should clear database name when undefined", () => {
            useStoreAction.setDbName("testdb")
            useStoreAction.setDbName(undefined)

            const state = useStore.getState()
            expect(state.instance.dbName).toBeUndefined()
        })
    })

})
