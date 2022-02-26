import {blue, green} from "@mui/material/colors";
import {JobStatus} from "./types";

export const nodeColor: { [key: string]: string } = {
    master: green[500],
    leader: green[500],
    replica: blue[500]
}

export const jobStatus: { [key: number]: { name: string, color: string, active: boolean } } = {
    [JobStatus.PENDING]: {name: "PENDING", color: "#b9b9b9", active: true},
    [JobStatus.RUNNING]: {name: "RUNNING", color: "#b97800", active: true},
    [JobStatus.FINISHED]: {name: "FINISHED", color: "#00b919", active: false},
    [JobStatus.FAILED]: {name: "FAILED", color: "#d20000", active: false},
    [JobStatus.STOPPED]: {name: "STOPPED", color: "#b9b9b9", active: false},
}
