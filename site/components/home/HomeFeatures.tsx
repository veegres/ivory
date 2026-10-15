import type {ReactNode} from "react"

import styles from "./Home.module.css"
import {HomeSection} from "./HomeSection"

const features: {title: string, text: string, link: string, icon: ReactNode}[] = [
    {
        title: "Deploy over plain SSH",
        text: "Pick a template, point it at your machines and press Deploy. A high-availability cluster comes up — no Kubernetes, Ansible or agent.",
        link: "/docs/deployment",
        icon: <><path d="M12 15V3" /><path d="m7 8 5-5 5 5" /><rect x="3" y="15" width="18" height="6" rx="1.5" /><path d="M7 18h.01" /></>,
    },
    {
        title: "Health at a glance",
        text: "Roles, replication lag, pending restarts and warnings for every node — including the one that quietly left the cluster.",
        link: "/docs/overview",
        icon: <path d="M3 12h4l3-8 4 16 3-8h4" />,
    },
    {
        title: "HA operations in a click",
        text: "Switchover, failover, reinitialise, restart, reload and pause — scheduled for later where the cluster supports it.",
        link: "/docs/overview#ha-operations",
        icon: <><path d="M4 8h13l-3-3" /><path d="M20 16H7l3 3" /></>,
    },
    {
        title: "The machine behind each node",
        text: "Containers, CPU, memory, network, processes and any log file on the host — read and controlled over SSH.",
        link: "/docs/node",
        icon: <><rect x="3" y="4" width="18" height="12" rx="1.5" /><path d="M8 20h8" /><path d="M12 16v4" /><path d="m7 9 2 2-2 2" /></>,
    },
    {
        title: "Query console and tools",
        text: "Saved diagnostic queries, a console for every engine, and maintenance tools like pgcompacttable with live output.",
        link: "/docs/node#database",
        icon: <><ellipse cx="12" cy="6" rx="8" ry="3" /><path d="M4 6v12c0 1.7 3.6 3 8 3s8-1.3 8-3V6" /><path d="M4 12c0 1.7 3.6 3 8 3s8-1.3 8-3" /></>,
    },
    {
        title: "Teams and permissions",
        text: "Basic, LDAP or OIDC sign-in, and a permission for every action — enforced the same in the UI and for AI agents.",
        link: "/docs/authentication",
        icon: <><path d="M12 3 4 6v6c0 4.5 3.4 8.2 8 9 4.6-.8 8-4.5 8-9V6l-8-3Z" /><path d="m9 12 2 2 4-4" /></>,
    },
]

export function HomeFeatures() {
    return (
        <HomeSection
            label="Features"
            title="Everything a cluster needs, in one place"
            intro="From the first deploy to the late-night fix, without switching tools or SSH-ing into hosts by hand."
        >
            <div className={`${styles.grid} ${styles.grid3}`}>
                {features.map(feature => (
                    <a key={feature.title} className={`${styles.cell} ${styles.link}`} href={feature.link}>
                        <svg className={styles.icon} viewBox="0 0 24 24" aria-hidden="true">{feature.icon}</svg>
                        <h3 className={styles.cellTitle}>{feature.title}</h3>
                        <p className={styles.cellText}>{feature.text}</p>
                        <span className={styles.more}>Learn more →</span>
                    </a>
                ))}
            </div>
        </HomeSection>
    )
}
