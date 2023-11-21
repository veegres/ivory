import {useEffect, useState} from "react";

export function useWindowSize() {
    const [windowWidth, setWindowWidth] = useState(window.innerWidth)
    const [windowHeight, setWindowHeight] = useState(window.innerHeight)

    useEffect(handleEffectResize, [])

    return [windowWidth, windowHeight]

    function handleEffectResize() {
        window.addEventListener('resize', setWindowDimensions);
        return () => {
            window.removeEventListener('resize', setWindowDimensions)
        }
    }

    function setWindowDimensions() {
        setWindowWidth(window.innerWidth)
        setWindowHeight(window.innerHeight)
    }
}
