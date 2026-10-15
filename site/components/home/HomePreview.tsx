import demo from "../../content/docs/images/demo.png"
import styles from "./Home.module.css"

export function HomePreview() {
    return (
        <section className={styles.section}>
            <div className={`${styles.frame} ${styles.preview}`}>
                <img
                    className={styles.shot}
                    src={demo}
                    alt="Ivory web UI showing a PostgreSQL Patroni cluster: node roles, replication lag and host metrics"
                    loading="lazy"
                />
            </div>
        </section>
    )
}
