import {Dispatch, SetStateAction, useEffect, useState} from "react";

export function useLocalStorageState<T>(key: string, initialState: T | (() => T)): [T, Dispatch<SetStateAction<T>>] {
    const localStorageValue = localStorage.getItem(key)
    const initValue = localStorageValue !== null ? JSON.parse(localStorageValue) as T : initialState
    const [value, setValue] = useState(initValue)
    useEffect(() => localStorage.setItem(key, JSON.stringify(value)), [value, key])
    return [value, setValue]
}
