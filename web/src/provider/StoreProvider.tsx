import {createContext, ReactNode, useContext, useState} from "react";
import {Cluster, Instance} from "../app/types";

// ACTIVE CLUSTER
interface ActiveClusterType { cluster?: Cluster, instance?: Instance, tab: number }
export const initialActiveCluster: ActiveClusterType = { cluster: undefined, instance: undefined, tab: 0 }

// STORE
interface StoreType { activeCluster: ActiveClusterType, activeNode: string, credentialsOpen: boolean }
interface StoreTypeParam { activeCluster?: ActiveClusterType, activeNode?: string, credentialsOpen?: boolean }
export const initialStore: StoreType = { activeCluster: initialActiveCluster, activeNode: "", credentialsOpen: false }

interface StoreContextType {
    store: StoreType
    setStore: (store: StoreTypeParam) => void
    isClusterActive: (name: string) => boolean
    isNodeActive: (name: string) => boolean
    isOverviewOpen: () => boolean
}

const StoreContext = createContext<StoreContextType>({
    store: initialStore, setStore: () => {},
    isClusterActive: () => false, isNodeActive: () => false,
    isOverviewOpen: () => false,
})

export function useStore() {
    return useContext(StoreContext)
}

export function StoreProvider(props: { children: ReactNode }) {
    const [state, setState] = useState(initialStore)
    const { activeCluster, activeNode } = state
    const value = {
        store: state,
        setStore: (store: StoreTypeParam) => setState({ ...state, ...store }),
        isClusterActive: (name: string) => name === activeCluster.cluster?.name,
        isNodeActive: (name: string) => name === activeNode,
        isOverviewOpen: () => !!activeCluster.instance && activeCluster.tab === 0
    }
    return (
        <StoreContext.Provider value={value}>
            {props.children}
        </StoreContext.Provider>
    )
}
