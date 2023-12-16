import {useEffect, useState} from "react";

export function useDebounce<T>(value: T, delay: number = 500): T {
    const [debouncedValue, setDebouncedValue] = useState<T>(value)

    useEffect(() => {
        const timer = setTimeout(() => setDebouncedValue(value), delay)

        return () => clearTimeout(timer)
    }, [value, delay])

    return debouncedValue
}

export function useDebounceFunction<T>(value: T, onChange: (v: T) => void, delay: number = 500) {
    useEffect(() => {
        const timeout = setTimeout(() => onChange(value), delay)
        return () => clearTimeout(timeout)
    }, [value, onChange, delay]);
}
