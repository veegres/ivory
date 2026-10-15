import mcp from "../../content/docs/images/mcp_example.png"
import styles from "./Home.module.css"
import {HomeSection} from "./HomeSection"

const prompts = [
    "Why is replica 2 lagging?",
    "What is blocking my queries?",
    "Switch the leader to node 3.",
]

export function HomeAgent() {
    return (
        <HomeSection label="AI-ready" title="Just ask your cluster" layout="side" intro={renderIntro()}>
            <div className={styles.agent}>
                <ul className={styles.prompts}>
                    {prompts.map(prompt => <li key={prompt} className={styles.prompt}>{prompt}</li>)}
                </ul>
                <img
                    className={styles.shot}
                    src={mcp}
                    alt="An AI agent answering questions about database clusters through Ivory's MCP server"
                    loading="lazy"
                />
            </div>
        </HomeSection>
    )

    function renderIntro() {
        return (
            <>
                <p>
                    Ivory has a built-in MCP server. Connect Claude, Cursor, Codex or any MCP client and talk to your
                    clusters in plain words — the agent reads their health, runs diagnostics and acts only where your
                    own permissions allow.
                </p>
                <a className={styles.action} href="/docs/mcp">Connect an agent →</a>
            </>
        )
    }
}
