import axios from "axios";
import {ClusterList, NodePatroni, Node} from "./types";

const api = axios.create({ baseURL: '/api' })

export const getNodePatroni = (node: String) => api.get<{ response: NodePatroni }>(`/node/${node}/patroni`).then((response) => response.data.response)
export const getNodeCluster = (node: String) => api.get<{ response: { members: Node[] } }>(`/node/${node}/cluster`).then((response) => response.data.response.members)
export const getNodeConfig = (node: String) => api.get(`/node/${node}/config`).then((response) => response.data.response)
export const getClusterList = () => api.get<{ response: ClusterList[] }>(`/cluster/list`).then((response) => response.data.response)


