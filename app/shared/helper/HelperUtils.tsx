import {
    BackupTwoTone,
    Block,
    CheckCircleOutlined,
    DnsTwoTone,
    FilePresentOutlined,
    HeartBrokenTwoTone,
    HelpOutlined,
    InfoTwoTone,
    KeyTwoTone,
    LockTwoTone,
    MenuOpen,
    Pause,
    PlayArrow,
    RuleTwoTone,
    SecurityTwoTone,
    Shield,
    UploadFileOutlined,
} from "@mui/icons-material"
import {SxProps, Theme} from "@mui/material"
import {materialDarkInit, materialLightInit} from "@uiw/codemirror-theme-material"
import {AxiosError} from "axios"
import dayjs from "dayjs"

import {CertType, FileUsageType} from "../../features/cert/api/CertType"
import {Node, NodeConfig, NodeOverview, Options} from "../../features/cluster/api/ClusterType"
import {
    DeployVar,
    KeeperConnection,
    KeeperDeploySpecResponse,
    KeeperOneRequest,
    KeeperPlugin,
    KeeperState,
    KeeperStatus,
    PlatformPlugin,
    PlatformVaultConnection,
    ReleaseStage,
    Role,
} from "../../features/node/api/NodeType"
import {Status as PermissionStatus} from "../../features/permission/api/PermissionType"
import {Connection as QueryConnection, DbPlugin, VarietyType} from "../../features/query/api/QueryType"
import {VaultType} from "../../features/vault/api/VaultType"
import {JobStatus} from "../../tools/pg_compacttable/api/job/PgCompactTableJobType"
import {EnumOptions, Links, Settings, SxPropsMap} from "./HelperType"

export const IvoryLinks: Links = {
    git: {name: "Github", link: "https://github.com/veegres/ivory"},
    docs: {name: "Docs", link: "https://github.com/veegres/ivory/blob/master/README.md"},
    repository: {name: "Repository", link: "https://hub.docker.com/r/veegres/ivory"},
    issues: {name: "Issues", link: "https://github.com/veegres/ivory/issues"},
    release: {name: "Releases", link: "https://github.com/veegres/ivory/releases"},
    sponsorship: {name: "Sponsorship", link: "https://boosty.to/anselvo/purchase/1454406"}
}

export const NodeColor: { [key in Role]: { label: "success" | "primary" | "error" | "info" | "warning", color: string } } = {
    leader: {label: "success", color: "success.main"},
    replica: {label: "info", color: "info.main"},
    unknown: {label: "warning", color:  "warning.main"},
}

export const JobOptions: { [key in JobStatus]: { name: string, color: string, active: boolean } } = {
    [JobStatus.PENDING]: {name: "PENDING", color: "rgb(169,169,169)", active: true},
    [JobStatus.UNKNOWN]: {name: "UNKNOWN", color: "rgb(91,59,0)", active: false},
    [JobStatus.RUNNING]: {name: "RUNNING", color: "rgba(255,166,0,0.9)", active: true},
    [JobStatus.FINISHED]: {name: "FINISHED", color: "rgba(0,185,25,0.9)", active: false},
    [JobStatus.FAILED]: {name: "FAILED", color: "rgba(210,0,0,0.9)", active: false},
    [JobStatus.STOPPED]: {name: "STOPPED", color: "rgb(185,185,185)", active: false},
}

export const VaultOptions: { [key in VaultType]: EnumOptions } = {
    [VaultType.DATABASE_PASSWORD]: {name: "DATABASE_PASSWORD", label: "Database Password", icon: <DnsTwoTone/>, key: "databaseId"},
    [VaultType.KEEPER_PASSWORD]: {name: "KEEPER_PASSWORD", label: "Keeper Password", icon: <HeartBrokenTwoTone/>, key: "keeperId"},
    [VaultType.SSH_PASSWORD]: {name: "SSH_PASSWORD", label: "SSH Password", icon: <LockTwoTone/>, key: "sshVaultId"},
    [VaultType.SSH_KEY]: {name: "SSH_KEY", label: "SSH Key", icon: <KeyTwoTone/>, key: "sshKeyId"},
}

