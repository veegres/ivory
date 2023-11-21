import {createContext, ReactNode, useContext, useMemo} from "react";
import {useQueryClient} from "@tanstack/react-query";
import {ActiveCluster, Cluster, DetectionType} from "../type/cluster";
import {InstanceMap, InstanceWeb} from "../type/instance";
import {getDomain} from "../app/utils";
import {useLocalStorageState} from "../hook/LocalStorage";

// STORE
interface StoreType {
    activeCluster?: ActiveCluster, activeClusterTab: number,
    activeInstance?: InstanceWeb,
    activeTags: string[],
    warnings: { [key: string]: boolean },
    settings: boolean,

    isClusterActive: (name: string) => boolean
    isClusterOverviewOpen: () => boolean
    isInstanceActive: (key: string) => boolean,
}
const initialStore: StoreType = {
    activeCluster: undefined, activeClusterTab: 0,
    activeInstance: undefined,
    activeTags: ["ALL"],
    warnings: {},
    settings: false,

    isClusterActive: () => false,
    isClusterOverviewOpen: () => false,
    isInstanceActive: () => false,
}

// ACTIONS
interface StoreActionType {
    setCluster: (cluster?: ActiveCluster) => void,
    setClusterInfo: (cluster: Cluster, defaultInstance: InstanceWeb, combinedInstanceMap: InstanceMap, warning: boolean) => void,
    setClusterInstance: (instance: InstanceWeb) => void,
    setClusterDetection: (detection: DetectionType) => void
    setClusterTab: (tab: number) => void,
    setWarnings: (name: string, warning: boolean) => void,
    setInstance: (instance?: InstanceWeb) => void,
    setTags: (tags: string[]) => void,
    toggleSettingsDialog: () => void,
    clear: () => void,
}
const initialStoreAction: StoreActionType = {
    setCluster: () => void 0,
    setClusterInfo: () => void 0,
    setClusterInstance: () => void 0,
    setClusterDetection: () => void 0,
    setClusterTab: () => void 0,
    setWarnings: () => void 0,
    setInstance: () => void 0,
    setTags: () => void 0,
    toggleSettingsDialog: () => void 0,
    clear: () => void 0,
}


// CREATE STORE CONTEXT
const StoreContext = createContext(initialStore)
const StoreActionContext = createContext(initialStoreAction)

export function useStore() {
    return useContext(StoreContext)
}

export function useStoreAction() {
    return useContext(StoreActionContext)
}

export function StoreProvider(props: { children: ReactNode }) {
    const [state, setState] = useLocalStorageState("store", initialStore)
    const queryClient = useQueryClient()

    const getters = useMemo(getGetters, [state])
    const setters = useMemo(getSetters, [])

    // we need two providers, because otherwise each team when you use only actions in some component
    // it will rerenders as well, event if you don't use any state inside
    return (
        <StoreContext.Provider value={getters}>
            <StoreActionContext.Provider value={setters}>
                {props.children}
            </StoreActionContext.Provider>
        </StoreContext.Provider>
    )

    function getGetters(): StoreType {
        return {
            ...state,
            isClusterActive, isClusterOverviewOpen, isInstanceActive,
        }
    }

    function getSetters(): StoreActionType {
        return {
            setCluster, setClusterInfo, setClusterInstance,
            setClusterDetection, setClusterTab, setWarnings,
            setInstance, setTags, toggleSettingsDialog, clear,
        }
    }

    // SETTERS
    function setCluster(cluster?: ActiveCluster) {
        setState(s => ({...s, activeCluster: cluster}))
    }
    function setClusterInfo(cluster: Cluster, defaultInstance: InstanceWeb, combinedInstanceMap: InstanceMap, warning: boolean) {
        setState(s => {
            if (!s.activeCluster) return s
            return ({...s, activeCluster: {...s.activeCluster, cluster, defaultInstance, combinedInstanceMap, warning}})
        })
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
