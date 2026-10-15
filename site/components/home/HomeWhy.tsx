import styles from "./Home.module.css"
import {HomeSection} from "./HomeSection"

const reasons = [
    {
        title: "Nothing to install",
        text: "Ivory reaches your machines over plain SSH and talks to each database in its own protocol. No agent beside every node, no Kubernetes, no Ansible, no cloud account — one Docker container and you are done.",
    },
    {
        title: "One workflow, every engine",
        text: "PostgreSQL, etcd, Redis, ClickHouse, MongoDB and ZooKeeper look and behave the same: the same health view, the same switchover button, the same console. Learn it once, use it everywhere.",
    },
    {
        title: "Fix it from anywhere",
        text: "A UI that works on your phone for the 3 a.m. page, and a built-in MCP server so your AI agent can dig into a cluster with you — under exactly the permissions you have.",
    },
]

export function HomeWhy() {
    return (
        <HomeSection label="Why Ivory" title="Why Ivory?" layout="side">
            <div className={`${styles.grid} ${styles.grid3}`}>
                {reasons.map(reason => (
                    <div key={reason.title} className={styles.cell}>
                        <h3 className={styles.cellTitle}>{reason.title}</h3>
                        <p className={styles.cellText}>{reason.text}</p>
                    </div>
                ))}
            </div>
        </HomeSection>
    )
}
