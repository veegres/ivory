import {MutableRefObject, useEffect, useState} from "react";

export function useClickDetection(ref: MutableRefObject<any>) {
    const [click, setClick] = useState<"inside" | "outside">()

    useEffect(handleEffect, [ref])

    return click

    function handleEffect() {
        function handleClickOutside(event: MouseEvent) {
            if (ref.current && ref.current.contains(event.target)) {
                setClick("inside")
            } else {
                setClick("outside")
            }
        }

        document.addEventListener("mousedown", handleClickOutside)
        return () => {
            document.removeEventListener("mousedown", handleClickOutside)
        }
    }
}
