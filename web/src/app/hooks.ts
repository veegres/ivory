import {useEffect, useState} from "react";
import {EventJob, EventStream, JobStatus} from "./types";
import {isJobEnded} from "./utils";
import {useQuery} from "react-query";
import {bloatApi} from "./api";

export function useEventJob(uuid: string, initStatus: JobStatus, isOpen: boolean): EventJob {
    const [logs, setLogs] = useState<string[]>([])
    const [status, setStatus] = useState<JobStatus>(initStatus)
    const [isEventSourceFetching, setEventSourceFetching] = useState<boolean>(false)

    const {isFetching} = useQuery(['node/bloat/logs', uuid], () => bloatApi.getLogs(uuid), {
        onSuccess: data => setLogs(data),
        enabled: isJobEnded(initStatus) && isOpen && logs.length === 0
    })

    useEffect(() => {
        setStatus(initStatus)
        if (isJobEnded(initStatus)) return

        const es = new EventSource(`/api/cli/bloat/${uuid}/stream`)
        const close = () => {
            es.close();
            setEventSourceFetching(false)
        }
        const addLog = (log: string) => setLogs((old) => [...old, log])
        es.onopen = () => setEventSourceFetching(true)
        es.addEventListener("log", (e) => addLog(e.data))
        es.addEventListener("server", (e) => addLog(e.data))
        es.addEventListener("status", (e) => setStatus(e.data))
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
