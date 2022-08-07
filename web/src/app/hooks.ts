import {useEffect, useState} from "react";
import {EventStream, JobStatus} from "./types";
import {useQueries, useQuery} from "react-query";
import {bloatApi, nodeApi} from "./api";
import {JobOptions} from "./utils";

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


export function useSmartClusterQuery(name: string, nodes: string[]) {
    const [nodeIndex, setNodeIndex] = useState(0)
    const [nodeQueries, setNodeQueries] = useState(getNodeQueries(nodeIndex, nodes))
    const [isRefetching, setIsRefetching] = useState(false)

    const result = useQueries(nodeQueries)

    useEffect(handleEffectIncrease, [nodeIndex, result, isRefetching])

    return {
        instance: nodes[nodeIndex],
        instanceResult: result[nodeIndex] ?? {},
        update: (nodes: string[]) => {
            setNodeIndex(0)
            setNodeQueries(getNodeQueries(nodeIndex, nodes))
        },
        refetch: () => {
            setNodeIndex(0)
            setIsRefetching(true)
        }
    }

    function handleEffectIncrease() {
        const node = result[nodeIndex] ?? {}
        if (node.isError && nodeIndex < result.length - 1) setNodeIndex(nodeIndex + 1)
        if (node.isIdle || (node.isError && isRefetching)) node.refetch()
        if (node.isSuccess || nodeIndex === result.length - 1) setIsRefetching(false)
    }

    /**
     * Create array of queries that should be requested be `useQueries()` in new rerender
     * @param currentIndex
     * @param nodes
     */
    function getNodeQueries(currentIndex: number, nodes: string[]) {
        return nodes.map((node) => ({
            queryKey: [`node/cluster`, node],
            queryFn: () => nodeApi.cluster(node),
            retry: 0, enabled: false
        }))
    }
}
