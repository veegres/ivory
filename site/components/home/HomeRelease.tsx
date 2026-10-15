import {useEffect, useState} from "react"

import styles from "./Home.module.css"

const latest = "https://github.com/veegres/ivory/releases/latest"

type Release = {tag_name: string, html_url: string}

// Asked of GitHub in the browser rather than at build time, so a release published after the last deploy still shows.
export function HomeRelease() {
    const [release, setRelease] = useState<Release>()

    useEffect(handleEffectFetch, [])

    return (
        <a className={styles.release} href={release?.html_url ?? latest} target="_blank" rel="noopener">
            {release ? <>Version <b>{release.tag_name.replace(/^v/, "")}</b></> : "Latest release"}
        </a>
    )

    function handleEffectFetch() {
        const controller = new AbortController()
        fetch("https://api.github.com/repos/veegres/ivory/releases/latest", {signal: controller.signal})
            .then(response => response.ok ? response.json() : undefined)
            .then(data => data?.tag_name && data?.html_url && setRelease(data))
            .catch(() => undefined)
        return () => controller.abort()
    }
}
