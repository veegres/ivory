// Adds posts tagged ivory-veegres from the Medium feed to components/ArticlesFeed.json; entries already there keep their hand-written fields, and a failed fetch leaves the file as it is.
import {readFileSync, writeFileSync} from "node:fs"
import {dirname, join} from "node:path"
import {fileURLToPath} from "node:url"

const file = join(dirname(fileURLToPath(import.meta.url)), "..", "components", "ArticlesFeed.json")
const feed = "https://medium.com/feed/@anselvo"
const marker = "ivory-veegres"
const wordsPerMinute = 265

const entities = {amp: "&", lt: "<", gt: ">", quot: "\"", apos: "'", nbsp: " "}

function decode(text) {
    return text.replace(/&(#x[0-9a-f]+|#\d+|[a-z]+);/gi, (match, name) => {
        if (name[0] === "#") {
            const code = name[1].toLowerCase() === "x" ? parseInt(name.slice(2), 16) : parseInt(name.slice(1), 10)
            return String.fromCodePoint(code)
        }
        return entities[name.toLowerCase()] ?? match
    })
}

function strip(html) {
    return decode(html.replace(/<[^>]+>/g, " ")).replace(/\s+/g, " ").trim()
}

function tag(xml, name) {
    const match = xml.match(new RegExp(`<${name}>(?:<!\\[CDATA\\[)?([\\s\\S]*?)(?:\\]\\]>)?</${name}>`))
    return match?.[1].trim()
}

function parseItem(xml) {
    const categories = [...xml.matchAll(/<category><!\[CDATA\[([^\]]+)\]\]><\/category>/g)].map(m => m[1])
    const title = tag(xml, "title")
    const link = tag(xml, "link")?.split("?")[0]
    const published = tag(xml, "pubDate")
    const content = tag(xml, "content:encoded") ?? ""
    if (!categories.includes(marker) || !title || !link || !published) return undefined
    const subtitle = content.match(/<h3[^>]*>([\s\S]*?)<\/h3>/)?.[1]
    const paragraph = content.match(/<p[^>]*>([\s\S]*?)<\/p>/)?.[1]
    return {
        title: decode(title),
        description: strip(subtitle ?? paragraph ?? ""),
        link,
        image: content.match(/<img[^>]+src="([^"]+)"/)?.[1] ?? "",
        date: new Date(published).toISOString().slice(0, 10),
        readTime: Math.max(1, Math.ceil(strip(content).split(" ").length / wordsPerMinute)),
        tags: categories.filter(c => c !== marker).slice(0, 3),
    }
}

async function main() {
    const response = await fetch(feed, {signal: AbortSignal.timeout(15000)})
    if (!response.ok) throw new Error(`HTTP ${response.status}`)
    const items = [...(await response.text()).matchAll(/<item>([\s\S]*?)<\/item>/g)].map(m => parseItem(m[1])).filter(Boolean)

    const current = JSON.parse(readFileSync(file, "utf8"))
    const known = new Set(current.map(a => a.link))
    const added = items.filter(item => !known.has(item.link))
    if (added.length === 0) return console.log("articles: up to date")

    const merged = [...added, ...current].sort((a, b) => b.date.localeCompare(a.date))
    writeFileSync(file, JSON.stringify(merged, undefined, 4) + "\n")
    console.log(`articles: added ${added.map(a => a.title).join(", ")}`)
}

main().catch(error => console.warn(`articles: feed not synced (${error.message}), keeping the current list`))
