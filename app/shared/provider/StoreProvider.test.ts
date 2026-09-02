import {beforeEach, describe, expect, it, rs} from "@rstest/core"

import {Cluster} from "../../features/cluster/api/ClusterType"
import {KeeperPlugin, NodeTabType} from "../../features/node/api/NodeType"
import {Type as QueryType} from "../../features/query/api/QueryType"
import {getDomain} from "../helper/HelperUtils"
import {createMockCluster, createMockNode} from "../test/TestMocks"
import * as actualAppProvider from "./AppProvider" with {rstest: "importActual"}
import {MainQueryClient} from "./AppProvider"
import {useStore, useStoreAction} from "./StoreProvider"

// Mock MainQueryClient from AppProvider
rs.mock("./AppProvider", () => ({
    ...actualAppProvider,
    MainQueryClient: {
        clear: rs.fn(),
        setDefaultOptions: rs.fn(),
    },
}))

describe("StoreProvider", () => {
    beforeEach(() => {
        // Reset store to initial state before each test
        useStore.setState(useStore.getInitialState())
        // Clear all mocks
        rs.clearAllMocks()
    })

    describe("Initial state", () => {
        it("should have correct initial values", () => {
            const state = useStore.getState()

            expect(state.searchCluster).toBe("")
            expect(state.activeCluster).toBeUndefined()
            expect(state.manualKeeper).toBeUndefined()
            expect(state.activeNode).toEqual({})
            expect(state.activeTags).toEqual(["ALL"])
            expect(state.warnings).toEqual({})
            expect(state.nodeState.nodeTab).toBe(NodeTabType.CONTAINER)
            expect(state.nodeState.queryTab).toBe(QueryType.CONSOLE)
            expect(state.nodeState.queryConsole).toBe("")
            expect(state.nodeState.dbName).toBeUndefined()
            expect(state.nodeState.dbSchema).toBeUndefined()
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
            const cluster: Cluster = createMockCluster({name: "test-cluster"})

            useStoreAction.setCluster(cluster)

            const state = useStore.getState()
            expect(state.activeCluster).toEqual(cluster)
        })

        it("should clear active cluster when undefined", () => {
            const cluster: Cluster = createMockCluster({name: "test-cluster"})

            useStoreAction.setCluster(cluster)
            useStoreAction.setCluster(undefined)

            const state = useStore.getState()
            expect(state.activeCluster).toBeUndefined()
        })
    })

    describe("setClusterKeeperPlugin", () => {
        it("should set active cluster keeper plugin", () => {
            useStoreAction.setClusterKeeperPlugin(KeeperPlugin.NATIVE_ETCD)

            const state = useStore.getState()
            expect(state.activeClusterKeeperPlugin).toBe(KeeperPlugin.NATIVE_ETCD)
        })

        it("should clear active cluster and manual keeper when plugin changes", () => {
            const cluster: Cluster = createMockCluster({name: "test-cluster"})
            useStoreAction.setCluster(cluster)
            useStoreAction.setClusterKeeper("test-keeper")

            useStoreAction.setClusterKeeperPlugin(KeeperPlugin.NATIVE_ETCD)

            const state = useStore.getState()
            expect(state.activeCluster).toBeUndefined()
            expect(state.manualKeeper).toBeUndefined()
        })
    })

    describe("setClusterDetection", () => {
        it("should update detection node", () => {
            const cluster: Cluster = createMockCluster({name: "test-cluster"})
            useStoreAction.setCluster(cluster)

            const manualKeeper = "test-keeper"
            useStoreAction.setClusterKeeper(manualKeeper)

            const state = useStore.getState()
            expect(state.manualKeeper).toEqual(manualKeeper)
        })

        it("should not update if no active cluster", () => {
            const stateBefore = useStore.getState()
            useStoreAction.setClusterKeeper("test-keeper")
            const stateAfter = useStore.getState()
            expect(stateAfter).toEqual(stateBefore)
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
            const cluster: Cluster = createMockCluster({name: "test-cluster"})

            useStoreAction.setCluster(cluster)

            const node = createMockNode({
                config: {name: "node1", host: "localhost", keeperPort: 8009},
            })

            useStoreAction.setNode(getDomain(node.config))

            const state = useStore.getState()
            expect(state.activeNode["test-cluster"]).toEqual(getDomain(node.config))
        })

        it("should remove active node when undefined", () => {
            const cluster: Cluster = createMockCluster({name: "test-cluster"})

            useStoreAction.setCluster(cluster)

            const node = createMockNode({
                config: {name: "node1", host: "localhost", keeperPort: 8009},
            })

            useStoreAction.setNode(getDomain(node.config))
            useStoreAction.setNode(undefined)

            const state = useStore.getState()
            expect(state.activeNode["test-cluster"]).toBeUndefined()
        })

        it("should not update if no active cluster", () => {
            const node = createMockNode()

            const stateBefore = useStore.getState()
            useStoreAction.setNode(getDomain(node.config))
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

    describe("setConsoleQuery", () => {
        it("should set console query", () => {
            useStoreAction.setConsoleQuery("SELECT * FROM users")

            const state = useStore.getState()
            expect(state.nodeState.queryConsole).toBe("SELECT * FROM users")
        })
    })

    describe("setNodeBody", () => {
        it("should set node body tab", () => {
            useStoreAction.setNodeBody(NodeTabType.DATABASE)

            const state = useStore.getState()
            expect(state.nodeState.nodeTab).toBe(NodeTabType.DATABASE)
        })
    })

    describe("setQueryTab", () => {
        it("should set query tab", () => {
            useStoreAction.setQueryTab(QueryType.ACTIVITY)

            const state = useStore.getState()
            expect(state.nodeState.queryTab).toBe(QueryType.ACTIVITY)
        })
    })

    describe("setDbName", () => {
        it("should set database name", () => {
            useStoreAction.setDbName("testdb")

            const state = useStore.getState()
            expect(state.nodeState.dbName).toBe("testdb")
        })

        it("should clear database name when undefined", () => {
            useStoreAction.setDbName("testdb")
            useStoreAction.setDbName(undefined)

            const state = useStore.getState()
            expect(state.nodeState.dbName).toBeUndefined()
        })
    })

    describe("setDbSchema", () => {
        it("should set database schema", () => {
            useStoreAction.setDbSchema("public")

            const state = useStore.getState()
            expect(state.nodeState.dbSchema).toBe("public")
        })

        it("should update database schema", () => {
            useStoreAction.setDbSchema("public")
            useStoreAction.setDbSchema("private")

            const state = useStore.getState()
            expect(state.nodeState.dbSchema).toBe("private")
        })

        it("should clear database schema when undefined", () => {
            useStoreAction.setDbSchema("public")
            useStoreAction.setDbSchema(undefined)

            const state = useStore.getState()
            expect(state.nodeState.dbSchema).toBeUndefined()
        })
    })

    describe("clear", () => {
        it("should reset store to initial state", () => {
            // Set various values
            useStoreAction.setSearchCluster("test-search")
            useStoreAction.setTags(["tag1", "tag2"])
            useStoreAction.setWarnings("test-warning", true)
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
            const cluster: Cluster = createMockCluster({name: "test-cluster"})

            useStoreAction.setCluster(cluster)
            useStoreAction.clear()

            const state = useStore.getState()
            expect(state.activeCluster).toBeUndefined()
        })

        it("should reset activeNode", () => {
            const cluster: Cluster = createMockCluster({name: "test-cluster"})

            useStoreAction.setCluster(cluster)

            const node = createMockNode({
                config: {name: "node1", host: "localhost", keeperPort: 8009},
            })

            useStoreAction.setNode(getDomain(node.config))
            useStoreAction.clear()

            const state = useStore.getState()
            expect(state.activeNode).toEqual({})
        })
    })

    // NOTE: v1.4.2 persisted under this same key and this same name, so the
    // version is the only thing that can tell the two shapes apart - hydrating
    // the old blob left activeCluster with no nodes and threw on first render
    describe("Persisted version", () => {
        it("should discard a v1 blob instead of merging it over the defaults", async () => {
            localStorage.setItem("store", JSON.stringify({
                version: 1,
                state: {
                    searchCluster: "old",
                    activeClusterTab: 2,
                    activeCluster: {cluster: {name: "old-cluster"}, detectBy: "host", warning: false},
                    activeInstance: {"old-cluster": "localhost:8008"},
                    instance: {body: 0, queryTab: 0, queryConsole: "select 1"},
                },
            }))

            await useStore.persist.rehydrate()

            const state = useStore.getState()
            expect(state.activeCluster).toBeUndefined()
            expect(state.searchCluster).toBe("")
            expect(state.activeNode).toEqual({})
            expect(state.nodeState.nodeTab).toBe(NodeTabType.CONTAINER)
            expect(state.nodeState.queryConsole).toBe("")
        })

        it("should keep a blob written by the current version", async () => {
            localStorage.setItem("store", JSON.stringify({
                version: useStore.persist.getOptions().version,
                state: {searchCluster: "kept", nodeState: {queryConsole: "select 1"}},
            }))

            await useStore.persist.rehydrate()

            const state = useStore.getState()
            expect(state.searchCluster).toBe("kept")
            expect(state.nodeState.queryConsole).toBe("select 1")
            expect(state.nodeState.nodeTab).toBe(NodeTabType.CONTAINER)
        })
    })

})
