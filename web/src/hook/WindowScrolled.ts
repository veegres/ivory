import {useWindowSize} from "./WindowSize";
import {useEffect, useState} from "react";

export function useWindowScrolled<T>(element: Element | null, elements: T[] = []): [boolean, T[]] {
    const [scrolled, setScrolled] = useState(false)
    const [width, height] = useWindowSize()
    useEffect(
        () => setScrolled(!!element && element.scrollWidth > element.clientWidth),
        [element, elements, width, height]
    )
    return [scrolled, elements]
}
