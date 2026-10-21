import styles from "./Articles.module.css"
import feed from "./ArticlesFeed.json"

export type Article = {
    title: string,
    description: string,
    link: string,
    image: string,
    date: string,
    readTime?: number,
    tags: string[],
}

// Kept in sync with the Medium tag feed by scripts/sync-articles.mjs, which adds new posts on every build.
export const articles: Article[] = feed

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
                        <span className={styles.meta}>medium.com · {formatArticleDate(article.date)}{article.readTime ? ` · ${article.readTime} min read` : ""}</span>
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
