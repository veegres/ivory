import {blue, green} from "@mui/material/colors";
import {JobStatus} from "./types";

export const nodeColor: { [key: string]: string } = {
    master: green[500],
    leader: green[500],
    replica: blue[500]
}

export const jobStatus: { [key: number]: { name: string, color: string } } = {
    [JobStatus.PENDING]: {name: "PENDING", color: "#b9b9b9"},
    [JobStatus.RUNNING]: {name: "RUNNING", color: "#b97800"},
    [JobStatus.FINISHED]: {name: "FINISHED", color: "#00b919"},
    [JobStatus.FAILED]: {name: "FAILED", color: "#d20000"},
}

export const isJobEnded = (status: JobStatus) => status === JobStatus.FAILED || status === JobStatus.FINISHED