// DbModel classifies a plugin by the data model it exposes (OLTP, OLAP, ...),
// shown next to the plugin's release stage and used to sort the plugin
// picker - see OptionsPlugins.tsx and ListKeepers.tsx.
export type DbModel = "OLTP" | "OLAP" | "Key-Value" | "Document" | "DCS"

export const KeeperPluginOptions: { [key in KeeperPlugin]: EnumOptions & {dbPlugin: DbPlugin, stage: ReleaseStage, model: DbModel} } = {
    [KeeperPlugin.PATRONI_POSTGRES]: {name: "patroni", label: "Patroni Postgres", icon: <HeartBrokenTwoTone/>, key: "patroni_postgres", dbPlugin: DbPlugin.POSTGRES, stage: ReleaseStage.STABLE, model: "OLTP"},
    [KeeperPlugin.NATIVE_POSTGRES]: {name: "postgres", label: "Postgres", icon: <HeartBrokenTwoTone/>, key: "native_postgres", dbPlugin: DbPlugin.POSTGRES, stage: ReleaseStage.BETA, model: "OLTP"},
    [KeeperPlugin.NATIVE_MONGO]: {name: "mongo", label: "MongoDB", icon: <HeartBrokenTwoTone/>, key: "native_mongo", dbPlugin: DbPlugin.MONGO, stage: ReleaseStage.ALPHA, model: "Document"},
    [KeeperPlugin.NATIVE_REDIS]: {name: "redis", label: "Redis", icon: <HeartBrokenTwoTone/>, key: "native_redis", dbPlugin: DbPlugin.REDIS, stage: ReleaseStage.ALPHA, model: "Key-Value"},
    [KeeperPlugin.NATIVE_CLICKHOUSE]: {name: "clickhouse", label: "ClickHouse", icon: <HeartBrokenTwoTone/>, key: "native_clickhouse", dbPlugin: DbPlugin.CLICKHOUSE, stage: ReleaseStage.ALPHA, model: "OLAP"},
    [KeeperPlugin.NATIVE_ETCD]: {name: "etcd", label: "Etcd", icon: <HeartBrokenTwoTone/>, key: "native_etcd", dbPlugin: DbPlugin.ETCD, stage: ReleaseStage.ALPHA, model: "DCS"},
    [KeeperPlugin.NATIVE_ZOOKEEPER]: {name: "zookeeper", label: "ZooKeeper", icon: <HeartBrokenTwoTone/>, key: "native_zookeeper", dbPlugin: DbPlugin.ZOOKEEPER, stage: ReleaseStage.ALPHA, model: "DCS"},
}

export const ReleaseStageOptions: { [key in ReleaseStage]: {label: string, description: string, color: "success" | "warning" | "error"} } = {
    [ReleaseStage.STABLE]: {label: "STABLE", description: "Production ready", color: "success"},
    [ReleaseStage.BETA]: {label: "BETA", description: "Mostly stable, expect rough edges", color: "warning"},
    [ReleaseStage.ALPHA]: {label: "ALPHA", description: "Experimental, use with caution", color: "error"},
}

export const DbModelOptions: { [key in DbModel]: {label: string, description: string} } = {
    "OLTP": {label: "OLTP", description: "Relational, transactional workloads"},
    "OLAP": {label: "OLAP", description: "Columnar, analytical workloads"},
    "Key-Value": {label: "Key-Value", description: "In-memory key-value store"},
    "Document": {label: "Document", description: "Document-oriented store"},
    "DCS": {label: "DCS", description: "Distributed coordination / consensus store"},
}

export const DbPluginOptions: { [key in DbPlugin]: EnumOptions & {model: DbModel} } = {
    [DbPlugin.POSTGRES]: {name: "POSTGRES", label: "Postgres", icon: <DnsTwoTone/>, key: "postgres", model: "OLTP"},
    [DbPlugin.MONGO]: {name: "MONGO", label: "MongoDB", icon: <DnsTwoTone/>, key: "mongo", model: "Document"},
    [DbPlugin.REDIS]: {name: "REDIS", label: "Redis", icon: <DnsTwoTone/>, key: "redis", model: "Key-Value"},
    [DbPlugin.CLICKHOUSE]: {name: "CLICKHOUSE", label: "ClickHouse", icon: <DnsTwoTone/>, key: "clickhouse", model: "OLAP"},
    [DbPlugin.ETCD]: {name: "ETCD", label: "Etcd", icon: <DnsTwoTone/>, key: "etcd", model: "DCS"},
    [DbPlugin.ZOOKEEPER]: {name: "ZOOKEEPER", label: "ZooKeeper", icon: <DnsTwoTone/>, key: "zookeeper", model: "DCS"},
}

