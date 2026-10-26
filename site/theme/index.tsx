import {useDark} from "@rspress/core/runtime"
import {DocContent, HomeLayout as BasicHomeLayout, Layout as BasicLayout} from "@rspress/core/theme-original"
import {type ComponentProps, useEffect} from "react"

import {HomeRelease} from "../components/home/HomeRelease"
import styles from "./HomeBody.module.css"

// Same colours as --rp-c-bg in styles/global.css and the manifests in content/public
const themeColor = {light: "#fffdf6", dark: "#15130f"}

function Layout() {
    const dark = useDark()
    useEffect(() => handleThemeColor(dark), [dark])
    return <BasicLayout HomeLayout={HomeLayout} afterFeatures={<HomeBody />} />
}

// Layout accepts afterHeroActions but never passes it on to the hero, so it is handed to HomeLayout here.
function HomeLayout(props: ComponentProps<typeof BasicHomeLayout>) {
    return <BasicHomeLayout {...props} afterHeroActions={<HomeRelease />} />
}

function HomeBody() {
    return (
        <div className={styles.body}>
            <DocContent components={undefined} />
        </div>
    )
}

// NOTE: the head follows the system before the page draws, a theme picked by the toggle must win over it, so one tag
//  without a media query is left; an installed site takes the manifest's colour instead, so the manifest follows too
function handleThemeColor(dark: boolean) {
    document.head.querySelector<HTMLLinkElement>("link[rel='manifest']")
        ?.setAttribute("href", dark ? "/manifest-dark.json" : "/manifest-light.json")
    const tags = document.head.querySelectorAll<HTMLMetaElement>("meta[name='theme-color']")
    tags.forEach((tag, index) => index === 0 ? tag.removeAttribute("media") : tag.remove())
    const tag = tags[0] ?? document.head.appendChild(Object.assign(document.createElement("meta"), {name: "theme-color"}))
    tag.content = dark ? themeColor.dark : themeColor.light
}

export * from "@rspress/core/theme-original"
export {Layout}
