import axios from "axios";
import {Cluster, NodePatroni, Node} from "./types";

const api = axios.create({ baseURL: '/api' })

export const nodeApi = {
    patroni: (node: String) => api.get<{ response: NodePatroni }>(`/node/${node}/patroni`).then((response) => response.data.response),
    cluster: (node: String) => api.get<{ response: { members: Node[] } }>(`/node/${node}/cluster`).then((response) => response.data.response.members),
    config: (node: String) => api.get(`/node/${node}/config`).then((response) => response.data.response)
}

export const clusterApi = {
    list: () => api.get<{ response: Cluster[] }>(`/cluster/list`).then((response) => response.data.response),
    create: (cluster: Cluster) => api.post(`/cluster/create`, cluster).then((response) => response.data.response)
}
