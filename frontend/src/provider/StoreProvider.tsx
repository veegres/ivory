import {createContext, ReactNode, useContext, useState} from "react";

interface StoreType {
    activeNode: string
}

interface StoreContextType {
    store: StoreType
    setStore: (store: StoreType) => void
}

const initialStore: StoreType = {
    activeNode: '',
}

const StoreContext = createContext<StoreContextType>({ store: initialStore, setStore: () => {} })

export function useStore() {
    return useContext(StoreContext)
}

export function StoreProvider(props: { children: ReactNode }) {
    const [state, setState] = useState(initialStore)

    const setStore = (store: StoreType) => {
        setState({ ...state, ...store })
    }

    return (
        <StoreContext.Provider value={{ store: state, setStore }}>
            {props.children}
        </StoreContext.Provider>
    )
}
