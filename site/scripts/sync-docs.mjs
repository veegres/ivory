// Copies the repository's .doc/ (and the root community files) into the site, so the Markdown stays readable on GitHub as the single source.
import {cpSync, mkdirSync, readdirSync, readFileSync, rmSync, writeFileSync} from "node:fs"
import {dirname, join, posix} from "node:path"
import {fileURLToPath} from "node:url"

const site = join(dirname(fileURLToPath(import.meta.url)), "..")
const repo = join(site, "..")
const source = join(repo, ".doc")
const target = join(site, "content", "docs")
const publicDir = join(site, "content", "public")
const github = "https://github.com/veegres/ivory/blob/master"

// Root files published beside .doc/, keyed by their repository path; the value is their path under docs/.
const community = {
    "CONTRIBUTING.md": "community/contributing.md",
    "SECURITY.md": "community/security.md",
    "CODE_OF_CONDUCT.md": "community/code-of-conduct.md",
}

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
    "databases.md": {
        title: "Supported databases — PostgreSQL, Patroni, etcd, Redis, ClickHouse, MongoDB, ZooKeeper",
        description: "Every database and cluster type Ivory supports, compared: HA operations, replication lag, health checks and query console for each engine.",
    },
    "databases/patroni.md": {
        title: "PostgreSQL + Patroni — Patroni GUI for HA PostgreSQL clusters",
        description: "How Ivory manages a Patroni PostgreSQL cluster: REST API discovery, roles and lag, switchover, failover, restart, reinitialise, pause, and Spilo deployment templates.",
    },
    "databases/postgres.md": {
        title: "PostgreSQL — primary/standby replication without Patroni",
        description: "How Ivory monitors and operates plain PostgreSQL streaming replication: roles, WAL lag, replication slot warnings, promotion, reload and deployment templates.",
    },
    "databases/etcd.md": {
        title: "etcd — etcd cluster GUI, members, leader and raft lag",
        description: "How Ivory manages an etcd cluster: member discovery, raft lag, alarms, moving the leader, an etcdctl-style console and deployment templates.",
    },
    "databases/redis.md": {
        title: "Redis — Redis replication monitoring and command console",
        description: "How Ivory monitors Redis master/replica replication: link health, persistence warnings, replica promotion, a Redis command console and deployment templates.",
    },
    "databases/clickhouse.md": {
        title: "ClickHouse — replicated ClickHouse cluster monitoring",
        description: "How Ivory monitors a replicated ClickHouse cluster: membership from system.clusters, replication delay and queue, coordination sessions, SQL console and templates.",
    },
    "databases/mongodb.md": {
        title: "MongoDB — replica set monitoring and step-down",
        description: "How Ivory manages a MongoDB replica set: member states, replication lag, heartbeats, primary step-down, config edits, a shell-style console and templates.",
    },
    "databases/zookeeper.md": {
        title: "ZooKeeper — ensemble monitoring and znode console",
        description: "How Ivory monitors a ZooKeeper ensemble: leader and followers, sync status, four-letter-word commands, a znode console and deployment templates.",
    },
    "tools/postgres/pg_compacttable.md": {
        title: "pgcompacttable — remove PostgreSQL table and index bloat online",
        description: "Measure and reduce PostgreSQL table and index bloat without heavy locks using pgcompacttable from Ivory, with background jobs, live output and bloat estimation queries.",
    },
    "community/contributing.md": {
        title: "Contributing to Ivory",
        description: "How to contribute to Ivory: picking an issue, setting up the local environment, the project's conventions, adding a database plugin and sending a pull request.",
    },
    "community/security.md": {
        title: "Security policy",
        description: "Which Ivory versions receive security fixes and how to report a vulnerability privately.",
    },
    "community/code-of-conduct.md": {
        title: "Code of conduct",
        description: "The Contributor Covenant code of conduct for everyone taking part in the Ivory community.",
    },
}

rmSync(target, {recursive: true, force: true})
mkdirSync(target, {recursive: true})

const files = readdirSync(source, {recursive: true}).map(f => f.split("\\").join("/"))
const pages = {}
const assets = {}
for (const file of files) {
    if (file.endsWith(".md")) pages[`.doc/${file}`] = file
    else if (file.includes(".")) {
        cpSync(join(source, file), join(target, file))
        assets[`.doc/${file}`] = file
    }
}
Object.assign(pages, community)

for (const [from, to] of Object.entries(pages)) {
    const page = meta[to]
    if (!page) throw new Error(`site/scripts/sync-docs.mjs: no title/description for ${from}`)
    const body = readFileSync(join(repo, from), "utf8")
        .replace(/^# \[☰]\([^)]*README\.md\) /m, "# ")
        .replace(/]\(([^)\s]+)\)/g, (match, link) => `](${rewriteLink(from, to, link)})`)
    const frontmatter = `---\ntitle: ${JSON.stringify(page.title)}\ndescription: ${JSON.stringify(page.description)}\n---\n\n`
    mkdirSync(dirname(join(target, to)), {recursive: true})
    writeFileSync(join(target, to), frontmatter + body)
}

const published = Object.values(pages)
for (const file of Object.keys(meta)) {
    if (!published.includes(file)) {
        throw new Error(`site/scripts/sync-docs.mjs: docs/${file} has a title/description but no longer exists`)
    }
}

mkdirSync(publicDir, {recursive: true})
cpSync(join(repo, "app", "shared", "assets", "ivory.png"), join(publicDir, "ivory.png"))
cpSync(join(repo, "app", "shared", "assets", "ivory-fav.png"), join(publicDir, "ivory-fav.png"))
cpSync(join(source, "images", "demo.png"), join(publicDir, "og-image.png"))

// A link to a file the site publishes stays on the site; anything else in the repository goes to GitHub.
function rewriteLink(from, to, link) {
    if (/^[a-z]+:|^#|^\//.test(link)) return link
    const [path, anchor] = link.split("#")
    const resolved = posix.normalize(posix.join(posix.dirname(from), path))
    const hash = anchor ? `#${anchor}` : ""
    if (resolved === "README.md") return `/${hash}`
    const page = pages[resolved] ?? assets[resolved]
    if (page) {
        const relative = posix.relative(posix.dirname(to), page)
        return (relative.startsWith(".") ? relative : `./${relative}`) + hash
    }
    if (resolved.startsWith("../")) throw new Error(`site/scripts/sync-docs.mjs: ${from} links outside the repository: ${link}`)
    return `${github}/${resolved}${hash}`
}
