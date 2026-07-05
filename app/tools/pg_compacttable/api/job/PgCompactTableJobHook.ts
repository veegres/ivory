import {useEffect, useState} from "react"

import {JobOptions} from "../../../../shared/helper/HelperUtils"
import {useRouterPgCompactTableLogs} from "../PgCompactTableHook"
import {PgCompactTableApi} from "../PgCompactTableRouter"
import {EventStreamType, EventType, JobStatus} from "./PgCompactTableJobType"

type Hook = {
    isFetching: boolean;
    logs: string[];
    status: { name: string; color: string; active: boolean }
}

export function useRouterPgCompactTableJob(uuid: string, initStatus: JobStatus, isOpen: boolean, refetchList: () => void): Hook {
    const [status, setStatus] = useState(JobOptions[initStatus])
    const [isEventSourceFetching, setEventSourceFetching] = useState<boolean>(false)

    const fetchFile = !status.active && isOpen
    const [logs, setLogs] = useState<string[]>([])
    const {data, error, isFetching, isError} = useRouterPgCompactTableLogs(uuid, fetchFile)

    useEffect(handleEffectStream, [uuid, initStatus, refetchList])

    return {
        isFetching: isEventSourceFetching || isFetching,
        logs: !fetchFile ? logs : (data ?? (!isError ? [] : [`[browser] streaming error: ${error.message ?? "unknown"}`])),
        status: status,
    }

    function handleEffectStream() {
        if (!JobOptions[initStatus].active) return

        const es = PgCompactTableApi.stream.fn(uuid)
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
    }
}
