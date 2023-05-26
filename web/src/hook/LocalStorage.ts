import {Dispatch, SetStateAction, useEffect, useState} from "react";

export function useLocalStorageState<T>(key: string, initialState: T | (() => T), sync = false): [T, Dispatch<SetStateAction<T>>] {
    const isObject = typeof initialState === "object"
    const localStorageValue = localStorage.getItem(key)
    const [value, setValue] = useState(getInitValue())

    useEffect(handleEffectStorageChange, [initialState, isObject, sync, key])
    useEffect(handleEffectItemUpdate, [key, value, isObject])

    return [value, setValue]

    function getInitValue() {
        if (localStorageValue === null) return initialState
        else {
            if (isObject) return JSON.parse(localStorageValue) as T
            else return localStorageValue as T
        }
    }

    function handleEffectItemUpdate() {
        if (value !== undefined && value !== null) {
            if (isObject) localStorage.setItem(key, JSON.stringify(value))
            else localStorage.setItem(key, value.toString())
        }
    }

    function handleEffectStorageChange() {
        if (sync) {
            window.addEventListener("storage", (e) => {
                if (e.key === key) {
                    if (e.newValue === null) setValue(initialState)
                    else {
                        if (isObject) setValue(JSON.parse(e.newValue) as T)
                        else setValue(e.newValue as T)
                    }
                }
            })
        }
    }
}
