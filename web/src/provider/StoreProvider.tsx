import {createContext, ReactNode, useContext, useState} from "react";

interface StoreType {
    activeCluster: string,
    activeNode: string
}

interface StoreTypeParam {
    activeCluster?: string
    activeNode?: string
}

interface StoreContextType {
    store: StoreType
    setStore: (store: StoreTypeParam) => void
    isClusterActive: (name: string) => boolean
    isNodeActive: (name: string) => boolean
}

const initialStore: StoreType = {
    activeCluster: '',
    activeNode: '',
}

const StoreContext = createContext<StoreContextType>({
    store: initialStore, setStore: () => {},
    isClusterActive: () => false, isNodeActive: () => false,
})

export function useStore() {
    return useContext(StoreContext)
}

export function StoreProvider(props: { children: ReactNode }) {
    const [state, setState] = useState(initialStore)

    const setStore = (store: StoreTypeParam) => setState({ ...state, ...store })
    const isClusterActive = (name: string) => name === state.activeCluster
    const isNodeActive = (name: string) => name === state.activeNode

    return (
        <StoreContext.Provider value={{ store: state, setStore, isClusterActive, isNodeActive }}>
            {props.children}
        </StoreContext.Provider>
    )
}
