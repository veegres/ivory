import {Cluster, InstanceLocal, InstanceDetection, OverviewMap} from "../app/types";
import {useQueries, useQuery} from "@tanstack/react-query";
import {nodeApi} from "../app/api";
import {useMemo, useRef, useState} from "react";
import {combineInstances, createInstanceColors} from "../app/utils";

export function useManualInstanceDetection(use: boolean, cluster: Cluster, selected: InstanceLocal): InstanceDetection {
    const query = useQuery(
        ["instance/overview", cluster.name, selected?.sidecar.host],
        () => { if (selected) return nodeApi.overview({ host: selected.sidecar.host, port: selected.sidecar.port, cluster: cluster.name }) },
        {enabled: use && !!selected}
    )

    const clusterInstances = useMemo(() => query.data ?? {}, [query.data])
    const colors = useMemo(() => createInstanceColors(clusterInstances), [clusterInstances])
    const [instances, warning] = useMemo(() => combineInstances(cluster.nodes, clusterInstances), [cluster.nodes, clusterInstances])
    if (!instances[selected.sidecar.host]) instances[selected.sidecar.host] = selected
    const instance = instances[selected.sidecar.host]

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
     * Either find leader or set query that we were sending request to
     */
    function handleMemoActiveInstance(instances: OverviewMap, name: string) {
        const values = Object.values(instances)
        const value = values.find(instance => instance.leader) ?? instances[name] ?? instances[safeNodeName.current]
        safeNodeName.current = value.sidecar.host
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
            const [host, port] = instance.split(":")
            return ({
                queryKey: ["instance/overview", name, instance],
                queryFn: () => nodeApi.overview({ host, port: parseInt(port), cluster: cluster.name }),
                retry: 0,
                enabled: use && index === j,
                onError: () => {
                    if (index < instances.length - 1) setIndex(index => index + 1)
                }
            })
        })
    }
}
