import {useEffect, useState} from "react"

import {useRouterBloatLogs} from "../api/bloat/hook"
import {EventStreamType, EventType, JobStatus} from "../api/bloat/job/type"
import {BloatApi} from "../api/bloat/router"
import {JobOptions} from "../app/utils"

type EventJob = {
    isFetching: boolean;
    logs: string[];
    status: { name: string; color: string; active: boolean }
}

export function useEventJob(uuid: string, initStatus: JobStatus, isOpen: boolean, refetchList: () => void): EventJob {
    const [status, setStatus] = useState(JobOptions[initStatus])
    const [isEventSourceFetching, setEventSourceFetching] = useState<boolean>(false)

    const fetchFile = !status.active && isOpen
    const [logs, setLogs] = useState<string[]>([])
    const {data, error, isFetching, isError} = useRouterBloatLogs(uuid, fetchFile)

    useEffect(() => {
        if (!JobOptions[initStatus].active) return

        const es = BloatApi.stream.fn(uuid)
        const close = () => {
            es.close()
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
            if (e.data === EventStreamType.END) {
                // NOTE: we need to refresh the job list because it keeps old job status
                refetchList()
                close()
            }
        })
        es.onerror = () => {
            setEventSourceFetching(false)
            addLog("[browser] streaming error: Trying to reestablish connection")
        }

        return () => es.close()
    }, [uuid, initStatus, refetchList])

    return {
        isFetching: isEventSourceFetching || isFetching,
        logs: !fetchFile ? logs : (data ?? (!isError ? [] : [`[browser] streaming error: ${error.message ?? "unknown"}`])),
        status: status,
    }
}
