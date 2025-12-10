import {
    FilePresentOutlined,
    HeartBroken,
    InfoTwoTone,
    Key,
    LockTwoTone,
    MenuOpen,
    Pause,
    PlayArrow,
    SecurityTwoTone,
    Shield,
    Storage,
    UploadFileOutlined
} from "@mui/icons-material"
import {SxProps, Theme} from "@mui/material"
import {blue, green, indigo, orange, red} from "@mui/material/colors"
import {materialDarkInit, materialLightInit} from "@uiw/codemirror-theme-material"
import {AxiosError} from "axios"
import dayjs from "dayjs"

import {JobStatus} from "../api/bloat/job/type"
import {CertType, FileUsageType} from "../api/cert/type"
import {ActiveCluster, ActiveInstance, Cluster} from "../api/cluster/type"
import {InstanceMap, InstanceRequest, InstanceWeb, Role, Sidecar, SidecarStatus} from "../api/instance/type"
import {PasswordType} from "../api/password/type"
import {Database, QueryConnection, QueryVariety} from "../api/query/type"
import {ColorsMap, EnumOptions, Links, Settings, SxPropsMap} from "./type"

export const IvoryLinks: Links = {
    git: {name: "Github", link: "https://github.com/veegres/ivory"},
    docs: {name: "Docs", link: "https://github.com/veegres/ivory/blob/master/README.md"},
    repository: {name: "Repository", link: "https://hub.docker.com/r/aelsergeev/ivory"},
    issues: {name: "Issues", link: "https://github.com/veegres/ivory/issues"},
    release: {name: "Releases", link: "https://github.com/veegres/ivory/releases"},
    sponsorship: {name: "Sponsorship", link: "https://patreon.com/anselvo"}
}

export const InstanceColor: { [key in Role]: string } = {
    leader: green[600],
    replica: blue[500],
    unknown: orange[500],
}

export const JobOptions: { [key in JobStatus]: { name: string, color: string, active: boolean } } = {
    [JobStatus.PENDING]: {name: "PENDING", color: "rgb(169,169,169)", active: true},
    [JobStatus.UNKNOWN]: {name: "UNKNOWN", color: "rgb(91,59,0)", active: false},
    [JobStatus.RUNNING]: {name: "RUNNING", color: "rgba(255,166,0,0.9)", active: true},
    [JobStatus.FINISHED]: {name: "FINISHED", color: "rgba(0,185,25,0.9)", active: false},
    [JobStatus.FAILED]: {name: "FAILED", color: "rgba(210,0,0,0.9)", active: false},
    [JobStatus.STOPPED]: {name: "STOPPED", color: "rgb(185,185,185)", active: false},
}

export const CredentialOptions: { [key in PasswordType]: EnumOptions } = {
    [PasswordType.POSTGRES]: {name: "POSTGRES", label: "Postgres Password", icon: <Storage/>, key: "postgresId"},
    [PasswordType.PATRONI]: {name: "PATRONI", label: "Patroni Password", icon: <HeartBroken/>, key: "patroniId"}
}

export const SidecarStatusOptions: { [key in SidecarStatus]: EnumOptions } = {
    [SidecarStatus.Active]: {name: "ACTIVE", label: "Activate Sidecar", icon: <Pause/>, color: green[600], key: "active"},
    [SidecarStatus.Paused]: {name: "PAUSED", label: "Pause Sidecar", icon: <PlayArrow/>, color: orange[500], key: "paused"}
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
    [Settings.PASSWORD]: {name: "PASSWORD", label: "Password Manager", icon: <LockTwoTone/>, key: "password"},
    [Settings.CERTIFICATE]: {name: "CERTIFICATE", label: "Certificate Manager", icon: <SecurityTwoTone/>, key: "cert"},
    [Settings.SECRET]: {name: "SECRET", label: "Secret Manager", icon: <Key/>, key: "secret"},
    [Settings.ABOUT]: {name: "ABOUT", label: "About", icon: <InfoTwoTone/>, key: "about"},
}

