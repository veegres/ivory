import {persist} from "zustand/middleware"
import {create} from "zustand/react"

import {Cluster} from "../../features/cluster/api/type"
import {NodeTabType} from "../../features/node/api/type"
import {Type as QueryType} from "../../features/query/api/type"
import {MainQueryClient} from "./AppProvider"

// STORE
interface Store {
    searchCluster: string,
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
    },
}

export const useStore = create(persist<Store>(
    () => ({
        searchCluster: "",
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
        },
    }),
    {name: "store", version: 1},
))

export const useStoreAction = {
    setCluster: setCluster,
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
}

// SETTERS
function setClusterSearch(search: string) {
    useStore.setState(s => ({...s, searchCluster: search}))
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
        delete storage.activeNode[clusterName]
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



