import {useState} from "react"

import styles from "./Home.module.css"
import {HomeSection} from "./HomeSection"

const command = "docker run -p 80:80 --restart always veegres/ivory"

const steps = [
    {title: "Run the container", text: "One image, no external dependencies — from Docker Hub or the GitHub Container registry."},
    {title: "Open the setup wizard", text: "Go to http://localhost, choose how people sign in and set the secret key."},
    {title: "Add a cluster", text: "Type its node addresses, let Ivory detect the rest from one, or deploy a new one."},
]

export function HomeStart() {
    const [copied, setCopied] = useState(false)

    return (
        <HomeSection
            label="Get started"
            title="Up and running in a minute"
            layout="side"
            intro={<a className={styles.action} href="/docs/getting-started">Read the guide →</a>}
        >
            <div className={styles.command}>
                <code>{command}</code>
                <button className={styles.copy} type="button" onClick={handleCopy}>{copied ? "Copied" : "Copy"}</button>
            </div>
            <ol className={`${styles.grid} ${styles.grid3} ${styles.steps}`}>
                {steps.map((step, index) => (
                    <li key={step.title} className={styles.cell}>
                        <span className={styles.step}>{String(index + 1).padStart(2, "0")}</span>
                        <h3 className={styles.cellTitle}>{step.title}</h3>
                        <p className={styles.cellText}>{step.text}</p>
                    </li>
                ))}
            </ol>
        </HomeSection>
    )

    function handleCopy() {
        navigator.clipboard.writeText(command).then(() => {
            setCopied(true)
            setTimeout(() => setCopied(false), 1500)
        })
    }
}
