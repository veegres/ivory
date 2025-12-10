import {persist} from "zustand/middleware"
import {create} from "zustand/react"

import {ActiveCluster, ActiveInstance, DetectionType} from "../api/cluster/type"
import {InstanceTabType, InstanceWeb} from "../api/instance/type"
import {QueryType} from "../api/query/type"
import {MainQueryClient} from "./AppProvider"

// STORE
interface Store {
    activeClusterTab: number,
    activeCluster?: ActiveCluster,
    activeInstance: ActiveInstance,
    activeTags: string[],
    warnings: { [key: string]: boolean },
    settings: boolean,
    instance: {
        body: InstanceTabType,
        queryTab: QueryType,
        queryConsole: string,
        dbName?: string,
    },
}

export const useStore = create(persist<Store>(
    () => ({
        activeClusterTab: 0,
        activeCluster: undefined,
        activeInstance: {},
        activeTags: ["ALL"],
        warnings: {},
        settings: false,
        instance: {
            body: InstanceTabType.CHART,
            queryTab: QueryType.CONSOLE,
            queryConsole: "",
            dbName: undefined,
        },
    }),
    {name: "store", version: 1}
))

export const useStoreAction = {
    setCluster: setCluster,
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
    useStore.setState(s => {
        const clusterName = s.activeCluster?.cluster.name
        if (!clusterName) return s
        if (instance) return {...s, activeInstance: {...s.activeInstance, [clusterName]: instance}}
        if (!s.activeInstance[clusterName]) return s
        const store = {...s}
        delete store.activeInstance[clusterName]
        return store
    })
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


