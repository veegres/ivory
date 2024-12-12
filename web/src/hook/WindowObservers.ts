import {useCallback, useEffect, useState} from "react";

export function useWindowSize(element: Element | null) {
    const [windowWidth, setWindowWidth] = useState(element?.clientWidth)
    const [windowHeight, setWindowHeight] = useState(element?.clientHeight)

    const callback = useCallback(setWindowDimensions, [element])
    useEffect(handleEffectResize, [element, callback])

    return [windowWidth, windowHeight]

    function handleEffectResize() {
        callback()
        const observer = new ResizeObserver(callback)
        if (element) observer.observe(element)
        return () => {
            observer.disconnect()
        }
    }

    function setWindowDimensions() {
        setWindowWidth(element?.clientWidth)
        setWindowHeight(element?.clientHeight)
    }
}

export function useWindowChildCount(element: Element | null) {
    const [count, setCount] = useState(element?.childElementCount)

    const callback = useCallback(setWindowDimensions, [element])
    useEffect(handleEffectResize, [element, callback])

    return [count]

    function handleEffectResize() {
        callback()
        const observer = new MutationObserver(callback)
        if (element) observer.observe(element, {childList: true})
        return () => {
            observer.disconnect()
        }
    }

    function setWindowDimensions() {
        setCount(element?.childElementCount)
    }
}

export function useWindowScrolled(element: Element | null): [boolean] {
    const [scrolled, setScrolled] = useState(false)
    const [width, height] = useWindowSize(element)
    const [count] = useWindowChildCount(element)

    useEffect(
        () => setScrolled(!!element && element.scrollWidth > element.clientWidth),
        [element, count, width, height]
    )
    return [scrolled]
}