export const KeeperStatusOptions: { [key in KeeperStatus]: EnumOptions } = {
    [KeeperStatus.Active]: {name: "ACTIVE", label: "Activate Keeper", icon: <Pause/>, color: "success.main", key: "active"},
    [KeeperStatus.Paused]: {name: "PAUSED", label: "Pause Keeper", icon: <PlayArrow/>, color: "warning.main", key: "paused"}
}

export const KeeperStateOptions: { [key in KeeperState]: { color: "success" | "warning" | "error" | "default" } } = {
    running: {color: "success"},
    starting: {color: "warning"},
    restarting: {color: "warning"},
    stopping: {color: "warning"},
    stopped: {color: "error"},
    failed: {color: "error"},
    unreachable: {color: "error"},
    unknown: {color: "default"},
}

export const CertOptions: { [key in CertType]: EnumOptions } = {
    [CertType.CLIENT_CA]: {name: "CLIENT_CA", label: "Client CA", icon: <Shield/>, badge: "CA", key: "clientCAId"},
    [CertType.CLIENT_CERT]: {name: "CLIENT_CERT", label: "Client Cert", icon: <Shield/>, badge: "C", key: "clientCertId"},
    [CertType.CLIENT_KEY]: {name: "CLIENT_KEY", label: "Client Key", icon: <Shield/>, badge: "K", key: "clientKeyId"}
}

export const FileUsageOptions: { [key in FileUsageType]: EnumOptions } = {
    [FileUsageType.UPLOAD]: {name: "UPLOAD", label: "Cert Upload", color: "secondary.main", icon: <UploadFileOutlined/>, key: "upload"},
    [FileUsageType.PATH]: {name: "PATH", label: "Cert Path", color: "secondary.main", icon: <FilePresentOutlined/>, key: "path"},
}

export const SettingOptions: { [key in Settings]: EnumOptions } = {
    [Settings.MENU]: {name: "MENU", label: "SETTINGS", icon: <MenuOpen/>, key: "menu"},
    [Settings.VAULT]: {name: "VAULT", label: "Vault Manager", icon: <LockTwoTone/>, key: "vault"},
    [Settings.CERTIFICATE]: {name: "CERTIFICATE", label: "Certificate Manager", icon: <SecurityTwoTone/>, key: "cert"},
    [Settings.PERMISSION]: {name: "PERMISSION", label: "Permission Manager", icon: <RuleTwoTone/>, key: "permission"},
    [Settings.SECRET]: {name: "SECRET", label: "Secret Manager", icon: <KeyTwoTone/>, key: "secret"},
    [Settings.BACKUP]: {name: "BACKUP", label: "Backup", icon: <BackupTwoTone/>, key: "backup"},
    [Settings.ABOUT]: {name: "ABOUT", label: "About", icon: <InfoTwoTone/>, key: "about"},
}

export const QueryVarietyOptions: { [key in VarietyType]: EnumOptions } = {
    [VarietyType.DatabaseSensitive]: {key: "DatabaseSensitive", label: "Database Sensitive", badge: "DS", color: "error", icon: <></>},
    [VarietyType.MasterOnly]: {key: "MasterOnly", label: "Master Only", badge: "MO", color: "success", icon: <></>},
    [VarietyType.ReplicaRecommended]: {key: "ReplicaRecommended", label: "Replica Recommended", badge: "RR", color: "info", icon: <></>},
}

export const PermissionOptions: { [key in PermissionStatus]: EnumOptions } = {
    [PermissionStatus.GRANTED]: {key: "Granted", label: "Granted", icon: <CheckCircleOutlined/>, color: "success.main"},
    [PermissionStatus.PENDING]: {key: "Pending", label: "Pending", icon: <HelpOutlined/>, color: "secondary.main"},
    [PermissionStatus.NOT_PERMITTED]: {key: "Not permitted", label: "Not permitted", icon: <Block/>, color: "error.main"},
}

