import styles from "./Home.module.css"
import {HomeSection} from "./HomeSection"

const ways = [
    {
        title: "⭐ Star on GitHub",
        text: "It takes one click, helps other people find Ivory, and keeps its author building it.",
        link: "https://github.com/veegres/ivory",
        action: "Star the repository →",
    },
    {
        title: "🧩 Add your database",
        text: "Every engine is a self-contained plugin. The repository ships an AI agent that writes one for your engine by the project's checklist.",
        link: "/docs/community/contributing",
        action: "Read the contributing guide →",
    },
    {
        title: "❤️ Support the project",
        text: "Ivory is free, open source and built in spare time. Sponsorship funds new database plugins and features.",
        link: "/support",
        action: "Become a sponsor →",
    },
]

export function HomeInvolved() {
    return (
        <HomeSection label="Community" title="Get involved" intro="Ivory is open source under GPL-3.0 and grows with the people who use it.">
            <div className={`${styles.grid} ${styles.grid3}`}>
                {ways.map(way => (
                    <a key={way.title} className={`${styles.cell} ${styles.link}`} href={way.link} {...getTarget(way.link)}>
                        <h3 className={styles.cellTitle}>{way.title}</h3>
                        <p className={styles.cellText}>{way.text}</p>
                        <span className={styles.more}>{way.action}</span>
                    </a>
                ))}
            </div>
        </HomeSection>
    )

    function getTarget(link: string) {
        return link.startsWith("http") ? {target: "_blank", rel: "noopener"} : {}
    }
}
