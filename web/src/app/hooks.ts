import {useEffect, useState} from "react";
import {EventJob, EventStatus} from "./types";

export function useEventJob(uuid: string): EventJob {
    const [logs, setLogs] = useState<string[]>([])
    const [isFetching, setIsFetching] = useState<boolean>(false)

    useEffect(() => {
        // TODO we shouldn't run it if status of a job is finished
        const es = new EventSource(`/api/cli/pgcompacttable/${uuid}/stream`)
        const close = () => {
            es.close();
            setIsFetching(false)
        }
        const logs = (log: string) => setLogs((old) => [...old, log])
        es.onopen = () => setIsFetching(true)
        es.addEventListener("log", (e) => logs(e.data))
        es.addEventListener("server", (e) => logs(e.data))
        es.addEventListener("status", (e) => {
            if (e.data === EventStatus.FINISH) close()
        })
        es.onerror = () => {
            setIsFetching(false);
            logs("Logs streaming error: Trying to reestablish it")
        }

        return () => es.close()
    }, [uuid])


    return {
        isFetching,
        logs
    }
}
