import {blue, green, orange} from "@mui/material/colors";
import {ColorsMap, CredentialType, JobStatus, InstanceMap} from "./types";
import {ReactElement} from "react";
import {HeartBroken, Storage} from "@mui/icons-material";

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

export const createColorsMap = (nodes: InstanceMap) => {
    return Object.values(nodes).reduce(
        (map, node) => {
            const color = node.leader ? "success" : "primary"
            map[node.host.toLowerCase()] = color
            map[node.api_domain.toLowerCase()] = color
            return map
        },
        {} as ColorsMap
    )
}

export const getPatroniDomain = (url: string) => url.split('/')[2]

export const shortUuid = (uuid: string) => uuid.substring(0, 8)

export const unicodeAnimal = [
    "🐘", "🐇", "🐈", "🐋", "🐒", "🐢", "🐣", "🐬", "🐉",
    "🐩", "🦄", "🦥", "🦫", "🦭", "🦋", "🦉", "🦎", "🦙",
    "🦦", "🦢", "🦤", "🦞", "🦒", "🦕", "🦔", "🦌", "🦜",
]
export const randomUnicodeAnimal = () => {
    return unicodeAnimal[Math.floor(Math.random() * unicodeAnimal.length)]
}