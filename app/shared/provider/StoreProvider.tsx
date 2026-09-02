import {persist} from "zustand/middleware"
import {create} from "zustand/react"

import {Cluster} from "../../features/cluster/api/ClusterType"
import {KeeperPlugin, NodeTabType} from "../../features/node/api/NodeType"
import {Type as QueryType} from "../../features/query/api/QueryType"
import {MainQueryClient} from "./AppProvider"

// STORE
interface Store {
    searchCluster: string,
    activeClusterKeeperPlugin: KeeperPlugin,
    activeCluster?: Cluster,
    manualKeeper?: string,
    activeNode: { [cluster: string]: string | undefined },
    activeTags: string[],
    warnings: { [key: string]: boolean },
    refresh: { [key: string]: [string, number] },
    nodeState: {
        nodeTab: NodeTabType,
        queryTab: QueryType,
        queryConsole: string,
        dbName?: string,
        dbSchema?: string,
        containerTab: number,
        systemTab: number,
        systemLogsPath: string,
        toolsTab?: number,
    },
}

const InitialStore: Store = {
    searchCluster: "",
    activeClusterKeeperPlugin: KeeperPlugin.PATRONI_POSTGRES,
    activeCluster: undefined,
    manualKeeper: undefined,
    activeNode: {},
    activeTags: ["ALL"],
    warnings: {},
    refresh: {},
    nodeState: {
        nodeTab: NodeTabType.CONTAINER,
        queryTab: QueryType.CONSOLE,
        queryConsole: "",
        dbName: undefined,
        dbSchema: undefined,
        containerTab: 0,
        systemTab: 0,
        systemLogsPath: "",
        toolsTab: undefined,
    },
}

export const useStore = create(persist<Store>(
    () => InitialStore,
    {
        name: "store",
        // NOTE: v1.4.2 persisted a completely different shape under this same
        //  key and this same version, so the old blob was merged over the new
        //  defaults instead of being discarded - activeCluster then had no
        //  nodes and the cluster list threw on the first render after upgrading
        version: 2,
        migrate: () => InitialStore,
        // NOTE: nodeState is nested, so the default shallow merge would drop any
        //  field added to it later that is missing from an older persisted blob.
        merge: (persisted, current) => {
            const state = persisted as Partial<Store> | undefined
            return {...current, ...state, nodeState: {...current.nodeState, ...state?.nodeState}}
        },
    },
))

export const useStoreAction = {
    setCluster: setCluster,
    setClusterKeeperPlugin: setClusterKeeperPlugin,
    setSearchCluster: setClusterSearch,
    setClusterKeeper: setClusterKeeper,
    setWarnings: setWarnings,
    setNode: setNode,
    setTags: setTags,
    clear: clear,
    setConsoleQuery: setConsoleQuery,
    setNodeBody: setNodeBody,
    setQueryTab: setQueryTab,
    setDbName: setDbName,
    setDbSchema: setDbSchema,
    setRefreshPeriod: setRefreshPeriod,
    setContainerTab: setContainerTab,
    setSystemTab: setSystemTab,
    setSystemLogsPath: setSystemLogsPath,
    setToolsTab: setToolsTab,
}

// SETTERS
function setClusterSearch(search: string) {
    useStore.setState(s => ({...s, searchCluster: search}))
}

function setClusterKeeperPlugin(plugin: KeeperPlugin) {
    useStore.setState(s => ({...s, activeClusterKeeperPlugin: plugin, activeCluster: undefined, manualKeeper: undefined}))
}

function setCluster(cluster?: Cluster) {
    useStore.setState(s => ({...s, activeCluster: cluster, manualKeeper: undefined}))
}

function setClusterKeeper(manualKeeper?: string) {
    useStore.setState(s => {
        if (!s.activeCluster) return s
        return {...s, manualKeeper: manualKeeper}
    })
}

function setWarnings(name: string, warning: boolean) {
    useStore.setState(s => ({...s, warnings: {...s.warnings, [name]: warning}}))
}

function setNode(node?: string) {
    useStore.setState(s => {
        const clusterName = s.activeCluster?.name
        if (!clusterName) return s
        if (node) return {...s, activeNode: {...s.activeNode, [clusterName]: node}}
        if (!s.activeNode[clusterName]) return s
        const store = {...s, activeNode: {...s.activeNode}}
        delete store.activeNode[clusterName]
        return store
    })
}

function setTags(tags: string[]) {
    useStore.setState(s => ({...s, activeTags: tags}))
}

function clear() {
    MainQueryClient.clear()
    useStore.setState(useStore.getInitialState)
}

function setConsoleQuery(q: string) {
    useStore.setState(s => ({...s, nodeState: {...s.nodeState, queryConsole: q}}))
}

function setNodeBody(t: NodeTabType) {
    useStore.setState(s => ({...s, nodeState: {...s.nodeState, nodeTab: t}}))
}

function setQueryTab(t: QueryType) {
    useStore.setState(s => ({...s, nodeState: {...s.nodeState, queryTab: t}}))
}

function setDbName(n?: string) {
    useStore.setState(s => ({...s, nodeState: {...s.nodeState, dbName: n}}))
}

function setDbSchema(n?: string) {
    useStore.setState(s => ({...s, nodeState: {...s.nodeState, dbSchema: n}}))
}

function setRefreshPeriod(key: string, period: [string, number]) {
    useStore.setState(s => ({...s, refresh: {...s.refresh, [key]: period}}))
}

function setContainerTab(t: number) {
    useStore.setState(s => ({...s, nodeState: {...s.nodeState, containerTab: t}}))
}

function setSystemTab(t: number) {
    useStore.setState(s => ({...s, nodeState: {...s.nodeState, systemTab: t}}))
}

function setSystemLogsPath(path: string) {
    useStore.setState(s => ({...s, nodeState: {...s.nodeState, systemLogsPath: path}}))
}

function setToolsTab(t: number) {
    useStore.setState(s => ({...s, nodeState: {...s.nodeState, toolsTab: t}}))
}



