import type {ReactNode} from "react"

import styles from "./Home.module.css"

type Props = {
    label: string,
    title: string,
    intro?: ReactNode,
    // "side" puts the heading in its own column beside the content, the way rust-lang.org's "Why Rust?" reads.
    layout?: "top" | "side",
    children: ReactNode,
}

export function HomeSection(props: Props) {
    const {label, title, intro, layout = "top", children} = props
    return (
        <section className={styles.section}>
            <div className={`${styles.frame} ${layout === "side" ? styles.side : ""}`}>
                <header className={styles.header}>
                    <p className={styles.label}>{label}</p>
                    <h2 className={styles.title}>{title}</h2>
                    {intro && <div className={styles.intro}>{intro}</div>}
                </header>
                <div className={styles.content}>{children}</div>
            </div>
        </section>
    )
}