export const getInitialNode = (config: NodeConfig): Node => {
    return ({
        config: config,
        warnings: ["no response from keeper"],
        keeper: {state: "unknown", role: "unknown", sync: false, lag: -1, pendingRestart: false},
    })
}

export const isConnectionEqual = (c1?: NodeConfig, c2?: NodeConfig): boolean => {
    return c1?.host === c2?.host && c1?.keeperPort === c2?.keeperPort
}

export function getQueryConnection(options: Options, host: string, port?: number): QueryConnection | undefined {
    if (!port) return
    const vaultId = options.vaults.databaseId
    const db = {plugin: options.plugins.database, host, port}
    const certs = options.tls.database ? options.certs : undefined
    return {db, certs, vaultId}
}

export function getPlatformConnection(options: Options, host: string, port?: number): PlatformVaultConnection | undefined {
    const vaultId = options.vaults.sshKeyId
    if (!port || !vaultId) return
    return {host, port, vaultId}
}

export function getKeeperConnection(host: string, port?: number): KeeperConnection | undefined {
    if (!port) return
    return {host, port}
}

export function getKeeperOneRequest(options: Options, host: string, port?: number): KeeperOneRequest | undefined {
    const con = getKeeperConnection(host, port)
    if (!con) return
    const vaultId = options.vaults.keeperId
    const certs = options.tls.keeper ? options.certs : undefined
    return {...con, certs, vaultId, plugin: options.plugins.keeper}
}

export const getDomain = (config: NodeConfig, simple: boolean = false) => {
    const host = config.host
    const keeperPort = config.keeperPort ? `:${config.keeperPort}` : simple ? "" : ":"
    const dbPort = simple ? "" : config.dbPort ? `:${config.dbPort}` : ":"
    const sshPort = simple ? "" : config.sshPort ? `:${config.sshPort}` : ":"
    return `${host.toLowerCase()}${keeperPort}${dbPort}${sshPort}`
}

export const getDomains = (nodes: NodeConfig[], simple: boolean = false) => {
    return nodes.map(value => getDomain(value, simple))
}

// NodeInputFormat declares how a node domain string is parsed: without a
// keeper port the format degrades to host:dbPort:sshPort and the keeper port
// mirrors the db port (plugins with no separate keeper API port); missing
// segments fall back to the provided defaults
export interface NodeInputFormat {
    withKeeperPort: boolean,
    defaults: {keeperPort?: number, dbPort?: number, sshPort?: number},
}

export const getNodeConfig = (domain: string, format?: NodeInputFormat): NodeConfig => {
    const [host, second, third, fourth] = domain.split(":")
    if (!format) {
        return {
            host: host.toLowerCase(),
            keeperPort: parseInt(second) || undefined,
            dbPort: parseInt(third) || undefined,
            sshPort: parseInt(fourth) || undefined,
        }
    }
    const {withKeeperPort, defaults} = format
    const dbPort = parseInt(withKeeperPort ? third : second) || defaults.dbPort
    const sshPort = parseInt(withKeeperPort ? fourth : third) || defaults.sshPort
    const keeperPort = withKeeperPort ? parseInt(second) || defaults.keeperPort : dbPort
    return {host: host.toLowerCase(), keeperPort, dbPort, sshPort}
}

export const getNodeConfigs = (domains: string[], format?: NodeInputFormat): NodeConfig[] => {
    return domains.map(value => getNodeConfig(value, format))
}

// getKeeperDefaultPort is the endpoint Ivory itself dials: the plugin's own
// keeper port when it has one, otherwise the database port.
export const getKeeperDefaultPort = (spec: KeeperDeploySpecResponse): number => {
    return spec.keeperPort ?? spec.dbPort
}

export interface DeployVarMeta {
    label: string,
    // example is what the variable turns into, shown next to it so the list
    // explains itself
    example: string,
    // secret means the browser must never hold the value: a preview shows the
    // mask, and only the server ever substitutes the real one
    secret: boolean,
}

// DeployPasswordMask is what a preview shows in place of a secret. It is not a
// placeholder - the command really does get a value there, this is only what
// the screen is allowed to render.
export const DeployPasswordMask = "*****"

