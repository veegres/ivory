import axios, {AxiosProgressEvent, AxiosRequestConfig} from "axios";
import {Instance, InstanceMap, InstanceRequest} from "../type/instance";
import {getDomain} from "./utils";
import {Cert, CertAddRequest, CertMap, CertType, CertUploadRequest} from "../type/cert";
import {Password, PasswordMap, PasswordType} from "../type/password";
import {SecretSetRequest, SecretStatus, SecretUpdateRequest} from "../type/secret";
import {
    Query,
    QueryChart,
    QueryChartRequest,
    QueryDatabasesRequest,
    QueryFields,
    QueryKillRequest,
    QueryRequest,
    QueryRunRequest,
    QuerySchemasRequest,
    QueryTablesRequest,
    QueryType
} from "../type/query";
import {AppConfig, AppInfo, Login, Response} from "../type/common";
import {Bloat, BloatRequest} from "../type/bloat";
import {Cluster, ClusterAuto, ClusterMap} from "../type/cluster";

export const api = axios.create({baseURL: '/api'})

// TODO we should simplify usage of react-query hooks
//      Possible solutions:
//      - provide queryKey with request in same object
//      - try to make args that can be used than directly in react-query hooks
//      - can be problems with same route as for DELETE / GET / POST


export const GeneralApi = {
    info: () => api
        .get<Response<AppInfo>>(`/info`)
        .then((response) => response.data.response),
    login: (req: Login) => api
        .post<Response<any>>(`/login`, req)
        .then((response) => response.data.response),
    setConfig: (request: AppConfig) => api
        .post<Response<string>>(`/initial/config`, request)
        .then((response) => response.data.response),
}

export const InitialApi = {
    setSecret: (request: SecretSetRequest) => api
        .post<Response<string>>(`/initial/secret`, request)
        .then((response) => response.data.response),
    erase: () => api
        .delete<Response<string>>(`/initial/erase`)
        .then((response) => response.data.response),
}

export const SafeApi = {
    getSecret: () => api
        .get<Response<SecretStatus>>(`/safe/secret`)
        .then((response) => response.data.response),
    changeSecret: (request: SecretUpdateRequest) => api
        .post<Response<string>>(`/safe/secret`, request)
        .then((response) => response.data.response),
    erase: () => api
        .delete<Response<string>>(`/safe/erase`)
        .then((response) => response.data.response),
}

export const InstanceApi = {
    overview: (request: InstanceRequest) => api
        .get<Response<Instance[]>>(`/instance/overview`, {params: {request: JSON.stringify(request)}})
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
        .get<Response<any>>(`/instance/config`, {params: {request: JSON.stringify(request)}})
        .then((response) => response.data.response),
    updateConfig: (request: InstanceRequest) => api
        .patch<Response<any>>(`/instance/config`, request)
        .then((response) => response.data.response),
    switchover: (request: InstanceRequest) => api
        .post<Response<string>>(`/instance/switchover`, request)
        .then((response) => response.data.response),
    deleteSwitchover: (request: InstanceRequest) => api
        .delete<Response<string>>(`/instance/switchover`, {params: {request: JSON.stringify(request)}})
        .then((response) => response.data.response),
    reinitialize: (request: InstanceRequest) => api
        .post<Response<string>>(`/instance/reinitialize`, request)
        .then((response) => response.data.response),
    restart: (request: InstanceRequest) => api
        .post<Response<string>>(`/instance/restart`, request)
        .then((response) => response.data.response),
    deleteRestart: (request: InstanceRequest) => api
        .delete<Response<string>>(`/instance/restart`, {params: {request: JSON.stringify(request)}})
        .then((response) => response.data.response),
    reload: (request: InstanceRequest) => api
        .post<Response<string>>(`/instance/reload`, request)
        .then((response) => response.data.response),
    failover: (request: InstanceRequest) => api
        .post<Response<string>>(`/instance/failover`, request)
        .then((response) => response.data.response),
    activate: (request: InstanceRequest) => api
        .post<Response<string>>(`/instance/activate`, request)
        .then((response) => response.data.response),
    pause: (request: InstanceRequest) => api
        .post<Response<string>>(`/instance/pause`, request)
        .then((response) => response.data.response),
}

export const ClusterApi = {
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
        .then((response) => response.data.response),
}

export const TagApi = {
    list: () => api
        .get<Response<string[]>>(`/tag`)
        .then((response) => response.data.response),
}

export const BloatApi = {
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

    stream: (uuid: string) => new EventSource(`/api/cli/bloat/job/${uuid}/stream`),
}

export const QueryApi = {
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
        .delete<Response<string>>(`/query/${uuid}`)
        .then((response) => response.data.response),
    run: (req: QueryRunRequest) => api
        .post<Response<QueryFields>>(`/query/run`, req)
        .then((response) => response.data.response),
    databases: (req: QueryDatabasesRequest) => api
        .post<Response<string[]>>(`/query/databases`, req)
        .then((response) => response.data.response),
    schemas: (req: QuerySchemasRequest) => api
        .post<Response<string[]>>(`/query/schemas`, req)
        .then((response) => response.data.response),
    tables: (req: QueryTablesRequest) => api
        .post<Response<string[]>>(`/query/tables`, req)
        .then((response) => response.data.response),
    chartCommon: (req: QueryChartRequest) => api
        .post<Response<QueryChart[]>>(`/query/chart/common`, req)
        .then((response) => response.data.response),
    chartDatabase: (req: QueryChartRequest) => api
        .post<Response<QueryChart[]>>(`/query/chart/database`, req)
        .then((response) => response.data.response),
    cancel: (req: QueryKillRequest) => api
        .post<Response<string>>(`/query/cancel`, req)
        .then((response) => response.data.response),
    terminate: (req: QueryKillRequest) => api
        .post<Response<string>>(`/query/terminate`, req)
        .then((response) => response.data.response),
    getLog: (uuid: string) => api
        .get<Response<QueryFields[]>>(`/query/log/${uuid}`)
        .then((response) => response.data.response),
    deleteLog: (uuid: string) => api
        .delete<Response<string>>(`/query/log/${uuid}`)
        .then((response) => response.data.response),
}

export const PasswordApi = {
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
        .then((response) => response.data.response),
}

export const CertApi = {
    list: (type?: CertType) => api
        .get<Response<CertMap>>(`/cert`, {params: {type}})
        .then((response => response.data.response)),
    upload: async (request: CertUploadRequest) => {
        const {file, type, setProgress} = request
        const formData = new FormData()
        formData.append("type", type.toString())
        formData.append("file", file)
        const config: AxiosRequestConfig = {
            headers: {"Content-Type": "multipart/form-data"},
            onUploadProgress: (progressEvent: AxiosProgressEvent) => {
                setProgress && setProgress(progressEvent)
            }
        }
        const response = await api
            .post<Response<Cert>>(`/cert/upload`, formData, config);
        return response.data.response;
    },
    add: (request: CertAddRequest) => api
        .post<Response<Cert>>(`/cert/add`, request)
        .then((response) => response.data.response),
    delete: (uuid: string) => api
        .delete(`/cert/${uuid}`)
        .then((response) => response.data.response),
}
