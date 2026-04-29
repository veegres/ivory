import {
    BackupTwoTone, Block, CheckCircleOutlined,
    FilePresentOutlined, HeartBroken, HelpOutline,
    InfoTwoTone, KeyTwoTone, LockTwoTone,
    MenuOpen, Pause, PlayArrow, RuleTwoTone,
    SecurityTwoTone, Shield, Storage, UploadFileOutlined,
} from "@mui/icons-material"
import {SxProps, Theme} from "@mui/material"
import {blue, green, indigo, orange, purple, red} from "@mui/material/colors"
import {materialDarkInit, materialLightInit} from "@uiw/codemirror-theme-material"
import {AxiosError} from "axios"
import dayjs from "dayjs"

import {JobStatus} from "../api/bloat/job/type"
import {CertType, FileUsageType} from "../api/cert/type"
import {Cluster, Node, NodeOverview} from "../api/cluster/type"
import {Database, DatabaseType, QueryVariety} from "../api/database/type"
import {KeeperType, Role, Status as KeeperStatus} from "../api/keeper/type"
import {Connection as NodeConnection, KeeperRequest} from "../api/node/type"
import {Status as PermissionStatus} from "../api/permission/type"
import {Connection as QueryConnection} from "../api/query/type"
import {VaultType} from "../api/vault/type"
import {EnumOptions, Links, Settings, SxPropsMap} from "./type"

export const IvoryLinks: Links = {
    git: {name: "Github", link: "https://github.com/veegres/ivory"},
    docs: {name: "Docs", link: "https://github.com/veegres/ivory/blob/master/README.md"},
    repository: {name: "Repository", link: "https://hub.docker.com/r/veegres/ivory"},
    issues: {name: "Issues", link: "https://github.com/veegres/ivory/issues"},
    release: {name: "Releases", link: "https://github.com/veegres/ivory/releases"},
    sponsorship: {name: "Sponsorship", link: "https://boosty.to/anselvo/purchase/1454406"}
}

export const NodeColor: { [key in Role]: { label: "success" | "primary" | "error" | "warning", color: string } } = {
    leader: {label: "success", color: green[600]},
    replica: {label: "primary", color: blue[500]},
    unknown: {label: "warning", color:  orange[500]},
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
    [VaultType.DATABASE_PASSWORD]: {name: "DATABASE_PASSWORD", label: "Database Password", icon: <Storage/>, key: "databaseId"},
    [VaultType.KEEPER_PASSWORD]: {name: "KEEPER_PASSWORD", label: "Keeper Password", icon: <HeartBroken/>, key: "keeperId"},
    [VaultType.SSH_PASSWORD]: {name: "SSH_PASSWORD", label: "SSH Password", icon: <LockTwoTone/>, key: "sshVaultId"},
    [VaultType.SSH_KEY]: {name: "SSH_KEY", label: "SSH Key", icon: <KeyTwoTone/>, key: "sshKeyId"},
}

export const KeeperStatusOptions: { [key in KeeperStatus]: EnumOptions } = {
    [KeeperStatus.Active]: {name: "ACTIVE", label: "Activate Keeper", icon: <Pause/>, color: green[600], key: "active"},
    [KeeperStatus.Paused]: {name: "PAUSED", label: "Pause Keeper", icon: <PlayArrow/>, color: orange[500], key: "paused"}
}

export const CertOptions: { [key in CertType]: EnumOptions } = {
    [CertType.CLIENT_CA]: {name: "CLIENT_CA", label: "Client CA", icon: <Shield/>, badge: "CA", key: "clientCAId"},
    [CertType.CLIENT_CERT]: {name: "CLIENT_CERT", label: "Client Cert", icon: <Shield/>, badge: "C", key: "clientCertId"},
    [CertType.CLIENT_KEY]: {name: "CLIENT_KEY", label: "Client Key", icon: <Shield/>, badge: "K", key: "clientKeyId"}
}

export const FileUsageOptions: { [key in FileUsageType]: EnumOptions } = {
    [FileUsageType.UPLOAD]: {name: "UPLOAD", label: "Cert Upload", color: indigo[300], icon: <UploadFileOutlined/>, key: "upload"},
    [FileUsageType.PATH]: {name: "PATH", label: "Cert Path", color: indigo[300], icon: <FilePresentOutlined/>, key: "path"},
}

export const SettingOptions: { [key in Settings]: EnumOptions } = {
    [Settings.MENU]: {name: "MENU", label: "Settings", icon: <MenuOpen/>, key: "menu"},
    [Settings.VAULT]: {name: "VAULT", label: "Vault Manager", icon: <LockTwoTone/>, key: "vault"},
    [Settings.CERTIFICATE]: {name: "CERTIFICATE", label: "Certificate Manager", icon: <SecurityTwoTone/>, key: "cert"},
    [Settings.PERMISSION]: {name: "PERMISSION", label: "Permission Manager", icon: <RuleTwoTone/>, key: "permission"},
    [Settings.SECRET]: {name: "SECRET", label: "Secret Manager", icon: <KeyTwoTone/>, key: "secret"},
    [Settings.BACKUP]: {name: "BACKUP", label: "Backup", icon: <BackupTwoTone/>, key: "backup"},
    [Settings.ABOUT]: {name: "ABOUT", label: "About", icon: <InfoTwoTone/>, key: "about"},
}

export const QueryVarietyOptions: { [key in QueryVariety]: EnumOptions } = {
    [QueryVariety.DatabaseSensitive]: {key: "DatabaseSensitive", label: "Database Sensitive", badge: "DS", color: red[900], icon: <></>},
    [QueryVariety.MasterOnly]: {key: "MasterOnly", label: "Master Only", badge: "MO", color: green[900], icon: <></>},
    [QueryVariety.ReplicaRecommended]: {key: "ReplicaRecommended", label: "Replica Recommended", badge: "RR", color: blue[900], icon: <></>},
}