export const DeployVarOptions: { [key in DeployVar]: DeployVarMeta } = {
    [DeployVar.Cluster]: {label: "Cluster", example: "my-cluster", secret: false},
    [DeployVar.Name]: {label: "Node Name", example: "node-1", secret: false},
    [DeployVar.Host]: {label: "Host", example: "10.0.0.1", secret: false},
    [DeployVar.SshPort]: {label: "SSH Port", example: "22", secret: false},
    [DeployVar.KeeperPort]: {label: "Keeper Port", example: "8008", secret: false},
    [DeployVar.DbPort]: {label: "Database Port", example: "5432", secret: false},
    [DeployVar.DbUser]: {label: "Database User", example: "postgres", secret: false},
    [DeployVar.DbPass]: {label: "Database Password", example: DeployPasswordMask, secret: true},
}

export const PlatformPluginOptions: { [key in PlatformPlugin]: EnumOptions } = {
    [PlatformPlugin.LINUX]: {label: "LINUX", name: "Linux", key: PlatformPlugin.LINUX, icon: <DnsTwoTone/>},
}

// PlaceholderPattern is the shape of a variable, shared so the editor's
// highlighting marks exactly what interpolation will act on. Copy it with
// new RegExp before handing it to anything that keeps lastIndex.
export const PlaceholderPattern = /{{\w+}}/g

// getPlaceholders lists every {{variable}} a command references, in order and
// deduplicated.
export const getPlaceholders = (text: string): string[] => {
    return [...new Set(text.match(PlaceholderPattern) ?? [])]
}

// DeployValues is what a deploy supplies to a command. The username is here
// because the form already knows it, so a preview may as well read like the
// real thing; the password is only ever the mask, since the browser is never
// given the value the server substitutes.
export interface DeployValues {
    cluster?: string,
    name?: string,
    host?: string,
    sshPort?: number,
    keeperPort?: number,
    dbPort?: number,
    dbUser?: string,
    dbPass?: string,
}

// interpolateCommand fills a command with a node's values so the user can read
// what will actually run. It mirrors keeper.Interpolate: a missing or empty
// value leaves the placeholder in place rather than substituting nothing.
export const interpolateCommand = (text: string, values: DeployValues): string => {
    const lookup: {[key in DeployVar]?: string} = {
        [DeployVar.Cluster]: values.cluster,
        [DeployVar.Name]: values.name,
        [DeployVar.Host]: values.host,
        [DeployVar.SshPort]: values.sshPort?.toString(),
        [DeployVar.KeeperPort]: values.keeperPort?.toString(),
        [DeployVar.DbPort]: values.dbPort?.toString(),
        [DeployVar.DbUser]: values.dbUser,
        [DeployVar.DbPass]: values.dbPass,
    }
    return text.replace(PlaceholderPattern, match => lookup[match as DeployVar] || match)
}

// getUnknownPlaceholders reports the ones outside the closed vocabulary. They
// are a validation error, never a new variable.
export const getUnknownPlaceholders = (text: string): string[] => {
    const known = Object.values(DeployVar) as string[]
    return getPlaceholders(text).filter(p => !known.includes(p))
}

export const getMainKeeper = (nodes: NodeOverview = {}, manual?: string): [string?, Node?] => {
    if (manual) return [manual, nodes[manual]]
    const list = Object.entries(nodes)
    return list.find(([_, v]) => v.keeper.role === "leader")
        ?? list.find(([_, v]) => v.keeper.role === "replica")
        ?? [undefined, undefined]
}

export const getDetectionItems = (mainNode: [string?, Node?], manual: boolean) => {
    const [domain, node] = mainNode
    const mainLabel = domain ?? "none"
    const mainRole: Role = node?.keeper.role ?? "unknown"
    return [
        {title: "Detection", label: manual ? "manual" : "auto", color: "secondary"},
        {title: "Main Keeper", label: mainLabel, color: NodeColor[mainRole].label}
    ]
}

export const getShortUuid = (uuid: string) => uuid.substring(0, 8)

export const UnicodeAnimal = [
    "🐘", "🐇", "🐈", "🐋", "🐒", "🐢", "🐣", "🐬", "🐉",
    "🐩", "🦄", "🦥", "🦫", "🦭", "🦋", "🦉", "🦎", "🦙",
    "🦦", "🦢", "🦤", "🦞", "🦒", "🦕", "🦔", "🦌", "🦜",
]
export const randomUnicodeAnimal = () => {
    return UnicodeAnimal[Math.floor(Math.random() * UnicodeAnimal.length)]
}

