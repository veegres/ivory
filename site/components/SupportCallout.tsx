import styles from "./SupportCallout.module.css"

export function SupportCallout() {
    return (
        <div className={styles.callout}>
            <div className={styles.text}>
                <p className={styles.title}>Like Ivory? Help it grow</p>
                <p className={styles.description}>
                    Ivory is free, open source and built in spare time. A star helps other people find it,
                    and sponsorship funds new database plugins and features.
                </p>
            </div>
            <div className={styles.actions}>
                <a className={styles.primary} href="https://github.com/veegres/ivory" target="_blank" rel="noopener">
                    ⭐ Star on GitHub
                </a>
                <a className={styles.secondary} href="/support">❤️ Support the project</a>
            </div>
        </div>
    )
}
