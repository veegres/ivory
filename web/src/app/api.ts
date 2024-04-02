import axios, {AxiosProgressEvent, AxiosRequestConfig} from "axios";
import {Instance, InstanceMap, InstanceRequest} from "../type/instance";
import {getDomain} from "./utils";
import {Cert, CertAddRequest, CertMap, CertType, CertUploadRequest} from "../type/cert";
import {Password, PasswordMap, PasswordType} from "../type/password";
import {SecretSetRequest, SecretUpdateRequest} from "../type/secret";
import {
    Query,
    QueryChart,
    QueryChartRequest, QueryConnection,
    QueryDatabasesRequest,
    QueryFields,
    QueryKillRequest,
    QueryRequest,
    QueryRunRequest,
    QuerySchemasRequest,
    QueryTablesRequest,
    QueryType
} from "../type/query";
import {AppConfig, AppInfo, Login, Response, Sidecar} from "../type/general";
import {Bloat, BloatRequest} from "../type/bloat";
import {Cluster, ClusterAuto, ClusterMap} from "../type/cluster";

export const api = axios.create({baseURL: '/api'})

export const GeneralApi = {
    info: {
        key: () => ["info"],
        fn: () => api.get<Response<AppInfo>>(`/info`)
            .then((response) => response.data.response),
    },
    login: {
        key: () => ["login"],
        fn: (req: Login) => api.post<Response<any>>(`/login`, req)
            .then((response) => response.data.response),
    },
    setConfig: {
        key: () => ["config", "set"],
        fn: (request: AppConfig) => api.post<Response<string>>(`/initial/config`, request)
            .then((response) => response.data.response),
    }
}

export const InitialApi = {
    setSecret: {
        key: () => ["secret", "set"],
        fn: (request: SecretSetRequest) => api.post<Response<string>>(`/initial/secret`, request)
            .then((response) => response.data.response),
    },
    erase: {
        key: () => ["erase"],
        fn: () => api.delete<Response<string>>(`/initial/erase`)
            .then((response) => response.data.response),
    },
}

export const SafeApi = {
    changeSecret: {
        key: () => ["secret", "change"],
        fn: (request: SecretUpdateRequest) => api.post<Response<string>>(`/safe/secret`, request)
            .then((response) => response.data.response),
    },
    erase: {
        key: () => ["erase"],
        fn: () => api.delete<Response<string>>(`/safe/erase`)
            .then((response) => response.data.response),
    },
}

export const InstanceApi = {
    overview: {
        key: (cluster: string) => ["instance", "overview", cluster],
        fn: (request: InstanceRequest) => api.get<Response<Instance[]>>(`/instance/overview`, {params: {request: JSON.stringify(request)}})
            .then<InstanceMap>((response) => response.data.response.reduce(
                (map, instance) => {
                    const leader = instance.role === "leader"
                    const domain = getDomain(instance.sidecar)
                    map[domain] = {...instance, leader, inCluster: true, inInstances: false}
                    return map
                },
                {} as InstanceMap
            )),
    },
    config: {
        key: (sidecar: Sidecar) => ["instance", "config", sidecar.host, sidecar.port],
        fn: (request: InstanceRequest) => api.get<Response<any>>(`/instance/config`, {params: {request: JSON.stringify(request)}})
            .then((response) => response.data.response),
    },
    updateConfig: {
        key: () => ["instance", "config", "update"],
        fn: (request: InstanceRequest) => api.patch<Response<any>>(`/instance/config`, request)
            .then((response) => response.data.response),
    },
    switchover: {
        key: () => ["instance", "switchover"],
        fn: (request: InstanceRequest) => api.post<Response<string>>(`/instance/switchover`, request)
            .then((response) => response.data.response),
    },
    deleteSwitchover: {
        key: () => ["instance", "switchover", "delete"],
        fn: (request: InstanceRequest) => api.delete<Response<string>>(`/instance/switchover`, {params: {request: JSON.stringify(request)}})
            .then((response) => response.data.response),
    },
    reinitialize: {
        key: () => ["instance", "reinitialize"],
        fn: (request: InstanceRequest) => api.post<Response<string>>(`/instance/reinitialize`, request)
            .then((response) => response.data.response),
    },
    restart: {
        key: () => ["instance", "restart"],
        fn: (request: InstanceRequest) => api.post<Response<string>>(`/instance/restart`, request)
            .then((response) => response.data.response),
    },
    deleteRestart: {
        key: () => ["instance", "restart", "delete"],
        fn: (request: InstanceRequest) => api.delete<Response<string>>(`/instance/restart`, {params: {request: JSON.stringify(request)}})
            .then((response) => response.data.response),
    },
    reload: {
        key: () => ["instance", "reload"],
        fn: (request: InstanceRequest) => api.post<Response<string>>(`/instance/reload`, request)
            .then((response) => response.data.response),
    },
    failover: {
        key: () => ["instance", "failover"],
        fn: (request: InstanceRequest) => api.post<Response<string>>(`/instance/failover`, request)
            .then((response) => response.data.response),
    },
    activate: {
        key: () => ["instance", "activate"],
        fn: (request: InstanceRequest) => api.post<Response<string>>(`/instance/activate`, request)
            .then((response) => response.data.response),
    },
    pause: {
        key: () => ["instance", "pause"],
        fn: (request: InstanceRequest) => api.post<Response<string>>(`/instance/pause`, request)
            .then((response) => response.data.response),
    },
}

