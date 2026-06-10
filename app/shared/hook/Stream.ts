import {useCallback, useEffect, useState} from "react"

import {EventStreamType, EventType} from "../../features/bloat/job/type"

export function useStream(url: string) {
    const [loading, setLoading] = useState(false)
    const [response, setResponse] = useState<string[]>([])

    const push = useCallback((type: EventType, message: string) =>
        setResponse((prev) => [...prev, `[${type}] ${message}`]), [])

    useEffect(() => {
        if (!url) return
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
    }, [push, url])

    return {loading, response}
}