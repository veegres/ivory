import {useEffect, useState} from "react";

export function useWindowSize() {
    const [windowWidth, setWindowWidth] = useState(window.innerWidth)
    const [windowHeight, setWindowHeight] = useState(window.innerHeight)

    useEffect(handleUseEffectResize, [])

    return [windowWidth, windowHeight]

    function handleUseEffectResize() {
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