export const ClusterApi = {
    list: {
        key: () => ["cluster", "list"],
        fn: (tags?: string[]) => api.get<Response<Cluster[]>>(`/cluster`, {params: {tags}})
            .then((response) => response.data.response.reduce(
                (map, cluster) => {
                    map[cluster.name] = cluster
                    return map
                },
                {} as ClusterMap
            )),
    },
    update: {
        key: () => ["cluster", "update"],
        fn: (cluster: Cluster) => api.put<Response<Cluster>>(`/cluster`, cluster)
            .then((response) => response.data.response),
    },
    createAuto: {
        key: () => ["cluster", "auto", "creation"],
        fn: (cluster: ClusterAuto) => api.post<Response<Cluster>>(`/cluster/auto`, cluster)
            .then((response) => response.data.response),
    },
    delete: {
        key: () => ["cluster", "delete"],
        fn: (name: string) => api.delete(`/cluster/${name}`)
            .then((response) => response.data.response),
    },
}

export const TagApi = {
    list: {
        key: () => ["tag", "list"],
        fn: () => api.get<Response<string[]>>(`/tag`)
            .then((response) => response.data.response),
    },
}

export const BloatApi = {
    list: {
        key: (cluster: string) => ["instance", "bloat", "list", cluster],
        fn: (cluster: string) => api.get<Response<Bloat[]>>(`/cli/bloat/cluster/${cluster}`)
            .then((response) => response.data.response),
    },
    logs: {
        key: (uuid: string) => ["instance", "bloat", "logs", uuid],
        fn: (uuid: string) => api.get<string>(`/cli/bloat/${uuid}/logs`, {responseType: "text"})
            .then(({data}) => data === "" ? [] : data.split("\n")),
    },

    start: {
        key: () => ["instance", "bloat", "job", "start"],
        fn: (ctr: BloatRequest) => api.post<Response<Bloat>>(`/cli/bloat/job/start`, ctr)
            .then((response) => response.data.response),
    },
    stop: {
        key: () => ["instance", "bloat", "job", "stop"],
        fn: (uuid: string) => api.post<Response<Bloat>>(`/cli/bloat/job/${uuid}/stop`)
            .then((response) => response.data.response),
    },
    delete: {
        key: () => ["instance", "bloat", "job", "delete"],
        fn: (uuid: string) => api.delete(`/cli/bloat/job/${uuid}/delete`)
            .then((response) => response.data.response),
    },

    stream: {
        key: () => ["instance", "bloat", "job", "stream"],
        fn: (uuid: string) => new EventSource(`/api/cli/bloat/job/${uuid}/stream`)
    },
}

