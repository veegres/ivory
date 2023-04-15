import axios from "axios";
import {InstanceResponse, InstanceInfo, InstanceMap, InstanceRequest} from "../type/Instance";
import {getDomain} from "./utils";
import {Cert, CertAddRequest, CertMap, CertType, CertUploadRequest} from "../type/cert";
import {Password, PasswordMap, PasswordType} from "../type/password";
import {SecretSetRequest, SecretStatus, SecretUpdateRequest} from "../type/secret";
import {
    Query, QueryChart,
    QueryChartRequest,
    QueryKillRequest,
    QueryRequest,
    QueryRunRequest,
    QueryRunResponse, QuerySearchRequest,
    QueryType
} from "../type/query";
import {AppInfo, Response} from "../type/common";
import {Bloat, BloatRequest} from "../type/bloat";
import {Cluster, ClusterAuto, ClusterMap} from "../type/cluster";

const api = axios.create({baseURL: '/api'})

// TODO we should simplify usage of react-query hooks
// Possible solutions:
// - provide queryKey with request in same object
// - try to make args that can be used than directly in react-query hooks
// - can be problems with same route as for DELETE / GET / POST

export const instanceApi = {
    info: (request: InstanceRequest) => api
        .get<Response<InstanceInfo>>(`/instance/info`, {params: request})
        .then((response) => response.data.response),
    overview: (request: InstanceRequest) => api
        .get<Response<InstanceResponse[]>>(`/instance/overview`, {params: request})
        .then<InstanceMap>((response) => response.data.response.reduce(
            (map, instance) => {
                const leader = instance.role === "leader"
                const domain = getDomain(instance.sidecar)
                map[domain] = {...instance, leader, inCluster: true, inInstances: false}
                return map
            },
            {} as InstanceMap
        )),
    config: (request: InstanceRequest) => api
        .get(`/instance/config`, {params: request})
        .then((response) => response.data.response),
    updateConfig: (request: InstanceRequest) => api
        .patch(`/instance/config`, request)
        .then((response) => response.data.response),
    switchover: (request: InstanceRequest) => api
        .post(`/instance/switchover`, request)
        .then((response) => response.data.response),
    reinitialize: (request: InstanceRequest) => api
        .post(`/instance/reinitialize`, request)
        .then((response) => response.data.response)
}

export const clusterApi = {
    get: (name: string) => api
        .get<Response<Cluster>>(`/cluster/${name}`)
        .then((response) => response.data.response),
    list: (tags?: string[]) => api
        .get<Response<Cluster[]>>(`/cluster`, {params: {tags}})
        .then((response) => response.data.response.reduce(
            (map, cluster) => {
                map[cluster.name] = cluster
                return map
            },
            {} as ClusterMap
        )),
    update: (cluster: Cluster) => api
        .put<Response<Cluster>>(`/cluster`, cluster)
        .then((response) => response.data.response),
    createAuto: (cluster: ClusterAuto) => api
        .post<Response<Cluster>>(`/cluster/auto`, cluster)
        .then((response) => response.data.response),
    delete: (name: string) => api
        .delete(`/cluster/${name}`)
        .then((response) => response.data.response)
}

export const tagApi = {
    list: () => api
        .get<Response<string[]>>(`/tag`)
        .then((response) => response.data.response),
}

export const bloatApi = {
    list: (name: string) => api
        .get<Response<Bloat[]>>(`/cli/bloat/cluster/${name}`)
        .then((response) => response.data.response),
    logs: (uuid: string) => api
        .get<string>(`/cli/bloat/${uuid}/logs`, {responseType: "text"})
        .then(({data}) => data === "" ? [] : data.split("\n")),

    start: (ctr: BloatRequest) => api
        .post<Response<Bloat>>(`/cli/bloat/job/start`, ctr)
        .then((response) => response.data.response),
    stop: (uuid: string) => api
        .post<Response<Bloat>>(`/cli/bloat/job/${uuid}/stop`)
        .then((response) => response.data.response),
    delete: (uuid: string) => api
        .delete(`/cli/bloat/job/${uuid}/delete`)
        .then((response) => response.data.response),

    stream: (uuid: string) => new EventSource(`/api/cli/bloat/job/${uuid}/stream`)
}

