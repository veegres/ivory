import path from "node:path"

import {defineConfig} from "@rspress/core"
import {pluginSitemap} from "@rspress/plugin-sitemap"

const siteUrl = "https://ivory.veegres.com"
const description = "Home for your database clusters: an open-source web UI and MCP server to deploy, watch and fix PostgreSQL & Patroni, etcd, Redis, ClickHouse, MongoDB and ZooKeeper — one Docker container, nothing to install on your hosts."

const structuredData = {
    "@context": "https://schema.org",
    "@type": "SoftwareApplication",
    name: "Ivory",
    url: siteUrl,
    image: `${siteUrl}/ivory-og.png`,
    description,
    applicationCategory: "DeveloperApplication",
    applicationSubCategory: "Database management",
    operatingSystem: "Linux, macOS, Windows (Docker)",
    license: "https://www.gnu.org/licenses/gpl-3.0.html",
    isAccessibleForFree: true,
    offers: {"@type": "Offer", price: "0", priceCurrency: "USD"},
    sameAs: ["https://github.com/veegres/ivory", "https://hub.docker.com/r/veegres/ivory"],
    keywords: "PostgreSQL GUI, Patroni UI, etcd GUI, Redis GUI, ClickHouse GUI, MongoDB GUI, ZooKeeper GUI, database cluster management, MCP server",
}

// Google takes the site name in search results from this, not from <title>; without it a subdomain shows its root domain.
const websiteData = {
    "@context": "https://schema.org",
    "@type": "WebSite",
    name: "Ivory",
    alternateName: ["Ivory by Veegres", "ivory.veegres.com"],
    url: `${siteUrl}/`,
}

export default defineConfig({
    root: "content",
    outDir: "build",
    lang: "en",
    title: "Ivory",
    description,
    icon: "/ivory-192.png",
    logo: "/ivory-fav.png",
    logoText: "Ivory",
    globalStyles: path.join(__dirname, "styles", "global.css"),
    llms: true,
    route: {
        cleanUrls: true,
    },
    head: [
        ["link", {rel: "apple-touch-icon", href: "/ivory-192.png"}],
        ["meta", {property: "og:site_name", content: "Ivory"}],
        ["meta", {property: "og:image", content: `${siteUrl}/ivory-og.png`}],
        ["meta", {property: "og:image:width", content: "1200"}],
        ["meta", {property: "og:image:height", content: "630"}],
        ["meta", {property: "og:image:alt", content: "Ivory — Deploy. Watch. Fix. Anywhere"}],
        ["meta", {name: "twitter:card", content: "summary_large_image"}],
        ["meta", {name: "twitter:image", content: `${siteUrl}/ivory-og.png`}],
        route => ["link", {rel: "canonical", href: `${siteUrl}${route.routePath}`}],
        route => ["meta", {property: "og:url", content: `${siteUrl}${route.routePath}`}],
        route => route.routePath === "/" ? `<script type="application/ld+json">${JSON.stringify(websiteData)}</script>` : undefined,
        route => route.routePath === "/" ? `<script type="application/ld+json">${JSON.stringify(structuredData)}</script>` : undefined,
    ],
    plugins: [
        pluginSitemap({siteUrl}),
    ],
    themeConfig: {
        nav: [
            {text: "Docs", link: "/docs/getting-started", activeMatch: "^/docs/"},
            {text: "MCP", link: "/docs/mcp"},
            {text: "Playground", link: "/playground"},
            {text: "Articles", link: "/articles"},
            {text: "Support ❤️", link: "/support"},
            {text: "Releases", link: "https://github.com/veegres/ivory/releases"},
        ],
        socialLinks: [
            {icon: "github", mode: "link", content: "https://github.com/veegres/ivory"},
        ],
        sidebar: {
            "/docs/": [
                {
                    text: "Get started",
                    items: [
                        {text: "What is Ivory", link: "/docs/about"},
                        {text: "Quick start", link: "/docs/getting-started"},
                        {text: "Core concepts", link: "/docs/concepts"},
                        {text: "FAQ", link: "/docs/faq"},
                    ],
                },
                {
                    text: "Using Ivory",
                    items: [
                        {text: "Add a cluster", link: "/docs/clusters"},
                        {text: "Deploy a cluster", link: "/docs/deployment"},
                        {text: "Cluster overview", link: "/docs/overview"},
                        {text: "Node view", link: "/docs/node"},
                        {
                            text: "Tools",
                            collapsible: true,
                            collapsed: true,
                            items: [
                                {text: "pgcompacttable (PostgreSQL)", link: "/docs/tools/postgres/pg_compacttable"},
                            ],
                        },
                    ],
                },
                {
                    text: "AI agents",
                    items: [
                        {text: "MCP server", link: "/docs/mcp"},
                    ],
                },
                {
                    text: "Supported databases",
                    items: [
                        {text: "Comparison", link: "/docs/databases"},
                        {text: "PostgreSQL + Patroni", link: "/docs/databases/patroni"},
                        {text: "PostgreSQL", link: "/docs/databases/postgres"},
                        {text: "etcd", link: "/docs/databases/etcd"},
                        {text: "Redis", link: "/docs/databases/redis"},
                        {text: "ClickHouse", link: "/docs/databases/clickhouse"},
                        {text: "MongoDB", link: "/docs/databases/mongodb"},
                        {text: "ZooKeeper", link: "/docs/databases/zookeeper"},
                    ],
                },
                {
                    text: "Self-hosting",
                    items: [
                        {text: "Installation & configuration", link: "/docs/configuration"},
                        {text: "Authentication & users", link: "/docs/authentication"},
                    ],
                },
                {
                    text: "Community",
                    items: [
                        {text: "Contributing", link: "/docs/community/contributing"},
                        {text: "Security policy", link: "/docs/community/security"},
                        {text: "Code of conduct", link: "/docs/community/code-of-conduct"},
                    ],
                },
            ],
        },
        footer: {
            message: `<div><a href="/docs/community/contributing">Contributing</a> · <a href="/docs/community/security">Security</a> · <a href="/docs/community/code-of-conduct">Code of conduct</a> · <a href="/support">Support Ivory</a></div><div>GPL-3.0 License © ${new Date().getFullYear()} Veegres Labs</div>`,
        },
    },
})
