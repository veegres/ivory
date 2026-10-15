import styles from "./Home.module.css"
import {HomeSection} from "./HomeSection"

const databases = [
    {name: "PostgreSQL + Patroni", cluster: "High-availability cluster with Patroni", link: "/docs/databases/patroni", stage: "Stable"},
    {name: "PostgreSQL", cluster: "Primary and standby replication", link: "/docs/databases/postgres", stage: "Stable"},
    {name: "etcd", cluster: "Raft cluster", link: "/docs/databases/etcd", stage: "Beta"},
    {name: "Redis", cluster: "Master and replica replication", link: "/docs/databases/redis", stage: "Beta"},
    {name: "ClickHouse", cluster: "Replicated cluster", link: "/docs/databases/clickhouse", stage: "Beta"},
    {name: "MongoDB", cluster: "Replica set", link: "/docs/databases/mongodb", stage: "Beta"},
    {name: "ZooKeeper", cluster: "Ensemble", link: "/docs/databases/zookeeper", stage: "Beta"},
]

export function HomeDatabases() {
    return (
        <HomeSection
            label="Supported"
            title="Supported databases and clusters"
            intro="Deploy a new cluster from a ready-made template, or add one you already run — every engine gets the same health view, operations and console."
        >
            <div className={`${styles.grid} ${styles.grid4}`}>
                {databases.map(database => (
                    <a key={database.link} className={`${styles.cell} ${styles.link}`} href={database.link}>
                        <span className={styles.stage}>{database.stage}</span>
                        <h3 className={styles.dbName}>{database.name}</h3>
                        <p className={styles.cellText}>{database.cluster}</p>
                    </a>
                ))}
                <a className={`${styles.cell} ${styles.link}`} href="/docs/databases#adding-your-database">
                    <span className={styles.stage}>Plugin</span>
                    <h3 className={styles.dbName}>Yours next?</h3>
                    <p className={styles.cellText}>Every engine is a plugin — add the one you run.</p>
                </a>
            </div>
        </HomeSection>
    )
}
