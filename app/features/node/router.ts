import {api} from "../api"
import {R} from "../management/type"
import {KeeperRequest, MetricsResponse,PlatformConnection,PlatformDeployRequest, PlatformLogsRequest} from "./type"


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
        key: (host: string) => ["node", "platform", "metrics", host],
        fn: (request: PlatformConnection) => api.get<R<MetricsResponse>>("/node/platform/metrics", {params: {request: JSON.stringify(request)}})
            .then((response) => response.data.response),
    },
    deployment: {
        list: {
            key: (host: string) => ["node", "platform", "deployment", "list", host],
            fn: (request: PlatformConnection) => api.get<R<any>>("/node/platform/deployment", {params: {request: JSON.stringify(request)}})
                .then((response) => response.data.response),
        },
        logs: {
            key: (host: string, name: string) => ["node", "platform", "deployment", "logs", host, name],
            fn: (request: PlatformLogsRequest) => api.get<R<any>>("/node/platform/deployment/logs", {params: {request: JSON.stringify(request)}})
                .then((response) => response.data.response),
        },
        deploy: {
            key: () => ["node", "platform", "deployment", "deploy"],
            fn: (request: PlatformDeployRequest) => api.post<R<any>>("/node/platform/deployment/deploy", request)
                .then((response) => response.data.response),
        },
        stop: {
            key: () => ["node", "platform", "deployment", "stop"],
            fn: (request: PlatformDeployRequest) => api.post<R<any>>("/node/platform/deployment/stop", request)
                .then((response) => response.data.response),
        },
        delete: {
            key: () => ["node", "platform", "deployment", "delete"],
            fn: (request: PlatformDeployRequest) => api.post<R<any>>("/node/platform/deployment/delete", request)
                .then((response) => response.data.response),
        }
    }
}
