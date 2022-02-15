import {useEffect, useState} from "react";
import {EventJob, EventStatus} from "./types";

export function useEventJob(uuid: string): EventJob {
    const [logs, setLogs] = useState<string[]>([])
    const [isFetching, setIsFetching] = useState<boolean>(true)

    useEffect(() => {
        const es = new EventSource(`/api/cli/pgcompacttable/${uuid}/stream`)
        es.addEventListener("log", (e) => setLogs((old) => [...old, e.data]))
        es.addEventListener("status", (e) => {
            if (e.data === EventStatus.FINISH) {
                es.close()
                setIsFetching(false)
            }
        })

        return () => es.close()
    }, [uuid])

    return {
        isFetching,
        logs
    }
}
