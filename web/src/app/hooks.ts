import {useEffect, useMemo, useState} from "react";
import {EventStream, Instance, InstanceMap, JobStatus} from "./types";
import {useQueries, useQuery} from "react-query";
import {bloatApi, nodeApi} from "./api";
import {createColorsMap, JobOptions} from "./utils";

export interface EventJob {
    isFetching: boolean;
    logs: string[];
    status: { name: string; color: string; active: boolean }
}

export function useEventJob(uuid: string, initStatus: JobStatus, isOpen: boolean): EventJob {
    const [logs, setLogs] = useState<string[]>([])
    const [status, setStatus] = useState(JobOptions[initStatus])
    const [isEventSourceFetching, setEventSourceFetching] = useState<boolean>(false)

    const {isFetching} = useQuery(['node/bloat/logs', uuid], () => bloatApi.logs(uuid), {
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

const initialInstance: Instance = { name: "-", host: "-", port: 0, role: "unknown", api_domain: "-", api_url: "-", lag: undefined, leader: false, state: "-", inInstances: true, inCluster: false }

export function useSmartClusterQuery(name: string, instanceNames: string[]) {
    const [index, setIndex] = useState(0)
    const [nodeQueries, setNodeQueries] = useState(getNodeQueries(instanceNames))
    const [isRefetching, setIsRefetching] = useState(false)

    const result = useQueries(nodeQueries)
    const instance = useMemo(() => result[index] ?? {}, [result, index])
    const instanceMap = useMemo(() => instance.data ?? {}, [instance])

    useEffect(handleEffectIncrease, [index, instance, isRefetching, result.length])

    const activeInstance = useMemo(() => handleMemoActiveInstance(index, instanceMap), [instanceMap, index])
    const colors = useMemo(() => createColorsMap(instanceMap), [instanceMap])
    const [instances, warning] = useMemo(() => handleMemoInstances(instanceMap, instanceNames), [instanceMap, instanceNames])

    return {
        instance: activeInstance,
        instances: instances,
        isFetching: instance.isFetching,
        warning: warning,
        colors,
        update: (nodes: string[]) => {
            setIndex(0)
            setNodeQueries(getNodeQueries(nodes))
        },
        refetch: () => {
            setIndex(0)
            setIsRefetching(true)
        }
    }

    /**
     * Fetching each instance while we don't have success request
     */
    function handleEffectIncrease() {
        if (instance.isError && index < result.length - 1) setIndex(index + 1)
        if (instance.isIdle || (instance.isError && isRefetching)) instance.refetch()
        if (instance.isSuccess || index === result.length - 1) setIsRefetching(false)
    }

    /**
     * Combine instances from patroni and from ivory
     */
    function handleMemoInstances(instanceMap: InstanceMap, instanceNames: string[]): [InstanceMap, boolean] {
        const map: InstanceMap = {}
        let warning: boolean = false

        for (const key in instanceMap) {
            if (instanceNames.includes(key)) {
                map[key] = { ...instanceMap[key], inInstances: true }
            } else {
                map[key] = { ...instanceMap[key], inInstances: false }
            }
        }

        for (const value of instanceNames) {
            if (!map[value]) {
                map[value] = initialInstance
            }
        }

        for (const key in map) {
            const value = map[key]
            if (!value.inInstances || !value.inCluster) {
                warning = true
            }
        }

        return [map, warning]
    }

    /**
     * Either find leader or set instance that we were sending request to
     */
    function handleMemoActiveInstance(index: number, instanceMap: InstanceMap) {
        return Object.values(instanceMap).find(instance => instance.leader) ?? instanceMap[index]
    }

    /**
     * Create array of queries that should be requested be `useQueries()` in new rerender
     * @param instances
     */
    function getNodeQueries(instances: string[]) {
        return instances.map((instance) => ({
            queryKey: [`node/cluster`, instance],
            queryFn: () => nodeApi.cluster(instance),
            retry: 0, enabled: false
        }))
    }
}
