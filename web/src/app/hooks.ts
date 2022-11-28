import {useEffect, useMemo, useState} from "react";
import {Cluster, EventStream, Instance, InstanceDetection, InstanceMap, JobStatus} from "./types";
import {QueryKey, useQueries, useQuery, useQueryClient} from "@tanstack/react-query";
import {bloatApi, nodeApi} from "./api";
import {combineInstances, createInstanceColors, getErrorMessage, JobOptions} from "./utils";
import {useSnackbar} from "notistack";

export interface EventJob {
    isFetching: boolean;
    logs: string[];
    status: { name: string; color: string; active: boolean }
}

export function useEventJob(uuid: string, initStatus: JobStatus, isOpen: boolean): EventJob {
    const [logs, setLogs] = useState<string[]>([])
    const [status, setStatus] = useState(JobOptions[initStatus])
    const [isEventSourceFetching, setEventSourceFetching] = useState<boolean>(false)

    const {isFetching} = useQuery(["node/bloat/logs", uuid], () => bloatApi.logs(uuid), {
        onSuccess: data => setLogs(data),
        refetchOnReconnect: false,
        refetchOnWindowFocus: false,
        enabled: !status.active && isOpen
    })

    useEffect(() => {
        setStatus(JobOptions[initStatus])
        if (!JobOptions[initStatus].active) return

        const es = bloatApi.stream(uuid)
        const close = () => {
            es.close();
            setEventSourceFetching(false)
        }
        const addLog = (log: string) => setLogs((old) => [...old, log])
        es.onopen = () => {
            setLogs([])
            setEventSourceFetching(true)
            addLog("Logs streaming open: New connection was established")
        }
        es.addEventListener("log", (e: MessageEvent<string>) => addLog(e.data))
        es.addEventListener("server", (e: MessageEvent<string>) => addLog(e.data))
        es.addEventListener("status", (e: MessageEvent<JobStatus>) => setStatus(JobOptions[e.data]))
        es.addEventListener("stream", (e: MessageEvent<EventStream>) => {
            if (e.data === EventStream.END) close()
        })
        es.onerror = () => {
            setEventSourceFetching(false);
            addLog("Logs streaming error: Trying to reestablish connection")
        }

        return () => es.close()
    }, [uuid, initStatus])

    return {
        isFetching: isEventSourceFetching || isFetching,
        logs,
        status
    }
}


export function useManualInstanceDetection(use: boolean, cluster: Cluster, selected: Instance): InstanceDetection {
    const query = useQuery(
        ["node/cluster", cluster.name, selected?.api_domain],
        () => { if (selected) return nodeApi.cluster(selected.api_domain) },
        {enabled: use && !!selected}
    )

    const clusterInstances = useMemo(() => query.data ?? {}, [query.data])
    const colors = useMemo(() => createInstanceColors(clusterInstances), [clusterInstances])
    const [instances, warning] = useMemo(() => combineInstances(cluster.nodes, clusterInstances), [cluster.nodes, clusterInstances])
    if (!instances[selected.api_domain]) instances[selected.api_domain] = selected
    const instance = instances[selected.api_domain]

    return {
        active: { cluster, instance, instances, warning },
        colors,
        fetching: query.isFetching,
        refetch: query.refetch,
    }
}

export function useAutoInstanceDetection(use: boolean, cluster: Cluster): InstanceDetection {
    const [index, setIndex] = useState(0)
    const activeNodeName = cluster.nodes[index]

    // we need disable cause it thinks that function can be updated
    // eslint-disable-next-line react-hooks/exhaustive-deps
    const instanceQueries = useMemo(() => getNodeQueries(index, cluster.name, cluster.nodes), [index, cluster.name, cluster.nodes])
    const queries = useQueries({ queries: instanceQueries })
    const query = queries[index] ?? {}

    const clusterInstances = useMemo(() => query.data ?? {}, [query.data])
    const colors = useMemo(() => createInstanceColors(clusterInstances), [clusterInstances])
    const [instances, warning] = useMemo(() => combineInstances(cluster.nodes, clusterInstances), [cluster.nodes, clusterInstances])

    const instance = useMemo(() => handleMemoActiveInstance(instances, activeNodeName), [instances, activeNodeName])

    return {
        active: { cluster, instance, instances, warning },
        colors,
        fetching: query.isFetching,
        refetch: () => setIndex(0)
    }

    /**
     * Either find leader or set query that we were sending request to
     */
    function handleMemoActiveInstance(instances: InstanceMap, name: string) {
        const values = Object.values(instances)
        return values.find(instance => instance.leader) ?? instances[name]
    }

    /**
     * Create array of queries that should be requested be `useQueries()` in new rerender
     * @param index
     * @param name
     * @param instances
     */
    function getNodeQueries(index: number, name: string, instances: string[]) {
        return instances.map((instance, j) => ({
            queryKey: ["node/cluster", name, instance],
            queryFn: () => nodeApi.cluster(instance),
            retry: 0,
            enabled: use && index === j,
            onError: () => {
                if (index < instances.length - 1) setIndex(index => index + 1)
            }
        }))
    }
}

/**
 * Simplify handling `onSuccess` and `onError` requests for react-query client
 * providing common approach with request refetch and custom toast messages for
 * mutation requests
 *
 * @param queryKey
 */
export function useMutationOptions(queryKey?: QueryKey) {
    const { enqueueSnackbar } = useSnackbar()
    const queryClient = useQueryClient();

    return {
        queryClient,
        onSuccess: queryKey === undefined ? undefined : () => queryClient.refetchQueries(queryKey),
        onError: (error: any) => enqueueSnackbar(getErrorMessage(error), { variant: "error" }),
    }
}

