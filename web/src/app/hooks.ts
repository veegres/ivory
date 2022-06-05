import {useEffect, useState} from "react";
import {EventStream, JobStatus} from "./types";
import {useQuery} from "react-query";
import {bloatApi} from "./api";
import {jobStatus} from "./utils";

export interface EventJob {
    isFetching: boolean;
    logs: string[];
    status: { name: string; color: string; active: boolean }
}

export function useEventJob(uuid: string, initStatus: JobStatus, isOpen: boolean): EventJob {
    const [logs, setLogs] = useState<string[]>([])
    const [status, setStatus] = useState(jobStatus[initStatus])
    const [isEventSourceFetching, setEventSourceFetching] = useState<boolean>(false)

    const {isFetching} = useQuery(['node/bloat/logs', uuid], () => bloatApi.logs(uuid), {
        onSuccess: data => setLogs(data),
        refetchOnReconnect: false,
        refetchOnWindowFocus: false,
        enabled: !status.active && isOpen
    })

    useEffect(() => {
        setStatus(jobStatus[initStatus])
        if (!jobStatus[initStatus].active) return

        const es = bloatApi.stream(uuid)
        const close = () => {
            es.close();
            setEventSourceFetching(false)
        }
        const addLog = (log: string) => setLogs((old) => [...old, log])
        es.onopen = () => setEventSourceFetching(true)
        es.addEventListener("log", (e) => addLog(e.data))
        es.addEventListener("server", (e) => addLog(e.data))
        es.addEventListener("status", (e) => setStatus(jobStatus[e.data]))
        es.addEventListener("stream", (e) => {
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
