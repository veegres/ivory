import {api} from "../api"
import {R} from "../management/type"
import {CloudConnection,ContainerLogsRequest, ContainerRequest, KeeperRequest, MetricsResponse} from "./type"


export const NodeApi = {
    overview: {
        key: (h: string, p?: number) => ["node", "db", "overview", h, p],
        fn: (request: KeeperRequest) => api.get<R<any>>("/node/db/overview", {params: {request: JSON.stringify(request)}})
            .then((response) => response.data.response),
    },
    config: {
        key: (h: string, p?: number) => ["node", "db", "config", h, p],
        fn: (request: KeeperRequest) => api.get<R<any>>("/node/db/config", {params: {request: JSON.stringify(request)}})
            .then((response) => response.data.response),
    },
    updateConfig: {
        key: () => ["node", "db", "config", "update"],
        fn: (request: KeeperRequest) => api.patch<R<any>>("/node/db/config", request)
            .then((response) => response.data.response),
    },
    switchover: {
        key: () => ["node", "db", "switchover"],
        fn: (request: KeeperRequest) => api.post<R<string>>("/node/db/switchover", request)
            .then((response) => response.data.response),
    },
    deleteSwitchover: {
        key: () => ["node", "db", "switchover", "delete"],
        fn: (request: KeeperRequest) => api.delete<R<string>>("/node/db/switchover", {params: {request: JSON.stringify(request)}})
            .then((response) => response.data.response),
    },
    reinitialize: {
        key: () => ["node", "db", "reinitialize"],
        fn: (request: KeeperRequest) => api.post<R<string>>("/node/db/reinitialize", request)
            .then((response) => response.data.response),
    },
    restart: {
        key: () => ["node", "db", "restart"],
        fn: (request: KeeperRequest) => api.post<R<string>>("/node/db/restart", request)
            .then((response) => response.data.response),
    },
    deleteRestart: {
        key: () => ["node", "db", "restart", "delete"],
        fn: (request: KeeperRequest) => api.delete<R<string>>("/node/db/restart", {params: {request: JSON.stringify(request)}})
            .then((response) => response.data.response),
    },
    reload: {
        key: () => ["node", "db", "reload"],
        fn: (request: KeeperRequest) => api.post<R<string>>("/node/db/reload", request)
            .then((response) => response.data.response),
    },
    failover: {
        key: () => ["node", "db", "failover"],
        fn: (request: KeeperRequest) => api.post<R<string>>("/node/db/failover", request)
            .then((response) => response.data.response),
    },
    activate: {
        key: () => ["node", "db", "activate"],
        fn: (request: KeeperRequest) => api.post<R<string>>("/node/db/activate", request)
            .then((response) => response.data.response),
    },
    pause: {
        key: () => ["node", "db", "pause"],
        fn: (request: KeeperRequest) => api.post<R<string>>("/node/db/pause", request)
            .then((response) => response.data.response),
    },
    metrics: {
        key: (host: string) => ["node", "cloud", "metrics", host],
        fn: (request: CloudConnection) => api.get<R<MetricsResponse>>("/node/cloud/metrics", {params: {request: JSON.stringify(request)}})
            .then((response) => response.data.response),
    },
    container: {
        list: {
            key: (host: string) => ["node", "cloud", "container", "list", host],
            fn: (request: CloudConnection) => api.get<R<any>>("/node/cloud/container", {params: {request: JSON.stringify(request)}})
                .then((response) => response.data.response),
        },
        logs: {
            key: (host: string, container: string) => ["node", "cloud", "container", "logs", host, container],
            fn: (request: ContainerLogsRequest) => api.get<R<any>>("/node/cloud/container/logs", {params: {request: JSON.stringify(request)}})
                .then((response) => response.data.response),
        },
        deploy: {
            key: () => ["node", "cloud", "container", "deploy"],
            fn: (request: ContainerRequest) => api.post<R<any>>("/node/cloud/container/deploy", request)
                .then((response) => response.data.response),
        },
        stop: {
            key: () => ["node", "cloud", "container", "stop"],
            fn: (request: ContainerRequest) => api.post<R<any>>("/node/cloud/container/stop", request)
                .then((response) => response.data.response),
        },
        run: {
            key: () => ["node", "cloud", "container", "run"],
            fn: (request: ContainerRequest) => api.post<R<any>>("/node/cloud/container/run", request)
                .then((response) => response.data.response),
        },
        delete: {
            key: () => ["node", "cloud", "container", "delete"],
            fn: (request: ContainerRequest) => api.post<R<any>>("/node/cloud/container/delete", request)
                .then((response) => response.data.response),
        }
    }
}
