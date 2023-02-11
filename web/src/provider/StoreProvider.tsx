import {createContext, ReactNode, useContext, useState} from "react";
import {ActiveCluster, ActiveInstance, DetectionType, InstanceLocal} from "../app/types";
import {useQueryClient} from "@tanstack/react-query";

// STORE
interface StoreType {
    activeCluster?: ActiveCluster, activeClusterTab: number,
    activeInstance?: ActiveInstance, activeInstanceTab: number,
    activeTags: string[],
    settings: boolean,
}
const initialStore: StoreType = {
    activeCluster: undefined, activeClusterTab: 0,
    activeInstance: undefined, activeInstanceTab: 0,
    activeTags: ["ALL"],
    settings: false,
}

// CONTEXT
interface StoreContextType {
    store: StoreType

    setCluster: (cluster?: ActiveCluster) => void,
    setClusterInstance: (instance: InstanceLocal) => void,
    setClusterDetection: (detection: DetectionType) => void
    setClusterTab: (tab: number) => void,
    isClusterActive: (name: string) => boolean
    isClusterOverviewOpen: () => boolean

    setInstance: (instance?: ActiveInstance) => void,
    setInstanceTab: (tab: number) => void,
    isInstanceActive: (instance?: ActiveInstance) => boolean,

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
    isClusterActive: () => false,
    isClusterOverviewOpen: () => false,

    setInstance: () => void 0,
    setInstanceTab: () => void 0,
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
    const [state, setState] = useState(initialStore)
    const queryClient = useQueryClient();

    const value = getStoreContext()
    return (
        <StoreContext.Provider value={value}>
            {props.children}
        </StoreContext.Provider>
    )

    function getStoreContext(): StoreContextType {
        const { activeCluster, activeInstance } = state

        return {
            store: state,

            setCluster: (cluster?: ActiveCluster) => setState(state => ({...state, activeCluster: cluster, activeInstance: undefined})),
            setClusterInstance: (instance: InstanceLocal) => {
                if (activeCluster) setState({...state, activeCluster: {...activeCluster, instance, detection: "manual"}})
            },
            setClusterDetection: (detection: DetectionType) => {
                if (activeCluster) setState({...state, activeCluster: {...activeCluster, detection}})
            },
            setClusterTab: (tab: number) => setState({...state, activeClusterTab: tab}),
            isClusterActive: (name: string) => name === activeCluster?.cluster.name,
            isClusterOverviewOpen: () => !!activeCluster && state.activeClusterTab === 0,

            setInstance: (instance?: ActiveInstance) => setState({...state, activeInstance: instance}),
            setInstanceTab: (tab: number) => setState({...state, activeInstanceTab: tab}),
            isInstanceActive: (instance?: ActiveInstance) => instance === activeInstance,

            setTags: (tags: string[]) => setState({...state, activeTags: tags}),
            toggleSettingsDialog: () => setState({...state, settings: !state.settings}),

            clear: () => {
                queryClient.clear()
                setState(initialStore)
            }
        }
    }
}
