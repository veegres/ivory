import {useEffect, useState} from "react";
import {JobOptions} from "../app/utils";
import {BloatApi} from "../app/api";
import {EventStreamType, EventType, JobStatus} from "../type/job";
import {useRouterBloatLogs} from "../router/bloat";

type EventJob = {
    isFetching: boolean;
    logs: string[];
    status: { name: string; color: string; active: boolean }
}

export function useEventJob(uuid: string, initStatus: JobStatus, isOpen: boolean): EventJob {
    const [logs, setLogs] = useState<string[]>([])
    const [status, setStatus] = useState(JobOptions[initStatus])
    const [isEventSourceFetching, setEventSourceFetching] = useState<boolean>(false)

    const {isFetching, data} = useRouterBloatLogs(uuid, !status.active && isOpen)

    useEffect(() => { if (data) setLogs(data) }, [data])
    useEffect(() => {
        setStatus(JobOptions[initStatus])
        if (!JobOptions[initStatus].active) return

        const es = BloatApi.stream.fn(uuid)
        const close = () => {
            es.close();
            setEventSourceFetching(false)
        }
        const addLog = (log: string) => setLogs((old) => [...old, log])
        es.onopen = () => {
            setLogs([])
            setEventSourceFetching(true)
            addLog("[browser] streaming open: New connection was established")
        }
        es.addEventListener(EventType.LOG, (e: MessageEvent<string>) => addLog(e.data))
        es.addEventListener(EventType.SERVER, (e: MessageEvent<string>) => addLog(e.data))
        es.addEventListener(EventType.STATUS, (e: MessageEvent<JobStatus>) => setStatus(JobOptions[e.data]))
        es.addEventListener(EventType.STREAM, (e: MessageEvent<EventStreamType>) => {
            if (e.data === EventStreamType.END) close()
        })
        es.onerror = () => {
            setEventSourceFetching(false);
            addLog("[browser] streaming error: Trying to reestablish connection")
        }

        return () => es.close()
    }, [uuid, initStatus])

    return {
        isFetching: isEventSourceFetching || isFetching,
        logs,
        status
    }
}
