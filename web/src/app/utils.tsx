import {blue, green, indigo, orange, purple} from "@mui/material/colors";
import {
    ColorsMap,
    CredentialType,
    DefaultInstance,
    JobStatus,
    InstanceMap,
    FileUsageType,
    CertType,
    EnumOptions, Sidecar, Settings
} from "./types";
import {
    FilePresentOutlined,
    HeartBroken, InfoTwoTone, Key,
    LockTwoTone,
    MenuOpen, SecurityTwoTone,
    Shield,
    Storage,
    UploadFileOutlined
} from "@mui/icons-material";
import {AxiosError} from "axios";
import {SxProps, Theme} from "@mui/material";

export const InstanceColor: { [key: string]: string } = {
    master: green[500],
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

export const CredentialOptions: { [key in CredentialType]: EnumOptions } = {
    [CredentialType.POSTGRES]: {name: "POSTGRES", label: "Postgres Password", color: blue[300], icon: <Storage/>},
    [CredentialType.PATRONI]: {name: "PATRONI", label: "Patroni Password", color: green[300], icon: <HeartBroken/>}
}

export const CertOptions: { [key in CertType]: EnumOptions } = {
    [CertType.CLIENT_CA]: {name: "CLIENT_CA", label: "Client CA", color: purple[300], icon: <Shield/>, badge: "CA"},
    [CertType.CLIENT_CERT]: {name: "CLIENT_CERT", label: "Client Cert", color: purple[300], icon: <Shield/>, badge: "C"},
    [CertType.CLIENT_KEY]: {name: "CLIENT_KEY", label: "Client Key", color: purple[300], icon: <Shield/>, badge: "K"}
}

export const FileUsageOptions: { [key in FileUsageType]: EnumOptions } = {
    [FileUsageType.UPLOAD]: {name: "UPLOAD", label: "Cert Upload", color: indigo[300], icon: <UploadFileOutlined/>},
    [FileUsageType.PATH]: {name: "PATH", label: "Cert Path", color: indigo[300], icon: <FilePresentOutlined/>},
}

export const SettingOptions: { [key in Settings]: EnumOptions } = {
    [Settings.MENU]: {name: "MENU", label: "Settings", icon: <MenuOpen/>},
    [Settings.PASSWORD]: {name: "PASSWORD", label: "Password Manager", icon: <LockTwoTone/>},
    [Settings.CERTIFICATE]: {name: "CERTIFICATE", label: "Certificate Manager", icon: <SecurityTwoTone/>},
    [Settings.SECRET]: {name: "SECRET", label: "Secret Manager", icon: <Key/>},
    [Settings.ABOUT]: {name: "ABOUT", label: "About", icon: <InfoTwoTone/>},
}

export const createInstanceColors = (instances: InstanceMap) => {
    return Object.values(instances).reduce(
        (map, instance) => {
            if (instance.inCluster) {
                const domain = getDomain(instance.sidecar)
                map[domain] = instance.leader ? "success" : "primary"
            }
            return map
        },
        {} as ColorsMap
    )
}

export const initialInstance = (sidecar?: Sidecar): DefaultInstance => {
    const defaultSidecar = sidecar ?? {host: "-", port: 8008}
    return ({
        state: "-",
        role: "unknown",
        lag: -1,
        sidecar: defaultSidecar,
        database: {host: "-", port: 0},
        leader: false,
        inInstances: true,
        inCluster: false
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

export const getHostAndPort = (domain: string) => {
    const [host, port] = domain.split(":")
    return {host, port: port ? parseInt(port) : 8008}
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
        if (error.response) message = error.response.data["error"]
        else if (error.message) message = error.message
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
    return [sx1, ...(Array.isArray(sx2) ? sx2 : [sx2])]
}
