import {blue, green, orange} from "@mui/material/colors";
import {ColorsMap, CredentialType, InstanceLocal, JobStatus, InstanceMap} from "./types";
import {ReactElement} from "react";
import {HeartBroken, Storage} from "@mui/icons-material";
import {AxiosError} from "axios";

export const InstanceColor: { [key: string]: string } = {
    master: green[500],
    leader: green[500],
    replica: blue[500],
    unknown: orange[500]
}

export const JobOptions: { [key in JobStatus]: { name: string, color: string, active: boolean } } = {
    [JobStatus.PENDING]: {name: "PENDING", color: "#b9b9b9", active: true},
    [JobStatus.UNKNOWN]: {name: "UNKNOWN", color: "#5b3b00", active: false},
    [JobStatus.RUNNING]: {name: "RUNNING", color: "#b97800", active: true},
    [JobStatus.FINISHED]: {name: "FINISHED", color: "#00b919", active: false},
    [JobStatus.FAILED]: {name: "FAILED", color: "#d20000", active: false},
    [JobStatus.STOPPED]: {name: "STOPPED", color: "#b9b9b9", active: false},
}

export const CredentialOptions: { [key in CredentialType]: { name: string, color: string, icon: ReactElement } } = {
    [CredentialType.POSTGRES]: {name: "POSTGRES", color: blue[300], icon: <Storage />},
    [CredentialType.PATRONI]: {name: "PATRONI", color: green[300], icon: <HeartBroken />}
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

export const initialInstance: (domain?: string) => InstanceLocal = (domain = "-") => {
    return ({
        state: "-",
        role: "unknown",
        lag: -1,
        sidecar: getHostAndPort(domain),
        database: {host: "-", port: 0},
        leader: false,
        inInstances: true,
        inCluster: false
    });
}

/**
 * Combine instances from patroni and from ivory
 */
export const combineInstances = (instanceNames: string[], instanceInCluster: InstanceMap): [InstanceMap, boolean] => {
    const map: InstanceMap = {}
    let warning: boolean = false

    for (const key in instanceInCluster) {
        if (instanceNames.includes(key)) {
            map[key] = { ...instanceInCluster[key], inInstances: true }
        } else {
            map[key] = { ...instanceInCluster[key], inInstances: false }
        }
    }

    for (const value of instanceNames) {
        if (!map[value]) {
            map[value] = initialInstance(value)
        }
    }

    for (const key in map) {
        const value = map[key]
        if (!value.inInstances || !value.inCluster) {
            warning = true
        }
    }

    return [map, warning]
}

export const getDomain = ({ host, port }: { host: string, port: number }) => {
    return `${host.toLowerCase()}${port ? `:${port}` : "" }`
}

export const getHostAndPort = (domain: string) => {
    const [host, port] = domain.split(":")
    return { host, port: port ? parseInt(port) : 0 }
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
