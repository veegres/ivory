import path from "node:path"

import {defineConfig} from "@rspress/core"
import {pluginSitemap} from "@rspress/plugin-sitemap"

const siteUrl = "https://ivory.veegres.com"
const description = "Open-source web UI and AI-ready MCP server for high-availability database clusters — PostgreSQL & Patroni, etcd, Redis, ClickHouse, MongoDB and ZooKeeper."

const structuredData = {
    "@context": "https://schema.org",
    "@type": "SoftwareApplication",
    name: "Ivory",
    url: siteUrl,
    image: `${siteUrl}/og-image.png`,
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

export default defineConfig({
    root: "content",
    outDir: "build",
    lang: "en",
    title: "Ivory",
    description,
    icon: "/ivory-fav.png",
    logo: "/ivory-fav.png",
    logoText: "Ivory",
    globalStyles: path.join(__dirname, "styles", "global.css"),
    llms: true,
    route: {
        cleanUrls: true,
    },
    head: [
        ["meta", {property: "og:site_name", content: "Ivory"}],
        ["meta", {property: "og:image", content: `${siteUrl}/og-image.png`}],
        ["meta", {name: "twitter:card", content: "summary_large_image"}],
        ["meta", {name: "twitter:image", content: `${siteUrl}/og-image.png`}],
        route => ["link", {rel: "canonical", href: `${siteUrl}${route.routePath}`}],
        route => ["meta", {property: "og:url", content: `${siteUrl}${route.routePath}`}],
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
            {text: "Support ❤️", link: "/support"},
            {text: "Releases", link: "https://github.com/veegres/ivory/releases"},
        ],
        socialLinks: [
            {icon: "github", mode: "link", content: "https://github.com/veegres/ivory"},
        ],
        sidebar: {
            "/docs/": [
                {
                    text: "Introduction",
                    items: [
                        {text: "About", link: "/docs/about"},
                        {text: "Getting started", link: "/docs/getting-started"},
                        {text: "FAQ", link: "/docs/faq"},
                    ],
                },
                {
                    text: "Clusters",
                    items: [
                        {text: "Cluster list", link: "/docs/clusters"},
                        {text: "Overview & HA operations", link: "/docs/overview"},
                        {text: "Deployment", link: "/docs/deployment"},
                    ],
                },
                {
                    text: "Nodes",
                    items: [
                        {text: "Node", link: "/docs/node"},
                    ],
                },
                {
                    text: "Tools",
                    items: [
                        {
                            text: "Postgres",
                            collapsible: false,
                            items: [
                                {text: "pgcompacttable", link: "/docs/tools/postgres/pg_compacttable"},
                            ],
                        },
                    ],
                },
                {
                    text: "Operating Ivory",
                    items: [
                        {text: "Authentication & users", link: "/docs/authentication"},
                        {text: "Configuration", link: "/docs/configuration"},
                        {text: "MCP server", link: "/docs/mcp"},
                    ],
                },
            ],
        },
        footer: {
            message: "Released under the GPL-3.0 License. · <a href=\"/support\">Support Ivory</a>",
        },
    },
})
