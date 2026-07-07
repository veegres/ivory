import {useCallback, useEffect, useState} from "react"

import {EventStreamType, EventType} from "../../tools/pg_compacttable/api/job/PgCompactTableJobType"

type Options = {
    enabled?: boolean,
    maxConsecutiveErrors?: number,
}

export function useStream(url: string, options?: Options) {
    const {enabled = true, maxConsecutiveErrors = 3} = options ?? {}
    const [loading, setLoading] = useState(false)
    const [response, setResponse] = useState<string[]>([])
    const [nonce, setNonce] = useState(0)

    const push = useCallback((type: EventType, message: string) =>
        setResponse((prev) => [...prev, `[${type}] ${message}`]), [])

    useEffect(() => {
        if (!url || !enabled) return
        const es = new EventSource(url)
        let consecutiveErrors = 0
        es.onopen = () => {
            consecutiveErrors = 0
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
            consecutiveErrors += 1
            if (consecutiveErrors >= maxConsecutiveErrors) {
                push(EventType.BROWSER, "connection failed repeatedly, giving up")
                es.close()
                return
            }
            push(EventType.BROWSER, "trying to reestablish connection")
        }
        return () => {
            es.close()
            setLoading(false)
            setResponse([])
        }
        // NOTE: nonce is not read, it only exists to let `reconnect` force this effect
        //  to tear down the current EventSource and open a fresh one from scratch.
    }, [push, url, enabled, nonce, maxConsecutiveErrors])

    return {loading, response, reconnect: () => setNonce(n => n + 1)}
}