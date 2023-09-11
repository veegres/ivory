import {createContext, ReactNode, useCallback, useContext} from "react";
import {useQueryClient} from "@tanstack/react-query";
import {ActiveCluster, DetectionType} from "../type/cluster";
import {InstanceWeb} from "../type/instance";
import {getDomain} from "../app/utils";
import {useLocalStorageState} from "../hook/LocalStorage";

// STORE
interface StoreType {
    activeCluster?: ActiveCluster, activeClusterTab: number,
    activeInstance?: InstanceWeb,
    activeTags: string[],
    warnings: { [key: string]: boolean },
    settings: boolean,
}
const initialStore: StoreType = {
    activeCluster: undefined, activeClusterTab: 0,
    activeInstance: undefined,
    activeTags: ["ALL"],
    warnings: {},
    settings: false,
}

// CONTEXT
interface StoreContextType {
    store: StoreType

    setCluster: (cluster?: ActiveCluster) => void,
    setClusterInstance: (instance: InstanceWeb) => void,
    setClusterDetection: (detection: DetectionType) => void
    setClusterTab: (tab: number) => void,
    setWarnings: (name: string, warning: boolean) => void,
    isClusterActive: (name: string) => boolean
    isClusterOverviewOpen: () => boolean

    setInstance: (instance?: InstanceWeb) => void,
    isInstanceActive: (key: string) => boolean,

    setTags: (tags: string[]) => void,
    toggleSettingsDialog: () => void,

    clear: () => void,
}
const initialStoreContext: StoreContextType = {
    store: initialStore,

    setCluster: () => void 0,
    setClusterInstance: () => void 0,
    setClusterDetection: () => void 0,
    setClusterTab: () => void 0,
    setWarnings: () => void 0,
    isClusterActive: () => false,
    isClusterOverviewOpen: () => false,

    setInstance: () => void 0,
    isInstanceActive: () => false,

    setTags: () => void 0,
    toggleSettingsDialog: () => void 0,

    clear: () => void 0,
}


// CREATE STORE CONTEXT
const StoreContext = createContext(initialStoreContext)

export function useStore() {
    return useContext(StoreContext)
}

export function StoreProvider(props: { children: ReactNode }) {
    const [state, setState] = useLocalStorageState("store", initialStore)
    const queryClient = useQueryClient()

    // setters
    const setClusterCallback = useCallback(setCluster, [setState])
    const setClusterInstanceCallback = useCallback(setClusterInstance, [setState])
    const setClusterDetectionCallback = useCallback(setClusterDetection, [setState])
    const setClusterTabCallback = useCallback(setClusterTab, [setState])
    const setWarningsCallback = useCallback(setWarnings, [setState])
    const setInstanceCallback = useCallback(setInstance, [setState])
    const setTagsCallback = useCallback(setTags, [setState])
    const toggleSettingsDialogCallback = useCallback(toggleSettingsDialog, [setState])
    const clearCallback = useCallback(clear, [setState, queryClient])

    // getters
    const isClusterActiveMemo = useCallback(isClusterActive, [state.activeCluster])
    const isClusterOverviewOpenMemo = useCallback(isClusterOverviewOpen, [state])
    const isInstanceActiveMemo = useCallback(isInstanceActive, [state.activeInstance])

    const value = getStoreContext()
    return (
        <StoreContext.Provider value={value}>
            {props.children}
        </StoreContext.Provider>
    )

    function getStoreContext(): StoreContextType {
        return {
            store: state,

            setCluster: setClusterCallback,
            setClusterInstance: setClusterInstanceCallback,
            setClusterDetection: setClusterDetectionCallback,
            setClusterTab: setClusterTabCallback,
            setWarnings: setWarningsCallback,
            setInstance: setInstanceCallback,
            setTags: setTagsCallback,
            toggleSettingsDialog: toggleSettingsDialogCallback,
            clear: clearCallback,

            isClusterActive: isClusterActiveMemo,
            isClusterOverviewOpen: isClusterOverviewOpenMemo,
            isInstanceActive: isInstanceActiveMemo,
        }
    }

    // SETTERS
    function setCluster(cluster?: ActiveCluster) {
        setState(s => ({...s, activeCluster: cluster}))
    }
    function setClusterInstance(instance: InstanceWeb) {
        setState(s => {
            if (!s.activeCluster) return s
            return {...s, activeCluster: {...s.activeCluster, defaultInstance: instance, detection: "manual"}}
        })
    }
    function setClusterDetection(detection: DetectionType) {
        setState(s => {
            if (!s.activeCluster) return s
            return {...s, activeCluster: {...s.activeCluster, detection}}
        })
    }
    function setClusterTab(tab: number) {
        setState(s => ({...s, activeClusterTab: tab}))
    }
    function setWarnings(name: string, warning: boolean) {
        setState(s => ({...s, warnings: {...s.warnings, [name]: warning}}))
    }
    function setInstance(instance?: InstanceWeb) {
        setState(s => ({...s, activeInstance: instance}))
    }
    function setTags(tags: string[]) {
        setState(s => ({...s, activeTags: tags}))
    }
    function toggleSettingsDialog() {
        setState(s => ({...s, settings: !s.settings}))
    }
    function clear() {
        queryClient.clear()
        setState(initialStore)
    }

    // GETTERS
    function isClusterActive(name: string) {
        return name === state.activeCluster?.cluster.name
    }
    function isClusterOverviewOpen() {
        return !!state.activeCluster && state.activeClusterTab === 0
    }
    function isInstanceActive(key: string) {
        return state.activeInstance ? getDomain(state.activeInstance.sidecar) === key : false
    }
}