export const QueryApi = {
    list: {
        key: (type?: QueryType) => ["query", "list", type],
        fn: (type?: QueryType) => api.get<Response<Query[]>>(`/query`, {params: {type}})
            .then((response) => response.data.response),
    },
    update: {
        key: () => ["query", "update"],
        fn: ({id, query}: { id: string, query: QueryRequest }) => api.put<Response<Query>>(`/query/${id}`, query)
            .then((response) => response.data.response),
    },
    create: {
        key: () => ["query", "create"],
        fn: (query: QueryRequest) => api.post<Response<Query>>(`/query`, query)
            .then((response) => response.data.response),
    },
    delete: {
        key: () => ["query", "delete"],
        fn: (uuid: string) => api.delete<Response<string>>(`/query/${uuid}`)
            .then((response) => response.data.response),
    },
    run: {
        key: (uuid?: string) => ["query", "run", uuid ?? "standalone"],
        fn: (req: QueryRunRequest) => api.post<Response<QueryFields>>(`/query/run`, req)
            .then((response) => response.data.response),
    },
    databases: {
        key: (con: QueryConnection) => ["query", "databases", con],
        fn: (req: QueryDatabasesRequest) => api.post<Response<string[]>>(`/query/databases`, req)
            .then((response) => response.data.response),
    },
    schemas: {
        key: (con: QueryConnection) => ["query", "schemas", con],
        fn: (req: QuerySchemasRequest) => api.post<Response<string[]>>(`/query/schemas`, req)
            .then((response) => response.data.response),
    },
    tables: {
        key: (con: QueryConnection, schema?: string) => ["query", "tables", con, schema ?? "-"],
        fn: (req: QueryTablesRequest) => api.post<Response<string[]>>(`/query/tables`, req)
            .then((response) => response.data.response),
    },
    chart: {
        key: (req: QueryChartRequest) => ["query", "chart", req.type, req.connection],
        fn: (req: QueryChartRequest) => api.post<Response<QueryChart>>(`/query/chart`, req)
            .then((response) => response.data.response),
    },
    cancel: {
        key: () => ["query", "cancel"],
        fn: (req: QueryKillRequest) => api.post<Response<string>>(`/query/cancel`, req)
            .then((response) => response.data.response),
    },
    terminate: {
        key: () => ["query", "terminate"],
        fn: (req: QueryKillRequest) => api.post<Response<string>>(`/query/terminate`, req)
            .then((response) => response.data.response),
    },

    getLog: {
        key: (uuid: string) => ["query", "log", uuid],
        fn: (uuid: string) => api.get<Response<QueryFields[]>>(`/query/log/${uuid}`)
            .then((response) => response.data.response),
    },
    deleteLog: {
        key: () => ["query", "log", "delete"],
        fn: (uuid: string) => api.delete<Response<string>>(`/query/log/${uuid}`)
            .then((response) => response.data.response),
    },
}

export const PasswordApi = {
    list: {
        key: (type?: PasswordType) => ["password", type],
        fn: (type?: PasswordType) => api.get<Response<PasswordMap>>(`/password`, {params: {type}})
            .then((response) => response.data.response),
    },
    create: {
        key: () => ["password", "create"],
        fn: (credential: Password) => api.post<Response<{ key: string, credential: Password }>>(`/password`, credential)
            .then((response) => response.data.response),
    },
    update: {
        key: () => ["password", "update"],
        fn: ({uuid, credential}: { uuid: string, credential: Password }) => api.patch<Response<Password>>(`/password/${uuid}`, credential)
            .then((response) => response.data.response),
    },
    delete: {
        key: () => ["password", "delete"],
        fn: (uuid: string) => api.delete(`/password/${uuid}`)
            .then((response) => response.data.response),
    },
}

export const CertApi = {
    list: {
        key: (type?: CertType) => ["cert", "list", type],
        fn: (type?: CertType) => api.get<Response<CertMap>>(`/cert`, {params: {type}})
            .then((response => response.data.response)),
    },
    upload: {
        key: () => ["cert", "upload"],
        fn: async (request: CertUploadRequest) => {
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
            const response = await api.post<Response<Cert>>(`/cert/upload`, formData, config)
            return response.data.response
        },
    },
    add: {
        key: () => ["cert", "add"],
        fn: (request: CertAddRequest) => api.post<Response<Cert>>(`/cert/add`, request)
            .then((response) => response.data.response),
    },
    delete: {
        key: () => ["cert", "delete"],
        fn: (uuid: string) => api.delete(`/cert/${uuid}`)
            .then((response) => response.data.response),
    },
}
