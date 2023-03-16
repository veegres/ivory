import {createContext, ReactNode, useContext, useState} from "react";
import {useQueryClient} from "@tanstack/react-query";
import {ActiveCluster, DetectionType} from "../type/cluster";
import {Instance} from "../type/Instance";
import {getDomain} from "../app/utils";

// STORE
interface StoreType {
    activeCluster?: ActiveCluster, activeClusterTab: number,
    activeInstance?: Instance,
    activeTags: string[],
    settings: boolean,
}
const initialStore: StoreType = {
    activeCluster: undefined, activeClusterTab: 0,
    activeInstance: undefined,
    activeTags: ["ALL"],
    settings: false,
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

    setInstance: (instance?: Instance) => void,
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
    isClusterActive: () => false,
    isClusterOverviewOpen: () => false,

    setInstance: () => void 0,
    isInstanceActive: (key: string) => false,

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

            setCluster: (cluster?: ActiveCluster) => setState(state => ({...state, activeCluster: cluster})),
            setClusterInstance: (instance: Instance) => {
                if (activeCluster) setState({...state, activeCluster: {...activeCluster, defaultInstance: instance, detection: "manual"}})
            },
            setClusterDetection: (detection: DetectionType) => {
                if (activeCluster) setState({...state, activeCluster: {...activeCluster, detection}})
            },
            setClusterTab: (tab: number) => setState({...state, activeClusterTab: tab}),
            isClusterActive: (name: string) => name === activeCluster?.cluster.name,
            isClusterOverviewOpen: () => !!activeCluster && state.activeClusterTab === 0,

            setInstance: (instance?: Instance) => setState({...state, activeInstance: instance}),
            isInstanceActive: (key: string) => activeInstance ? getDomain(activeInstance.sidecar) === key : false,

            setTags: (tags: string[]) => setState({...state, activeTags: tags}),
            toggleSettingsDialog: () => setState({...state, settings: !state.settings}),

            clear: () => {
                queryClient.clear()
                setState(initialStore)
            }
        }
    }
}