export const QueryVarietyOptions: { [key in QueryVariety]: EnumOptions } = {
    [QueryVariety.DatabaseSensitive]: {key: "DatabaseSensitive", label: "Database Sensitive", badge: "DS", color: red[900], icon: <></>},
    [QueryVariety.MasterOnly]: {key: "MasterOnly", label: "Master Only", badge: "MO", color: green[900], icon: <></>},
    [QueryVariety.ReplicaRecommended]: {key: "ReplicaRecommended", label: "Replica Recommended", badge: "RR", color: blue[900], icon: <></>},
}

export const createInstanceColors = (instances: InstanceMap) => {
    return Object.values(instances).reduce(
        (map, instance) => {
            const domain = getDomain(instance.sidecar)
            map[domain] = !instance.inCluster || !instance.inInstances ? "warning" : (instance.leader ? "success" : "primary")
            return map
        },
        {} as ColorsMap
    )
}

export const initialInstance = (sidecar?: Sidecar): InstanceWeb => {
    const defaultSidecar = sidecar ?? {host: "-", port: 8008}
    return ({
        state: "-",
        role: "unknown",
        lag: -1,
        pendingRestart: false,
        sidecar: defaultSidecar,
        database: {host: "-", port: 0},
        leader: false,
        inInstances: true,
        inCluster: false,
    })
}

export const isSidecarEqual = (sidecar1?: Sidecar, sidecar2?: Sidecar): boolean => {
    return sidecar1?.host === sidecar2?.host && sidecar1?.port === sidecar2?.port
}

/**
 * Combine instances from patroni and from ivory
 */
export const combineInstances = (instanceNames: Sidecar[], instanceInCluster: InstanceMap) => {
    const map: InstanceMap = {}
    let warning: boolean = false

    for (const key in instanceInCluster) {
        if (getDomains(instanceNames).includes(key)) {
            map[key] = {...instanceInCluster[key], inInstances: true}
        } else {
            map[key] = {...instanceInCluster[key], inInstances: false}
        }
    }

    for (const value of instanceNames) {
        const domain = getDomain(value)
        if (!map[domain]) {
            map[domain] = initialInstance(value)
        }
    }

    for (const key in map) {
        const value = map[key]
        if (!value.inInstances || !value.inCluster) {
            warning = true
        }
    }

    return {
        combinedInstanceMap: map,
        warning,
    }
}

export const getDomain = ({host, port}: Sidecar) => {
    return `${host.toLowerCase()}${port ? `:${port}` : ""}`
}

export const getDomains = (sidecars: Sidecar[]) => {
    return sidecars.map(value => getDomain(value))
}

export const getActiveInstance = (activeInstance: ActiveInstance, activeCluster?: ActiveCluster) => {
    return activeCluster && activeInstance[activeCluster.cluster.name]
}

export const getIsActiveInstance = (key: string, instance?: InstanceWeb) => {
    return instance ? getDomain(instance.sidecar) === key : false
}

export function getQueryConnection(cluster: Cluster, db: Database): QueryConnection {
    const credentialId = cluster.credentials.postgresId
    const certs = cluster.tls.database ? cluster.certs : undefined
    return {db, certs, credentialId}
}

export function getSidecarConnection(cluster: Cluster, sidecar: Sidecar): InstanceRequest {
    const credentialId = cluster.credentials.patroniId
    const certs = cluster.tls.sidecar ? cluster.certs : undefined
    return {sidecar, certs, credentialId}
}

export const getSidecar = (domain: string): Sidecar => {
    const [host, port] = domain.split(":")
    return {host, port: port ? parseInt(port) : 8008}
}

export const getSidecars = (domains: string[]): Sidecar[] => {
    return domains.map(value => getSidecar(value))
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
     * This function is needed to fix typescript issues when
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

