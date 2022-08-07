import axios from "axios";
import {
    Cluster,
    ClusterMap,
    CompactTable,
    CompactTableRequest, Credential,
    CredentialMap,
    InstanceOverview,
    InstanceResponse,
    Response,
    SecretSetRequest,
    SecretStatus,
    SecretUpdateRequest,
    Instance
} from "./types";
import {getPatroniDomain} from "./utils";

const api = axios.create({ baseURL: '/api' })

export const nodeApi = {
    overview: (node: String) =>
        api.get<Response<InstanceOverview>>(`/node/${node}/overview`).then((response) => response.data.response),
    cluster: (node: String) =>
        api.get<Response<{ members: InstanceResponse[] }>>(`/node/${node}/cluster`)
            .then<Instance[]>((response) => response.data.response.members.map((node) => ({
                ...node,
                api_domain: getPatroniDomain(node.api_url),
                isLeader: node.role === "leader"
            }))),
    config: (node: String) =>
        api.get(`/node/${node}/config`).then((response) => response.data.response),
    updateConfig: ({node, config}: { node: string, config: string }) =>
        api.patch(`/node/${node}/config`, config).then((response) => response.data.response),
    switchover: ({node, leader, candidate}: { node: string, leader: string, candidate?: string }) =>
        api.post(`/node/${node}/switchover`, {candidate, leader}).then((response) => response.data.response),
    reinitialize: (node: string, force: boolean = true) =>
        api.post(`/node/${node}/reinitialize`, {force: force}).then((response) => response.data.response)
}

export const clusterApi = {
    get: (name: string) =>
        api.get<Response<Cluster>>(`/cluster/${name}`).then((response) => response.data.response),
    list: () =>
        api.get<Response<Cluster[]>>(`/cluster`).then((response) => response.data.response.reduce(
            (prev, current) => {
                prev[current.name] = current.nodes;
                return prev
            },
            {} as ClusterMap
        )),
    update: (cluster: Cluster) =>
        api.put<Response<Cluster>>(`/cluster`, cluster).then((response) => response.data.response),
    delete: (name: string) =>
        api.delete(`/cluster/${name}`).then((response) => response.data.response)
}

export const bloatApi = {
    list: (name: string) =>
        api.get<Response<CompactTable[]>>(`/cli/bloat/cluster/${name}`).then((response) => response.data.response),
    logs: (uuid: string) =>
        api.get<string>(`/cli/bloat/${uuid}/logs`, {responseType: "text"}).then(({data}) => data === "" ? [] : data.split("\n")),

    start: (ctr: CompactTableRequest) =>
        api.post<Response<CompactTable>>(`/cli/bloat/job/start`, ctr).then((response) => response.data.response),
    stop: (uuid: string) =>
        api.post<Response<CompactTable>>(`/cli/bloat/job/${uuid}/stop`).then((response) => response.data.response),
    delete: (uuid: string) =>
        api.delete(`/cli/bloat/job/${uuid}/delete`).then((response) => response.data.response),
    stream: (uuid: string) => new EventSource(`/api/cli/bloat/job/${uuid}/stream`)
}

export const secretApi = {
    get: () =>
        api.get<Response<SecretStatus>>(`/secret`).then((response) => response.data.response),
    set: (request: SecretSetRequest) =>
        api.post<Response<string>>(`/secret/set`, request).then((response) => response.data.response),
    update: (request: SecretUpdateRequest) =>
        api.post<Response<string>>(`/secret/update`, request).then((response) => response.data.response),
    clean: () =>
        api.post<Response<string>>(`/secret/clean`).then((response) => response.data.response),
}

export const credentialApi = {
    get: () =>
        api.get<Response<CredentialMap>>(`/credential`).then((response) => response.data.response),
    create: (credential: Credential) =>
        api.post<Response<{key: string, credential: Credential}>>(`/credential`, credential).then((response) => response.data.response),
    update: ({ uuid, credential }: { uuid: string, credential: Credential }) =>
        api.patch<Response<Credential>>(`/credential/${uuid}`, credential).then((response) => response.data.response),
    delete: (uuid: string) =>
        api.delete(`/credential/${uuid}`).then((response) => response.data.response)
}