export const PermissionOptions: { [key in PermissionStatus]: EnumOptions } = {
    [PermissionStatus.GRANTED]: {key: "Granted", label: "Granted", icon: <CheckCircleOutlined/>, color: "success.main"},
    [PermissionStatus.PENDING]: {key: "Pending", label: "Pending", icon: <HelpOutline/>, color: "secondary.main"},
    [PermissionStatus.NOT_PERMITTED]: {key: "Not permitted", label: "Not permitted", icon: <Block/>, color: "error.main"},
}

export const DefaultKeeperPorts = {
    [KeeperType.PATRONI]: 8008,
    [KeeperType.POSTGRES]: 5432,
}

export const DefaultDatabasePorts = {
    [DatabaseType.POSTGRES]: 5432,
    [DatabaseType.ETCD]: 2379,
}

export const initialNode = (kt: KeeperType, dbt: DatabaseType, domain: string): Node => {
    const connection = getNodeConnection(kt, dbt, domain)
    return ({
        connection: connection,
        warnings: ["no response from keeper"],
        keeper: {
            state: "-",
            role: "unknown",
            lag: -1,
            pendingRestart: false,
            discoveredHost: connection.host,
            discoveredKeeperPort: kt,
            discoveredDbPort: dbt,
        },
    })
}

export const isConnectionEqual = (c1?: NodeConnection, c2?: NodeConnection): boolean => {
    return c1?.host === c2?.host && c1?.keeperPort === c2?.keeperPort && c1?.sshKeyId === c2?.sshKeyId
}

export function getQueryConnection(cluster: Cluster, db: Database): QueryConnection {
    const vaultId = cluster.vaults.databaseId
    const certs = cluster.tls.database ? cluster.certs : undefined
    return {db, certs, vaultId}
}

export function getKeeperRequest(cluster: Cluster, host: string, port: number): KeeperRequest {
    const vaultId = cluster.vaults.keeperId
    const certs = cluster.tls.keeper ? cluster.certs : undefined
    return {host, port, certs, vaultId, type: cluster.keeperType}
}

export const getDomain = (connection: NodeConnection, simple: boolean = true) => {
    const host = connection.host
    const keeperPort = connection.keeperPort ? `:${connection.keeperPort}` : ""
    const dbPort = !simple && connection.dbPort ? `:${connection.dbPort}` : ""
    const sshPort = !simple && connection.sshPort ? `:${connection.sshPort}` : ""
    const sshKeyId = !simple && connection.sshKeyId ? `:${connection.sshKeyId}` : ""
    return `${host.toLowerCase()}${keeperPort}${dbPort}${sshPort}${sshKeyId}`
}

export const getDomains = (nodes: NodeConnection[], simple: boolean = true) => {
    return nodes.map(value => getDomain(value, simple))
}

export const getNodeConnection = (kt: KeeperType, dbt: DatabaseType, domain: string): NodeConnection => {
    const [host, keeperPort, dbPort, sshPort, sshKeyId] = domain.split(":")
    return {
        host: host.toLowerCase(),
        keeperPort: parseInt(keeperPort) || DefaultKeeperPorts[kt],
        dbPort: parseInt(dbPort) || DefaultDatabasePorts[dbt],
        sshPort: parseInt(sshPort) || 22,
        sshKeyId: sshKeyId,
    }
}

export const getNodeConnections = (kt: KeeperType, dbt: DatabaseType, domains: string[]): NodeConnection[] => {
    return domains.map(value => getNodeConnection(kt, dbt, value))
}

export const getMainKeeper = (nodes: NodeOverview = {}, detectedKeeper?: string): [string?, Node?] => {
    return Object.entries(nodes).find(([_, v]) => v.keeper.role === "leader") ?? [detectedKeeper, nodes[detectedKeeper ?? ""]]
}

export const getDetectionItems = (nodes: NodeOverview = {}, detectedKeeper?: string, manualKeeper?: string) => {
    const detection = manualKeeper ? "manual" : "auto"
    const [domain, node] = getMainKeeper(nodes, manualKeeper ?? detectedKeeper)
    const mainLabel = domain ?? "none"
    const mainRole = node?.keeper.role ?? "unknown"
    return [
        {title: "Detection", label: detection, bgColor: purple[400]},
        {title: "Main Keeper", label: mainLabel, bgColor: NodeColor[mainRole].color}
    ]
}

export const shortUuid = (uuid: string) => uuid.substring(0, 8)

export const unicodeAnimal = [
    "🐘", "🐇", "🐈", "🐋", "🐒", "🐢", "🐣", "🐬", "🐉",
    "🐩", "🦄", "🦥", "🦫", "🦭", "🦋", "🦉", "🦎", "🦙",
    "🦦", "🦢", "🦤", "🦞", "🦒", "🦕", "🦔", "🦌", "🦜",
]
export const randomUnicodeAnimal = () => {
    return unicodeAnimal[Math.floor(Math.random() * unicodeAnimal.length)]
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
     * This function is needed to fix TypeScript issues when
     * `sx` can be an array and `SxProps` can be an array type
     *
     * https://github.com/mui/material-ui/issues/29900
     *
     * @param sx1
     * @param sx2
     */
    merge: (sx1?: SxProps<Theme>, sx2?: SxProps<Theme>) => [...(Array.isArray(sx1) ? sx1 : [sx1]), ...(Array.isArray(sx2) ? sx2 : [sx2])],
    style: {
        paper: {backgroundImage: "linear-gradient(rgba(255, 255, 255, 0.09), rgba(255, 255, 255, 0.09))"},
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
