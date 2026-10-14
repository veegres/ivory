// Copies the repository's .doc/ into the site, so the Markdown stays readable on GitHub as the single source.
import {cpSync, mkdirSync, readdirSync, readFileSync, rmSync, writeFileSync} from "node:fs"
import {dirname, join, posix} from "node:path"
import {fileURLToPath} from "node:url"

const site = join(dirname(fileURLToPath(import.meta.url)), "..")
const repo = join(site, "..")
const source = join(repo, ".doc")
const target = join(site, "content", "docs")
const publicDir = join(site, "content", "public")
const github = "https://github.com/veegres/ivory/blob/master"

// What a search engine shows for each page; a doc missing here fails the build so none ships without one.
const meta = {
    "about.md": {
        title: "About Ivory — open-source database cluster management UI",
        description: "What Ivory is: a self-hosted web UI and MCP server to deploy, monitor and operate high-availability PostgreSQL/Patroni, etcd, Redis, ClickHouse, MongoDB and ZooKeeper clusters.",
    },
    "getting-started.md": {
        title: "Getting started — install Ivory with Docker",
        description: "Run Ivory with one docker run command, complete the setup wizard and add your first PostgreSQL, Patroni, etcd, Redis, ClickHouse, MongoDB or ZooKeeper cluster.",
    },
    "faq.md": {
        title: "FAQ",
        description: "Answers about upgrading Ivory, where it stores data, agents, credentials encryption, reverse proxies, HTTPS, AI agents and using it from a phone.",
    },
    "clusters.md": {
        title: "Clusters — add, auto-detect and organize database clusters",
        description: "Register an existing database cluster manually, auto-detect it from a single node address, or deploy a new one — then tag and filter your clusters in Ivory.",
    },
    "overview.md": {
        title: "Cluster overview — health, replication lag, switchover and failover",
        description: "Ivory's cluster dashboard: node roles, replication lag, pending restarts and warnings, plus HA operations — switchover, failover, reinitialise, restart, reload and pause.",
    },
    "node.md": {
        title: "Node — VM metrics, containers, keeper, query console and tools",
        description: "Per-node operations in Ivory: host CPU, memory and logs over SSH, Docker container lifecycle, keeper configuration, SQL query builder and console, and maintenance tools.",
    },
    "deployment.md": {
        title: "Deployment — deploy HA database clusters over SSH with Docker",
        description: "Deploy a high-availability PostgreSQL/Patroni, etcd, Redis, ClickHouse, MongoDB or ZooKeeper cluster to your own VMs over SSH from reusable templates — no Kubernetes or Ansible.",
    },
    "authentication.md": {
        title: "Authentication & users — Basic, LDAP, OIDC/SSO and permissions",
        description: "Secure Ivory with Basic, LDAP or OIDC/SSO sign-in, register users with one-time links, manage superusers and grant granular per-user permissions.",
    },
    "configuration.md": {
        title: "Configuration — persistence, upgrades, TLS and environment variables",
        description: "Run Ivory in production: Docker data persistence, upgrading with backup and restore, environment variables, HTTPS/TLS, reverse-proxy sub-path and secret encryption.",
    },
    "mcp.md": {
        title: "MCP server — connect Claude, Cursor or Codex to your database clusters",
        description: "Ivory's built-in Model Context Protocol server lets AI agents inspect and operate PostgreSQL, Patroni, etcd, Redis, ClickHouse, MongoDB and ZooKeeper clusters under per-user permissions.",
    },
    "tools/postgres/pg_compacttable.md": {
        title: "pgcompacttable — remove PostgreSQL table and index bloat online",
        description: "Measure and reduce PostgreSQL table and index bloat without heavy locks using pgcompacttable from Ivory, with background jobs, live output and bloat estimation queries.",
    },
}

rmSync(target, {recursive: true, force: true})
mkdirSync(target, {recursive: true})

const files = readdirSync(source, {recursive: true}).map(f => f.split("\\").join("/"))
for (const file of files) {
    if (!file.endsWith(".md")) {
        if (file.includes(".")) cpSync(join(source, file), join(target, file))
        continue
    }
    const page = meta[file]
    if (!page) throw new Error(`site/scripts/sync-docs.mjs: no title/description for .doc/${file}`)
    const body = readFileSync(join(source, file), "utf8")
        .replace(/^# \[☰]\([^)]*README\.md\) /m, "# ")
        .replace(/]\(([^)\s]+)\)/g, (match, link) => `](${rewriteLink(file, link)})`)
    const frontmatter = `---\ntitle: ${JSON.stringify(page.title)}\ndescription: ${JSON.stringify(page.description)}\n---\n\n`
    mkdirSync(dirname(join(target, file)), {recursive: true})
    writeFileSync(join(target, file), frontmatter + body)
}

for (const file of Object.keys(meta)) {
    if (!files.includes(file)) {
        throw new Error(`site/scripts/sync-docs.mjs: .doc/${file} has a title/description but no longer exists`)
    }
}

mkdirSync(publicDir, {recursive: true})
cpSync(join(repo, "app", "shared", "assets", "ivory.png"), join(publicDir, "ivory.png"))
cpSync(join(repo, "app", "shared", "assets", "ivory-fav.png"), join(publicDir, "ivory-fav.png"))
cpSync(join(source, "images", "demo.png"), join(publicDir, "og-image.png"))

function rewriteLink(file, link) {
    if (/^[a-z]+:|^#|^\//.test(link)) return link
    const resolved = posix.normalize(posix.join(posix.dirname(file), link))
    if (resolved.startsWith("../")) return `${github}/${resolved.slice(3)}`
    return link.startsWith(".") ? link : `./${link}`
}