export const getErrorMessage = (error: any): string => {
    let message = "unknown"
    if (error instanceof AxiosError) {
        if (error.response) {
            if (error.response.data) {
                if (error.response.data["error"] !== undefined) message = error.response.data["error"]
                else message = error.response.data
            } else {
                message = `${error.response.status} ${error.response.statusText}`
            }
        } else {
            if (error.message) message = error.message
            else message = "unknown"
        }
    }
    if (typeof error === "string") message = error
    return message
}

export const getPostgresUrl = (con: QueryConnection) => {
    return `postgres://${con.db.host}:${con.db.port}/${con.db.name ?? "postgres"}`
}

// CodeMirror theme
export const CodeThemes = {
    dark: materialDarkInit({settings: {background: "transparent", gutterActiveForeground: "rgba(255,255,255,0.3)", selection: "rgba(255,255,255,0.1)"}}),
    light: materialLightInit({settings: {background: "transparent", gutterActiveForeground: "rgba(0,0,0,0.3)", selection: "rgba(0,0,0,0.1)"}}),
}

export const SxPropsFormatter = {
    /**
     * Merges two `SxProps` values into a flat sx array. A plain `sx={[a, b]}`
     * doesn't compile when an element is itself typed `SxProps<Theme>` (it may
     * be an array, and `SxProps` doesn't allow nested arrays). MUI closed
     * https://github.com/mui/material-ui/issues/29900 without changing the
     * type — the docs officially recommend this `Array.isArray` spread instead.
     */
    merge: (sx1?: SxProps<Theme>, sx2?: SxProps<Theme>) => [...(Array.isArray(sx1) ? sx1 : [sx1]), ...(Array.isArray(sx2) ? sx2 : [sx2])],
    style: {
        paper: {backgroundImage: "linear-gradient(rgba(255, 255, 255, 0.09), rgba(255, 255, 255, 0.09))"},
        pageMargin: {margin: {xs: "0 8px", sm: "0 5%"}},
        bgImageError: (theme) => ({backgroundImage: `linear-gradient(${theme.palette.error.dark}12, ${theme.palette.error.dark}12)`}),
        bgImageSelected: (theme) => ({backgroundImage: `linear-gradient(${theme.palette.action.hover}, ${theme.palette.action.hover})`}),
    } as SxPropsMap
}

export const DateTimeFormatter = {
    format: "YYYY-MM-DD HH:mm",
    formatWithTimezone: "YYYY-MM-DD HH:mm Z",
    utc: (value: string) => dayjs.utc(value).local().format(DateTimeFormatter.formatWithTimezone)
}

export const SizeFormatter = {
    format: Intl.NumberFormat("en", {notation: "compact", style: "unit", unit: "byte", unitDisplay: "narrow"}),
    pretty: (size: number) => SizeFormatter.format.format(size)
}

export const printLogs = (title: string, logs: string[]) => {
    // 1. Convert the array into HTML list items safely
    const listItemsHTML = logs
        .map(item => `<div>${item}</div>`)
        .join("")

    // 2. Create the hidden iframe
    const iframe = document.createElement("iframe")
    iframe.style.position = "fixed"
    iframe.style.width = "0"
    iframe.style.height = "0"
    iframe.style.border = "0"

    // Build the complete HTML document string
    // 3. Use srcdoc instead of document.write()
    iframe.srcdoc = `
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <title>${title}</title>
            <style>
                body { font-family: monospace; }
                h3 { color: #2c3e50; border-bottom: 2px solid #ecf0f1; padding-bottom: 10px }
            </style>
        </head>
        <body>
            <h3>${title}</h3>
            <div>
                ${listItemsHTML}
            </div>
        </body>
        </html>
    `

    // Append it to the document
    document.body.appendChild(iframe)

    // 4. Wait for the iframe to load before executing window functions
    iframe.onload = () => {
        const contentWindow = iframe.contentWindow

        // Explicit type guard check eliminates the 'possibly null' warning
        if (contentWindow) {
            contentWindow.focus()
            contentWindow.print()
        }

        // Cleanup DOM after a short timeout
        setTimeout(() => document.body.removeChild(iframe), 1000)
    }
}
