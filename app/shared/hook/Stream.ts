import {useCallback, useEffect, useState} from "react"

import {EventStreamType, EventType} from "../../tools/pg_compacttable/api/job/PgCompactTableJobType"

type Options = {
    enabled?: boolean,
}

export function useStream(url: string, options?: Options) {
    const {enabled = true} = options ?? {}
    const [loading, setLoading] = useState(false)
    const [response, setResponse] = useState<string[]>([])
    const [nonce, setNonce] = useState(0)

    const push = useCallback((type: EventType, message: string) =>
        setResponse((prev) => [...prev, `[${type}] ${message}`]), [])

    useEffect(() => {
        if (!url || !enabled) return
        const es = new EventSource(url)
        es.onopen = () => {
            setLoading(true)
            push(EventType.BROWSER, "new connection was established")
        }
        Object.values(EventType).forEach((type) => {
            es.addEventListener(type, (e) => push(type, e.data))
        })
        es.addEventListener(EventType.STREAM, (e) => {
            if (e.data === EventStreamType.END) {
                setLoading(false)
                es.close()
            }
        })
        es.onerror = () => {
            setLoading(false)
            push(EventType.BROWSER, "trying to reestablish connection")
        }
        return () => {
            es.close()
            setLoading(false)
            setResponse([])
        }
        // NOTE: nonce is not read, it only exists to let `reconnect` force this effect
        //  to tear down the current EventSource and open a fresh one from scratch.
    }, [push, url, enabled, nonce])

    return {loading, response, reconnect: () => setNonce(n => n + 1)}
}