import {DocContent, HomeLayout as BasicHomeLayout, Layout as BasicLayout} from "@rspress/core/theme-original"
import type {ComponentProps} from "react"

import {HomeRelease} from "../components/home/HomeRelease"
import styles from "./HomeBody.module.css"

function Layout() {
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

export * from "@rspress/core/theme-original"
export {Layout}
