import {blue, green, indigo, orange, purple, red} from "@mui/material/colors";
import {
    FilePresentOutlined,
    HeartBroken,
    InfoTwoTone,
    Key,
    LockTwoTone,
    MenuOpen,
    SecurityTwoTone,
    Shield,
    Storage,
    UploadFileOutlined
} from "@mui/icons-material";
import {AxiosError} from "axios";
import {SxProps, Theme} from "@mui/material";
import {materialDarkInit, materialLightInit} from "@uiw/codemirror-theme-material";
import {PasswordType} from "../type/password";
import {ColorsMap, EnumOptions, FileUsageType, Links, Settings, Sidecar} from "../type/common";
import {CertType} from "../type/cert";
import {InstanceMap, InstanceWeb, Role} from "../type/instance";
import {JobStatus} from "../type/job";
import {QueryVariety} from "../type/query";

export const IvoryLinks: Links = {
    git: {name: "Github", link: "https://github.com/veegres/ivory"},
    docs: {name: "Docs", link: "https://github.com/veegres/ivory/blob/master/README.md"},
    repository: {name: "Repository", link: "https://hub.docker.com/r/aelsergeev/ivory"},
    issues: {name: "Issues", link: "https://github.com/veegres/ivory/issues"},
    release: {name: "Releases", link: "https://github.com/veegres/ivory/releases"},
    sponsorship: {name: "Sponsorship", link: "https://patreon.com/anselvo"}
}

export const InstanceColor: { [key in Role]: string } = {
    leader: green[500],
    replica: blue[500],
    unknown: orange[500],
}

export const JobOptions: { [key in JobStatus]: { name: string, color: string, active: boolean } } = {
    [JobStatus.PENDING]: {name: "PENDING", color: "#a9a9a9", active: true},
    [JobStatus.UNKNOWN]: {name: "UNKNOWN", color: "#5b3b00", active: false},
    [JobStatus.RUNNING]: {name: "RUNNING", color: "rgba(255,166,0,0.9)", active: true},
    [JobStatus.FINISHED]: {name: "FINISHED", color: "rgba(0,185,25,0.9)", active: false},
    [JobStatus.FAILED]: {name: "FAILED", color: "rgba(210,0,0,0.9)", active: false},
    [JobStatus.STOPPED]: {name: "STOPPED", color: "#b9b9b9", active: false},
}

export const CredentialOptions: { [key in PasswordType]: EnumOptions } = {
    [PasswordType.POSTGRES]: {name: "POSTGRES", label: "Postgres Password", color: blue[300], icon: <Storage/>, key: "postgresId"},
    [PasswordType.PATRONI]: {name: "PATRONI", label: "Patroni Password", color: green[300], icon: <HeartBroken/>, key: "patroniId"}
}

export const CertOptions: { [key in CertType]: EnumOptions } = {
    [CertType.CLIENT_CA]: {name: "CLIENT_CA", label: "Client CA", color: purple[300], icon: <Shield/>, badge: "CA", key: "clientCAId"},
    [CertType.CLIENT_CERT]: {name: "CLIENT_CERT", label: "Client Cert", color: purple[300], icon: <Shield/>, badge: "C", key: "clientCertId"},
    [CertType.CLIENT_KEY]: {name: "CLIENT_KEY", label: "Client Key", color: purple[300], icon: <Shield/>, badge: "K", key: "clientKeyId"}
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
    });
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

/**
 * This function is needed to fix typescript issues when
 * `sx` can be an array and `SxProps` can be an array type
 *
 * https://github.com/mui/material-ui/issues/29900
 *
 * @param sx1
 * @param sx2
 */
export const mergeSxProps = (sx1?: SxProps<Theme>, sx2?: SxProps<Theme>) => {
    return [...(Array.isArray(sx1) ? sx1 : [sx1]), ...(Array.isArray(sx2) ? sx2 : [sx2])]
}

// CodeMirror theme
export const CodeThemes = {
    dark: materialDarkInit({settings: {background: "transparent", gutterActiveForeground: "rgba(255,255,255,0.3)", selection: "rgba(255,255,255,0.1)"}}),
    light: materialLightInit({settings: {background: "transparent", gutterActiveForeground: "rgba(0,0,0,0.3)", selection: "rgba(0,0,0,0.1)"}}),
}
