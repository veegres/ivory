import {DocContent, Layout as BasicLayout} from "@rspress/core/theme-original"

import styles from "./HomeBody.module.css"

function Layout() {
    return <BasicLayout afterFeatures={<HomeBody />} />
}

function HomeBody() {
    return (
        <div className={styles.body}>
            <div className="rp-doc rspress-doc">
                <DocContent components={undefined} />
            </div>
        </div>
    )
}

export * from "@rspress/core/theme-original"
export {Layout}
