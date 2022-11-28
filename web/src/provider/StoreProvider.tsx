import {createContext, ReactNode, useContext, useState} from "react";
import {ActiveCluster, DetectionType, Instance} from "../app/types";

// STORE
interface StoreType {
    activeCluster?: ActiveCluster, activeClusterTab: number,
    activeInstance?: string, activeInstanceTab: number,
    credentialsOpen: boolean,
    certsOpen: boolean,
}
const initialStore: StoreType = {
    activeCluster: undefined, activeClusterTab: 0,
    activeInstance: undefined, activeInstanceTab: 0,
    credentialsOpen: false, certsOpen: false,
}

// CONTEXT
interface StoreContextType {
    store: StoreType

    setCluster: (cluster?: ActiveCluster) => void,
    setClusterInstance: (instance: Instance) => void,
    setClusterDetection: (detection: DetectionType) => void
    setClusterTab: (tab: number) => void,
    isClusterActive: (name: string) => boolean
    isClusterOverviewOpen: () => boolean

    setInstance: (instance?: string) => void,
    setInstanceTab: (tab: number) => void,
    isInstanceActive: (name: string) => boolean,

    toggleCredentialsWindow: () => void,
    toggleCertsWindow: () => void,
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

    toggleCredentialsWindow: () => void 0,
    toggleCertsWindow: () => void 0,
}


// CREATE STORE CONTEXT
const StoreContext = createContext(initialStoreContext)

export function useStore() {
    return useContext(StoreContext)
}

export function StoreProvider(props: { children: ReactNode }) {
    const [state, setState] = useState(initialStore)
    const { activeCluster, activeInstance } = state

    const value: StoreContextType = {
        store: state,

        setCluster: (cluster?: ActiveCluster) => setState(state => ({...state, activeCluster: cluster, activeInstance: undefined})),
        setClusterInstance: (instance: Instance) => {
            if (activeCluster) setState({...state, activeCluster: {...activeCluster, instance, detection: "manual"}})
        },
        setClusterDetection: (detection: DetectionType) => {
            if (activeCluster) setState({...state, activeCluster: {...activeCluster, detection}})
        },
        setClusterTab: (tab: number) => setState({...state, activeClusterTab: tab}),
        isClusterActive: (name: string) => name === activeCluster?.cluster.name,
        isClusterOverviewOpen: () => !!activeCluster && state.activeClusterTab === 0,

        setInstance: (instance?: string) => setState({...state, activeInstance: instance}),
        setInstanceTab: (tab: number) => setState({...state, activeInstanceTab: tab}),
        isInstanceActive: (name: string) => name === activeInstance,

        toggleCredentialsWindow: () => setState({...state, credentialsOpen: !state.credentialsOpen}),
        toggleCertsWindow: () => setState({...state, certsOpen: !state.certsOpen}),
    }
    return (
        <StoreContext.Provider value={value}>
            {props.children}
        </StoreContext.Provider>
    )
}
