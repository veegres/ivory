import {useEffect, useState} from "react";
import {JobOptions} from "../app/utils";
import {useQuery} from "@tanstack/react-query";
import {bloatApi} from "../app/api";
import {EventStreamType, EventType, JobStatus} from "../type/job";

type EventJob = {
    isFetching: boolean;
    logs: string[];
    status: { name: string; color: string; active: boolean }
}

export function useEventJob(uuid: string, initStatus: JobStatus, isOpen: boolean): EventJob {
    const [logs, setLogs] = useState<string[]>([])
    const [status, setStatus] = useState(JobOptions[initStatus])
    const [isEventSourceFetching, setEventSourceFetching] = useState<boolean>(false)

    const {isFetching} = useQuery(
        ["instance/bloat/logs", uuid],
        () => bloatApi.logs(uuid),
        {
            onSuccess: (data) => setLogs(data),
            refetchOnReconnect: false,
            refetchOnWindowFocus: false,
            enabled: !status.active && isOpen
        }
    )

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
        es.addEventListener(EventType.LOG, (e: MessageEvent<string>) => addLog(e.data))
        es.addEventListener(EventType.SERVER, (e: MessageEvent<string>) => addLog(e.data))
        es.addEventListener(EventType.STATUS, (e: MessageEvent<JobStatus>) => setStatus(JobOptions[e.data]))
        es.addEventListener(EventType.STREAM, (e: MessageEvent<EventStreamType>) => {
            if (e.data === EventStreamType.END) close()
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
