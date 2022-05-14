import {blue, green} from "@mui/material/colors";
import {ColorsMap, JobStatus, Node} from "./types";

export const nodeColor: { [key: string]: string } = {
    master: green[500],
    leader: green[500],
    replica: blue[500]
}

export const jobStatus: { [key: number]: { name: string, color: string, active: boolean } } = {
    [JobStatus.PENDING]: {name: "PENDING", color: "#b9b9b9", active: true},
    [JobStatus.UNKNOWN]: {name: "UNKNOWN", color: "#5b3b00", active: false},
    [JobStatus.RUNNING]: {name: "RUNNING", color: "#b97800", active: true},
    [JobStatus.FINISHED]: {name: "FINISHED", color: "#00b919", active: false},
    [JobStatus.FAILED]: {name: "FAILED", color: "#d20000", active: false},
    [JobStatus.STOPPED]: {name: "STOPPED", color: "#b9b9b9", active: false},
}

export const createColorsMap = (nodes?: Node[]) => {
    return nodes?.reduce(
        (map, node) => {
            const color = node.isLeader ? "success" : "primary"
            map[node.host] = color
            map[node.api_domain] = color
            return map
        },
        {} as ColorsMap
    )
}

export const activeNode = (nodes?: Node[]) => {
    const leader = nodes?.find(node => node.isLeader)
    return leader ? leader.api_domain : ''
}

export const getPatroniDomain = (url: string) => url.split('/')[2]


export const unicodeAnimal = [
    "🐘", "🐇", "🐈", "🐋", "🐒", "🐢", "🐣", "🐬", "🐉",
    "🐩", "🦄", "🦥", "🦫", "🦭", "🦋", "🦉", "🦎", "🦙",
    "🦦", "🦢", "🦤", "🦞", "🦒", "🦕", "🦔", "🦌", "🦜",
]
export const randomUnicodeAnimal = () => {
    return unicodeAnimal[Math.floor(Math.random() * unicodeAnimal.length)]
}
