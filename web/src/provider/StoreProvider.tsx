import {ActiveCluster, Cluster, DetectionType} from "../api/cluster/type";
import {InstanceMap, InstanceTabType, InstanceWeb} from "../api/instance/type";
import {getDomain} from "../app/utils";
import {QueryType} from "../api/query/type";
import {create} from "zustand/react";
import {MainQueryClient} from "../app/main";
import {persist} from "zustand/middleware";

// STORE
interface Store {
    activeCluster?: ActiveCluster, activeClusterTab: number,
    activeInstance?: InstanceWeb,
    activeTags: string[],
    warnings: { [key: string]: boolean },
    settings: boolean,
    instance: {
        body: InstanceTabType,
        queryTab: QueryType,
        queryConsole: string,
        dbName?: string,
    },

    isClusterActive: (name: string) => boolean
    isClusterOverviewOpen: () => boolean
    isInstanceActive: (key: string) => boolean,
}

export const useStore = create(persist<Store>(
    (_, get) => ({
        activeCluster: undefined, activeClusterTab: 0,
        activeInstance: undefined,
        activeTags: ["ALL"],
        warnings: {},
        settings: false,
        instance: {
        body: InstanceTabType.CHART,
            queryTab: QueryType.CONSOLE,
            queryConsole: "",
            dbName: undefined,
        },
        isClusterActive: (name: string) => isClusterActive(get(), name),
        isClusterOverviewOpen: () => isClusterOverviewOpen(get()),
        isInstanceActive: (key: string) => isInstanceActive(get(), key),
    }),
    {name: "store", version: 1}
))

export const useStoreAction = {
    setCluster: setCluster,
    setClusterInfo: setClusterInfo,
    setClusterInstance: setClusterInstance,
    setClusterDetection: setClusterDetection,
    setClusterTab: setClusterTab,
    setWarnings: setWarnings,
    setInstance: setInstance,
    setTags: setTags,
    toggleSettingsDialog: toggleSettingsDialog,
    clear: clear,
    setConsoleQuery: setConsoleQuery,
    setInstanceBody: setInstanceBody,
    setQueryTab: setQueryTab,
    setDbName: setDbName,
}

// SETTERS
function setCluster(cluster?: ActiveCluster) {
    useStore.setState(s => ({...s, activeCluster: cluster}))
}
function setClusterInfo(cluster: Cluster, defaultInstance: InstanceWeb, combinedInstanceMap: InstanceMap, warning: boolean) {
    useStore.setState(s => {
        if (!s.activeCluster) return s
        return ({...s, activeCluster: {...s.activeCluster, cluster, defaultInstance, combinedInstanceMap, warning}})
    })
}
function setClusterInstance(instance: InstanceWeb) {
    useStore.setState(s => {
        if (!s.activeCluster) return s
        return {...s, activeCluster: {...s.activeCluster, defaultInstance: instance, detection: "manual"}}
    })
}
function setClusterDetection(detection: DetectionType) {
    useStore.setState(s => {
        if (!s.activeCluster) return s
        return {...s, activeCluster: {...s.activeCluster, detection}}
    })
}
function setClusterTab(tab: number) {
    useStore.setState(s => ({...s, activeClusterTab: tab}))
}
function setWarnings(name: string, warning: boolean) {
    useStore.setState(s => ({...s, warnings: {...s.warnings, [name]: warning}}))
}
function setInstance(instance?: InstanceWeb) {
    useStore.setState(s => ({...s, activeInstance: instance}))
}
function setTags(tags: string[]) {
    useStore.setState(s => ({...s, activeTags: tags}))
}
function toggleSettingsDialog() {
    useStore.setState(s => ({...s, settings: !s.settings}))
}
function clear() {
    MainQueryClient.clear()
    useStore.setState(useStore.getInitialState)
}
function setConsoleQuery(q: string) {
    useStore.setState(s => ({...s, instance: {...s.instance, queryConsole: q}}))
}
function setInstanceBody(t: InstanceTabType) {
    useStore.setState(s => ({...s, instance: {...s.instance, body: t}}))
}
function setQueryTab(t: QueryType) {
    useStore.setState(s => ({...s, instance: {...s.instance, queryTab: t}}))
}
function setDbName(n?: string) {
    useStore.setState(s => ({...s, instance: {...s.instance, dbName: n}}))
}

// GETTERS
function isClusterActive(state: Store, name: string) {
    return name === state.activeCluster?.cluster.name
}
function isClusterOverviewOpen(state: Store) {
    return !!state.activeCluster && state.activeClusterTab === 0
}
function isInstanceActive(state: Store, key: string) {
    return state.activeInstance ? getDomain(state.activeInstance.sidecar) === key : false
}

