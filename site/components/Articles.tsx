import styles from "./Articles.module.css"

export type Article = {
    title: string,
    description: string,
    link: string,
    image: string,
    date: string,
    readTime: number,
    tags: string[],
}

// Newest first; a new Medium post goes on top. The image is the post's own cover, served from Medium's CDN.
export const articles: Article[] = [
    {
        title: "Ivory v2.1.0 — Your Database Clusters Can Now Talk to AI",
        description: "A built-in MCP server, one deploy flow for clusters and nodes, smarter sessions and health checks you can trust a little more.",
        link: "https://anselvo.medium.com/ivory-v2-1-0-your-database-clusters-can-now-talk-to-ai-a56c39abfb79",
        image: "https://cdn-images-1.medium.com/max/1024/1*wwaJFjCuSMIyWXFq1ZfSCg.jpeg",
        date: "2026-10-15",
        readTime: 6,
        tags: ["AI", "MCP", "Release"],
    },
    {
        title: "Ivory v2.0.0 — Beyond Postgres",
        description: "From a Postgres and Patroni tool to a pluggable architecture: Keepers, etcd, Redis, ClickHouse, MongoDB and ZooKeeper, cluster deployment and user management.",
        link: "https://anselvo.medium.com/ivory-v2-0-0-beyond-postgres-4cfe865f2128",
        image: "https://cdn-images-1.medium.com/max/1024/1*x2uAOFbnGxrMesQbniV6Vg.png",
        date: "2026-09-08",
        readTime: 5,
        tags: ["Plugins", "Deployment", "Release"],
    },
    {
        title: "Ivory 1.4.0",
        description: "LDAP and OIDC/SSO sign-in, a full permissions system, cluster detection moved to the backend, backups and query improvements.",
        link: "https://anselvo.medium.com/ivory-1-4-0-aae9610f1fd5",
        image: "https://cdn-images-1.medium.com/max/1024/1*ql2QYUdx_gbflPBsSDLSRg.png",
        date: "2025-12-30",
        readTime: 2,
        tags: ["Security", "Release"],
    },
    {
        title: "Ivory v1.3.0 is Out!",
        description: "A simpler, more powerful query builder: a free-form console, prepared statement parameters and labels for your queries.",
        link: "https://anselvo.medium.com/ivory-v1-3-0-is-out-postgres-patroni-cluster-management-tool-356b4b7eb61d",
        image: "https://cdn-images-1.medium.com/max/1024/1*88N25Uix7z9f34xbMuVapg.jpeg",
        date: "2024-04-29",
        readTime: 2,
        tags: ["Queries", "Release"],
    },
    {
        title: "Troubleshoot Postgres clusters easily with Ivory",
        description: "Presenting Ivory: why it was built and how it helps you reach a Patroni cluster fast when your elephant suddenly takes a nap.",
        link: "https://anselvo.medium.com/troubleshoot-postgres-clusters-easily-with-ivory-eb0a6d61695a",
        image: "https://cdn-images-1.medium.com/max/612/0*KT8fx60SXj7hqM2Z.png",
        date: "2023-06-30",
        readTime: 3,
        tags: ["PostgreSQL", "Patroni"],
    },
]

// Formatted in UTC so the build and the browser agree on the day.
export function formatArticleDate(date: string) {
    return new Date(date).toLocaleDateString("en-US", {year: "numeric", month: "short", day: "numeric", timeZone: "UTC"})
}

export function Articles() {
    return (
        <div className={styles.list}>
            {articles.map(article => (
                <a key={article.link} className={styles.card} href={article.link} target="_blank" rel="noopener">
                    <div className={styles.cover}>
                        <img className={styles.image} src={article.image} alt="" loading="lazy"/>
                    </div>
                    <div className={styles.body}>
                        <span className={styles.meta}>medium.com · {formatArticleDate(article.date)} · {article.readTime} min read</span>
                        <h3 className={styles.title}>{article.title}</h3>
                        <p className={styles.description}>{article.description}</p>
                        <div className={styles.tags}>
                            {article.tags.map(tag => <span key={tag} className={styles.tag}>{tag}</span>)}
                        </div>
                    </div>
                </a>
            ))}
        </div>
    )
}