export const infoApi = {
    get: () => api
        .get<Response<AppInfo>>(`/info`)
        .then((response) => response.data.response),
    erase: () => api
        .post<Response<string>>(`/erase`)
        .then((response) => response.data.response)
}

export const queryApi = {
    list: (type?: QueryType) => api
        .get<Response<Query[]>>(`/query`, {params: {type}})
        .then((response) => response.data.response),
    update: ({id, query}: { id: string, query: QueryRequest }) => api
        .put<Response<Query>>(`/query/${id}`, query)
        .then((response) => response.data.response),
    create: (query: QueryRequest) => api
        .post<Response<Query>>(`/query`, query)
        .then((response) => response.data.response),
    delete: (uuid: string) => api
        .delete(`/query/${uuid}`)
        .then((response) => response.data.response),
    run: (req: QueryRunRequest) => api
        .post<Response<QueryRunResponse>>(`/query/run`, req)
        .then((response) => response.data.response),
    databases: (req: QuerySearchRequest) => api
        .post<Response<string[]>>(`/query/databases`, req)
        .then((response) => response.data.response),
    schemas: (req: QuerySearchRequest) => api
        .post<Response<string[]>>(`/query/schemas`, req)
        .then((response) => response.data.response),
    tables: (req: QuerySearchRequest) => api
        .post<Response<string[]>>(`/query/tables`, req)
        .then((response) => response.data.response),
    chartCommon: (req: QueryChartRequest) => api
        .post<Response<QueryChart[]>>(`/query/chart/common`, req)
        .then((response) => response.data.response),
    chartDatabase: (req: QueryChartRequest) => api
        .post<Response<QueryChart[]>>(`/query/chart/database`, req)
        .then((response) => response.data.response),
    cancel: (req: QueryKillRequest) => api
        .post<Response<QueryRunResponse>>(`/query/cancel`, req)
        .then((response) => response.data.response),
    terminate: (req: QueryKillRequest) => api
        .post<Response<QueryRunResponse>>(`/query/terminate`, req)
        .then((response) => response.data.response),
}

export const secretApi = {
    get: () => api
        .get<Response<SecretStatus>>(`/secret`)
        .then((response) => response.data.response),
    set: (request: SecretSetRequest) => api
        .post<Response<string>>(`/secret/set`, request)
        .then((response) => response.data.response),
    update: (request: SecretUpdateRequest) => api
        .post<Response<string>>(`/secret/update`, request)
        .then((response) => response.data.response),
    clean: () => api
        .post<Response<string>>(`/secret/clean`)
        .then((response) => response.data.response),
}

export const passwordApi = {
    list: (type?: PasswordType) => api
        .get<Response<PasswordMap>>(`/password`, {params: {type}})
        .then((response) => response.data.response),
    create: (credential: Password) => api
        .post<Response<{ key: string, credential: Password }>>(`/password`, credential)
        .then((response) => response.data.response),
    update: ({uuid, credential}: { uuid: string, credential: Password }) => api
        .patch<Response<Password>>(`/password/${uuid}`, credential)
        .then((response) => response.data.response),
    delete: (uuid: string) => api
        .delete(`/password/${uuid}`)
        .then((response) => response.data.response)
}

export const certApi = {
    list: (type?: CertType) => api
        .get<Response<CertMap>>(`/cert`, {params: {type}})
        .then((response => response.data.response)),
    upload: (request: CertUploadRequest) => {
        const {file, type, setProgress} = request
        const formData = new FormData()
        formData.append("type", type.toString())
        formData.append("file", file)
        const config = {
            headers: {"Content-Type": "multipart/form-data"},
            onUploadProgress: (progressEvent: ProgressEvent) => {
                setProgress && setProgress(progressEvent)
            }
        }
        return api
            .post<Response<Cert>>(`/cert/upload`, formData, config)
            .then((response) => response.data.response)
    },
    add: (request: CertAddRequest) => api
        .post<Response<Cert>>(`/cert/add`, request)
        .then((response) => response.data.response),
    delete: (uuid: string) => api
        .delete(`/cert/${uuid}`)
        .then((response) => response.data.response)
}
