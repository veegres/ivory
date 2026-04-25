import {beforeEach, describe, expect, it, vi} from "vitest"

import type {ActiveCluster} from "../../src/api/cluster/type"
import {QueryType} from "../../src/api/database/type"
import {NodeTabType} from "../../src/api/node/type"
import {getDomain} from "../../src/app/utils"
import {MainQueryClient} from "../../src/provider/AppProvider"
import {useStore, useStoreAction} from "../../src/provider/StoreProvider"
import {createMockCluster, createMockNode} from "../test-helpers"

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
        // Clear all mocks
        vi.clearAllMocks()
    })

    describe("Initial state", () => {
        it("should have correct initial values", () => {
            const state = useStore.getState()

            expect(state.searchCluster).toBe("")
            expect(state.activeClusterTab).toBe(0)
            expect(state.activeCluster).toBeUndefined()
            expect(state.activeNode).toEqual({})
            expect(state.activeTags).toEqual(["ALL"])
            expect(state.warnings).toEqual({})
            expect(state.settings).toBe(false)
            expect(state.node.body).toBe(NodeTabType.CHART)
            expect(state.node.queryTab).toBe(QueryType.CONSOLE)
            expect(state.node.queryConsole).toBe("")
            expect(state.node.dbName).toBeUndefined()
            expect(state.node.dbSchema).toBeUndefined()
        })
    })

    describe("setSearchCluster", () => {
        it("should set search cluster text", () => {
            useStoreAction.setSearchCluster("test-search")

            const state = useStore.getState()
            expect(state.searchCluster).toBe("test-search")
        })

        it("should update search cluster text", () => {
            useStoreAction.setSearchCluster("first")
            useStoreAction.setSearchCluster("second")

            const state = useStore.getState()
            expect(state.searchCluster).toBe("second")
        })

        it("should clear search cluster when empty string", () => {
            useStoreAction.setSearchCluster("test")
            useStoreAction.setSearchCluster("")

            const state = useStore.getState()
            expect(state.searchCluster).toBe("")
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
        it("should update detection node", () => {
            const cluster: ActiveCluster = {
                cluster: createMockCluster({name: "test-cluster"}),
                warning: false,
            }

            useStoreAction.setCluster(cluster)

            const detectNode = createMockNode({
                connection: {sshKeyId: "00000000-0000-0000-0000-000000000000", host: "localhost", keeperPort: 8009},
            })
            useStoreAction.setClusterDetection(detectNode)

            const state = useStore.getState()
            expect(state.activeCluster?.detectBy).toEqual(detectNode)
        })

        it("should not update if no active cluster", () => {
            const node = createMockNode()

            const stateBefore = useStore.getState()
            useStoreAction.setClusterDetection(node)
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

    describe("setNode", () => {
        it("should set active node for cluster", () => {
            const cluster: ActiveCluster = {
                cluster: createMockCluster({name: "test-cluster"}),
                warning: false,
            }

            useStoreAction.setCluster(cluster)

            const node = createMockNode({
                connection: {sshKeyId: "00000000-0000-0000-0000-000000000000", host: "localhost", keeperPort: 8009},
            })

            useStoreAction.setNode(getDomain(node.connection))

            const state = useStore.getState()
            expect(state.activeNode["test-cluster"]).toEqual(getDomain(node.connection))
        })

        it("should remove active node when undefined", () => {
            const cluster: ActiveCluster = {
                cluster: createMockCluster({name: "test-cluster"}),
                warning: false,
            }

            useStoreAction.setCluster(cluster)

            const node = createMockNode({
                connection: {sshKeyId: "00000000-0000-0000-0000-000000000000", host: "localhost", keeperPort: 8009},
            })

            useStoreAction.setNode(getDomain(node.connection))
            useStoreAction.setNode(undefined)

            const state = useStore.getState()
            expect(state.activeNode["test-cluster"]).toBeUndefined()
        })

        it("should not update if no active cluster", () => {
            const node = createMockNode()

            const stateBefore = useStore.getState()
            useStoreAction.setNode(getDomain(node.connection))
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
            expect(state.node.queryConsole).toBe("SELECT * FROM users")
        })
    })

    describe("setNodeBody", () => {
        it("should set node body tab", () => {
            useStoreAction.setNodeBody(NodeTabType.QUERY)

            const state = useStore.getState()
            expect(state.node.body).toBe(NodeTabType.QUERY)
        })
    })

    describe("setQueryTab", () => {
        it("should set query tab", () => {
            useStoreAction.setQueryTab(QueryType.ACTIVITY)

            const state = useStore.getState()
            expect(state.node.queryTab).toBe(QueryType.ACTIVITY)
        })
    })

    describe("setDbName", () => {
        it("should set database name", () => {
            useStoreAction.setDbName("testdb")

            const state = useStore.getState()
            expect(state.node.dbName).toBe("testdb")
        })

        it("should clear database name when undefined", () => {
            useStoreAction.setDbName("testdb")
            useStoreAction.setDbName(undefined)

            const state = useStore.getState()
            expect(state.node.dbName).toBeUndefined()
        })
    })

    describe("setDbSchema", () => {
        it("should set database schema", () => {
            useStoreAction.setDbSchema("public")

            const state = useStore.getState()
            expect(state.node.dbSchema).toBe("public")
        })

        it("should update database schema", () => {
            useStoreAction.setDbSchema("public")
            useStoreAction.setDbSchema("private")

            const state = useStore.getState()
            expect(state.node.dbSchema).toBe("private")
        })

        it("should clear database schema when undefined", () => {
            useStoreAction.setDbSchema("public")
            useStoreAction.setDbSchema(undefined)

            const state = useStore.getState()
            expect(state.node.dbSchema).toBeUndefined()
        })
    })

    describe("clear", () => {
        it("should reset store to initial state", () => {
            // Set various values
            useStoreAction.setSearchCluster("test-search")
            useStoreAction.setClusterTab(2)
            useStoreAction.setTags(["tag1", "tag2"])
            useStoreAction.setWarnings("test-warning", true)
            useStoreAction.toggleSettingsDialog()
            useStoreAction.setConsoleQuery("SELECT * FROM users")
            useStoreAction.setDbName("testdb")
            useStoreAction.setDbSchema("public")

            // Clear store
            useStoreAction.clear()

            const state = useStore.getState()
            const initialState = useStore.getInitialState()

            expect(state).toEqual(initialState)
        })

        it("should call MainQueryClient.clear", () => {
            useStoreAction.clear()

            expect(MainQueryClient.clear).toHaveBeenCalled()
        })

        it("should reset activeCluster", () => {
            const cluster: ActiveCluster = {
                cluster: createMockCluster({name: "test-cluster"}),
                warning: false,
            }

            useStoreAction.setCluster(cluster)
            useStoreAction.clear()

            const state = useStore.getState()
            expect(state.activeCluster).toBeUndefined()
        })

        it("should reset activeNode", () => {
            const cluster: ActiveCluster = {
                cluster: createMockCluster({name: "test-cluster"}),
                warning: false,
            }

            useStoreAction.setCluster(cluster)

            const node = createMockNode({
                connection: {sshKeyId: "00000000-0000-0000-0000-000000000000", host: "localhost", keeperPort: 8009},
            })

            useStoreAction.setNode(getDomain(node.connection))
            useStoreAction.clear()

            const state = useStore.getState()
            expect(state.activeNode).toEqual({})
        })
    })

})
