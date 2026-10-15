import {articles, formatArticleDate} from "../Articles"
import styles from "./Home.module.css"
import {HomeSection} from "./HomeSection"

export function HomeArticles() {
    return (
        <HomeSection label="Articles" title="Read about Ivory" intro={renderIntro()}>
            <div className={`${styles.grid} ${styles.grid3}`}>
                {articles.slice(0, 3).map(article => (
                    <a key={article.link} className={`${styles.cell} ${styles.link}`} href={article.link} target="_blank" rel="noopener">
                        <img className={styles.cover} src={article.image} alt="" loading="lazy"/>
                        <span className={styles.meta}>{formatArticleDate(article.date)} · {article.readTime} min read</span>
                        <h3 className={styles.cellTitle}>{article.title}</h3>
                        <p className={styles.cellText}>{article.description}</p>
                    </a>
                ))}
            </div>
        </HomeSection>
    )

    function renderIntro() {
        return <>Release notes and the thinking behind them. <a className={styles.action} href="/articles">All articles →</a></>
    }
}
