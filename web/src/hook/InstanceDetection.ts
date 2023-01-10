import {Cluster, InstanceLocal, InstanceDetection, InstanceMap} from "../app/types";
import {useQueries, useQuery} from "@tanstack/react-query";
import {instanceApi} from "../app/api";
import {useEffect, useMemo, useRef, useState} from "react";
import {combineInstances, createInstanceColors, getDomain, getHostAndPort} from "../app/utils";

type ManualState = {
    instance: InstanceLocal,
    instances: InstanceMap,
}

export function useManualInstanceDetection(use: boolean, cluster: Cluster, state: ManualState): InstanceDetection {
    const { instance: selected } = state

    const query = useQuery(
        ["manual", "instance/overview", cluster.name, selected.sidecar.host, selected.sidecar.port],
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

// TODO consider moving this to backend
export function useAutoInstanceDetection(use: boolean, cluster: Cluster): InstanceDetection {
    const [index, setIndex] = useState(0)
    const activeNodeName = cluster.nodes[index]
    const safeNodeName = useRef(activeNodeName)

    // we need disable cause it thinks that function can be updated
    // eslint-disable-next-line react-hooks/exhaustive-deps
    const instanceQueries = useMemo(() => getNodeQueries(index, cluster.name, cluster.nodes), [index, cluster.name, cluster.nodes])
    const queries = useQueries({ queries: instanceQueries })
    const query = useMemo(() => queries[index] ?? {}, [queries, index])

    const clusterInstances = useMemo(() => query.data ?? {}, [query.data])
    const colors = useMemo(() => createInstanceColors(clusterInstances), [clusterInstances])
    const [instances, warning] = useMemo(() => combineInstances(cluster.nodes, clusterInstances), [cluster.nodes, clusterInstances])

    const instance = useMemo(() => handleMemoActiveInstance(instances, activeNodeName), [instances, activeNodeName])

    // TODO it helps to refetch current instance query by using just instance sidecar info, but when we don't have leader it points doesn't right query...
    useEffect(
        () => { cluster.nodes.forEach((node, i) => {
            if (node === getDomain(instance.sidecar)) setIndex(i)
        })},
        [cluster.nodes, instance.sidecar]
    )

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
     * Either find leader or set query that we send request to
     */
    function handleMemoActiveInstance(instances: InstanceMap, name: string) {
        const values = Object.values(instances)
        const value = values.find(instance => instance.leader) ?? instances[name] ?? instances[safeNodeName.current]
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
                queryKey: ["auto", "instance/overview", name, domain.host, domain.port],
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
