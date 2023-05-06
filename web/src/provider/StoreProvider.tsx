import {createContext, ReactNode, useContext} from "react";
import {useQueryClient} from "@tanstack/react-query";
import {ActiveCluster, DetectionType} from "../type/cluster";
import {InstanceWeb} from "../type/Instance";
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
    const queryClient = useQueryClient();

    const value = getStoreContext()
    return (
        <StoreContext.Provider value={value}>
            {props.children}
        </StoreContext.Provider>
    )

    function getStoreContext(): StoreContextType {
        return {
            store: state,

            setCluster: (cluster?: ActiveCluster) => {
                setState(s => ({...s, activeCluster: cluster}))
            },
            setClusterInstance: (instance: InstanceWeb) => {
                setState(s => {
                    if (!s.activeCluster) return s
                    return {...s, activeCluster: {...s.activeCluster, defaultInstance: instance, detection: "manual"}}
                })
            },
            setClusterDetection: (detection: DetectionType) => {
                setState(s => {
                    if (!s.activeCluster) return s
                    return {...s, activeCluster: {...s.activeCluster, detection}}
                })
            },
            setClusterTab: (tab: number) => {
                setState(s => ({...s, activeClusterTab: tab}))
            },
            setWarnings: (name: string, warning: boolean) => {
                setState(s => ({...s, warnings: {...s.warnings, [name]: warning}}))
            },
            setInstance: (instance?: InstanceWeb) => {
                setState(s => ({...s, activeInstance: instance}))
            },
            setTags: (tags: string[]) => {
                setState(s => ({...s, activeTags: tags}))
            },
            toggleSettingsDialog: () => {
                setState(s => ({...s, settings: !s.settings}))
            },
            clear: () => {
                queryClient.clear()
                setState(initialStore)
            },


            isClusterActive: (name: string) => name === state.activeCluster?.cluster.name,
            isClusterOverviewOpen: () => !!state.activeCluster && state.activeClusterTab === 0,
            isInstanceActive: (key: string) => state.activeInstance ? getDomain(state.activeInstance.sidecar) === key : false,
        }
    }
}
