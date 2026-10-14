import tier1 from "../content/docs/images/support_tier_1.png"
import tier2 from "../content/docs/images/support_tier_2.png"
import tier3 from "../content/docs/images/support_tier_3.png"
import styles from "./SupportTiers.module.css"

const tiers = [
    {
        title: "General Development Support",
        description: "Keeps Ivory moving: maintenance, bug fixes, releases, documentation and answering issues.",
        image: tier1,
        link: "https://boosty.to/anselvo/purchase/1454406?ssource=DIRECT&share=subscription_link",
    },
    {
        title: "New Plugin Integration Support",
        description: "Funds support for more databases and keepers — every engine is a plugin with its own HA view and query console.",
        image: tier2,
        link: "https://boosty.to/anselvo/purchase/4049256?ssource=DIRECT&share=subscription_link",
    },
    {
        title: "New Features Support",
        description: "Funds new capabilities across every engine: deployment, monitoring, HA operations, maintenance tools and MCP.",
        image: tier3,
        link: "https://boosty.to/anselvo/purchase/4049261?ssource=DIRECT&share=subscription_link",
    },
]

export function SupportTiers() {
    return (
        <div className={styles.grid}>
            {tiers.map(tier => (
                <a key={tier.title} className={styles.card} href={tier.link} target="_blank" rel="noopener sponsored">
                    <img className={styles.image} src={tier.image} alt={`Ivory ${tier.title}`} loading="lazy" width={1179} height={643}/>
                    <div className={styles.body}>
                        <h3 className={styles.title}>{tier.title}</h3>
                        <p className={styles.description}>{tier.description}</p>
                        <span className={styles.button}>Support on Boosty</span>
                    </div>
                </a>
            ))}
        </div>
    )
}
