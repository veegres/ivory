import {Cluster, InstanceLocal, InstanceDetection, InstanceMap} from "../app/types";
import {useQueries, useQuery} from "@tanstack/react-query";
import {instanceApi} from "../app/api";
import {useEffect, useMemo, useRef, useState} from "react";
import {combineInstances, createInstanceColors, getDomain, getHostAndPort, initialInstance} from "../app/utils";

type ManualState = {
    instance: InstanceLocal,
    instances: InstanceMap,
}

export function useManualInstanceDetection(use: boolean, cluster: Cluster, state: ManualState): InstanceDetection {
    const { instance: selected } = state

    const query = useQuery(
        ["instance/overview", cluster.name, selected.sidecar.host, selected.sidecar.port],
        () => instanceApi.overview({ ...selected.sidecar, cluster: cluster.name }),
        {enabled: use && !!selected}
    )

    const clusterInstances = useMemo(() => query.data ?? state.instances, [query.data, state.instances])
    const colors = useMemo(() => createInstanceColors(clusterInstances), [clusterInstances])
    const [instances, warning] = useMemo(() => combineInstances(cluster.nodes, clusterInstances), [cluster.nodes, clusterInstances])

    const instance = instances[getDomain(selected.sidecar)]

    return {
        active: { cluster, instance, instances, warning },
        colors,
        fetching: query.isFetching,
        refetch: query.refetch,
    }
}

// TODO #1 consider moving this to backend
// TODO #2 investigate why it sends a lot of request after switchover or when chose manual instance or refetch (probably problem in index change `enabled: use && index === j`)
export function useAutoInstanceDetection(use: boolean, cluster: Cluster): InstanceDetection {
    const [index, setIndex] = useState(0)
    const [nodes, setNodes] = useState(cluster.nodes)
    const activeNodeName = nodes[index]
    const safeNodeName = useRef(activeNodeName)
    const safeData = useRef<InstanceMap>({})

    // we need disable cause it thinks that function can be updated
    // eslint-disable-next-line react-hooks/exhaustive-deps
    const instanceQueries = useMemo(() => getNodeQueries(index, cluster.name, nodes), [index, cluster.name, nodes])
    const queries = useQueries({ queries: instanceQueries })
    const query = useMemo(() => queries[index] ?? {}, [queries, index])

    const clusterInstances = useMemo(handleMemoClusterInstances, [query.data])
    const colors = useMemo(() => createInstanceColors(clusterInstances), [clusterInstances])
    const [instances, warning] = useMemo(() => combineInstances(cluster.nodes, clusterInstances), [cluster.nodes, clusterInstances])
    const instance = useMemo(() => handleMemoActiveInstance(instances, activeNodeName), [instances, activeNodeName])

    useEffect(handleUseEffectSetIndex, [nodes, instance.sidecar])
    useEffect(handleUseEffectSetNodes, [cluster.nodes, instances])

    return {
        active: { cluster, instance, instances, warning },
        colors,
        fetching: query.isFetching,
        refetch: () => {
            if (index === 0) query.refetch()
            else setIndex(0)
        }
    }

    /**
     * We want to update index each time when primary instance was changed (either
     * this is a missed leader or potentially something else)
     */
    function handleUseEffectSetIndex() {
        const index = nodes.findIndex(node => node === getDomain(instance.sidecar))
        if (index !== -1) setIndex(index)
    }

    /**
     * When we get combined instances from request and from cluster list,
     * we need to have local cluster instances list and add missed instances from
     * request there
     */
    function handleUseEffectSetNodes() {
        const newNodes = [...cluster.nodes]
        for (const key of Object.keys(instances)) {
            if (!newNodes.includes(key)) newNodes.push(key)
        }
        setNodes(newNodes)
    }

    /**
     * When there is new request query, we will get undefined for `query.data`
     * we don't want to lose previous data that is why we need ref `safeData`
     */
    function handleMemoClusterInstances() {
        if (!query.data) return safeData.current
        safeData.current = query.data
        return query.data
    }

    /**
     * Either find leader or set query that we send request to.
     * We need ref `safeNodeName` to be able to chose previous instance.
     */
    function handleMemoActiveInstance(instances: InstanceMap, name: string) {
        const values = Object.values(instances)
        const value = values.find(instance => instance.leader) ?? instances[name] ?? instances[safeNodeName.current] ?? initialInstance()
        safeNodeName.current = getDomain(value.sidecar)
        return value
    }

    /**
     * Create array of queries that should be requested be `useQueries()` in new rerender
     * @param index
     * @param name
     * @param instances
     */
    function getNodeQueries(index: number, name: string, instances: string[]) {
        return instances.map((instance, j) => {
            const domain = getHostAndPort(instance)
            return ({
                queryKey: ["instance/overview", name, domain.host, domain.port],
                queryFn: () => instanceApi.overview({ ...domain, cluster: cluster.name }),
                retry: 0,
                enabled: use && index === j,
                onError: () => {
                    if (index < instances.length - 1) setIndex(index => index + 1)
                }
            })
        })
    }
}
