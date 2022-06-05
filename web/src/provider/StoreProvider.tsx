import {createContext, ReactNode, useContext, useState} from "react";
import {Node} from "../app/types";

interface StoreType {
    activeCluster: { name: string, leader?: Node, tab: number }
    activeNode: string
    credentialsOpen: boolean
}

interface StoreTypeParam {
    activeCluster?: { name: string, leader?: Node, tab: number }
    activeNode?: string
    credentialsOpen?: boolean
}

interface StoreContextType {
    store: StoreType
    setStore: (store: StoreTypeParam) => void
    isClusterActive: (name: string) => boolean
    isNodeActive: (name: string) => boolean
    isOverviewOpen: () => boolean
}

const initialStore: StoreType = {
    activeCluster: { name: "", leader: undefined, tab: 0 },
    activeNode: "",
    credentialsOpen: false,
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
        isClusterActive: (name: string) => name === activeCluster.name,
        isNodeActive: (name: string) => name === activeNode,
        isOverviewOpen: () => !!activeCluster.leader && activeCluster.tab === 0
    }
    return (
        <StoreContext.Provider value={value}>
            {props.children}
        </StoreContext.Provider>
    )
}
