import {api} from "../../api"
import {R} from "../../management/api/type"
import {
    KeeperOneRequest,
    KeeperOneResponse,
    PlatformActionRequest,
    PlatformLogsRequest,
    PlatformMetricsResponse,
    PlatformUpRequest,
    PlatformVaultConnection,
} from "./type"


export const NodeApi = {
    overview: {
        key: (h: string, p?: number) => ["node", "db", "overview", h, p],
        fn: (request: KeeperOneRequest) => api.get<R<KeeperOneResponse>>("/node/db/overview", {params: {request: JSON.stringify(request)}})
            .then((response) => response.data.response),
    },
    config: {
        key: (h: string, p?: number) => ["node", "db", "config", h, p],
        fn: (request: KeeperOneRequest) => api.get<R<KeeperOneResponse>>("/node/db/config", {params: {request: JSON.stringify(request)}})
            .then((response) => response.data.response),
    },
    updateConfig: {
        key: () => ["node", "db", "config", "update"],
        fn: (request: KeeperOneRequest) => api.patch<R<KeeperOneResponse>>("/node/db/config", request)
            .then((response) => response.data.response),
    },
    switchover: {
        key: () => ["node", "db", "switchover"],
        fn: (request: KeeperOneRequest) => api.post<R<string>>("/node/db/switchover", request)
            .then((response) => response.data.response),
    },
    deleteSwitchover: {
        key: () => ["node", "db", "switchover", "delete"],
        fn: (request: KeeperOneRequest) => api.delete<R<string>>("/node/db/switchover", {params: {request: JSON.stringify(request)}})
            .then((response) => response.data.response),
    },
    reinitialize: {
        key: () => ["node", "db", "reinitialize"],
        fn: (request: KeeperOneRequest) => api.post<R<string>>("/node/db/reinitialize", request)
            .then((response) => response.data.response),
    },
    restart: {
        key: () => ["node", "db", "restart"],
        fn: (request: KeeperOneRequest) => api.post<R<string>>("/node/db/restart", request)
            .then((response) => response.data.response),
    },
    deleteRestart: {
        key: () => ["node", "db", "restart", "delete"],
        fn: (request: KeeperOneRequest) => api.delete<R<string>>("/node/db/restart", {params: {request: JSON.stringify(request)}})
            .then((response) => response.data.response),
    },
    reload: {
        key: () => ["node", "db", "reload"],
        fn: (request: KeeperOneRequest) => api.post<R<string>>("/node/db/reload", request)
            .then((response) => response.data.response),
    },
    failover: {
        key: () => ["node", "db", "failover"],
        fn: (request: KeeperOneRequest) => api.post<R<string>>("/node/db/failover", request)
            .then((response) => response.data.response),
    },
    activate: {
        key: () => ["node", "db", "activate"],
        fn: (request: KeeperOneRequest) => api.post<R<string>>("/node/db/activate", request)
            .then((response) => response.data.response),
    },
    pause: {
        key: () => ["node", "db", "pause"],
        fn: (request: KeeperOneRequest) => api.post<R<string>>("/node/db/pause", request)
            .then((response) => response.data.response),
    },
    metrics: {
        key: (host: string) => ["node", "platform", "metrics", host],
        fn: (request: PlatformVaultConnection) => api.get<R<PlatformMetricsResponse>>("/node/platform/metrics", {params: {request: JSON.stringify(request)}})
            .then((response) => response.data.response),
    },
    logs: {
        url: (request: PlatformLogsRequest) => `/api/node/platform/logs?request=${encodeURIComponent(JSON.stringify(request))}`,
    },
    deployment: {
        list: {
            key: (request: PlatformVaultConnection) => ["node", "platform", "container", "list", request.host, request.port, request.vaultId],
            fn: (request: PlatformVaultConnection) => api.get<R<string[]>>("/node/platform/container", {params: {request: JSON.stringify(request)}})
                .then((response) => response.data.response),
        },
        logs: {
            url: (request: PlatformLogsRequest) => `/api/node/platform/container/logs?request=${encodeURIComponent(JSON.stringify(request))}`,
        },
        up: {
            key: () => ["node", "platform", "container", "up"],
            fn: (request: PlatformUpRequest) => api.post<R<string[]>>("/node/platform/container/up", request)
                .then((response) => response.data.response),
        },
        start: {
            key: () =>  ["node", "platform", "container", "start"],
            fn: (request: PlatformActionRequest) => api.post<R<string[]>>("/node/platform/container/start", request)
                .then((response) => response.data.response),
        },
        stop: {
            key: () =>  ["node", "platform", "container", "stop"],
            fn: (request: PlatformActionRequest) => api.post<R<string[]>>("/node/platform/container/stop", request)
                .then((response) => response.data.response),
        },
        restart: {
            key: () =>  ["node", "platform", "container", "restart"],
            fn: (request: PlatformActionRequest) => api.post<R<string[]>>("/node/platform/container/restart", request)
                .then((response) => response.data.response),
        },
        down: {
            key: () =>  ["node", "platform", "container", "down"],
            fn: (request: PlatformActionRequest) => api.post<R<string[]>>("/node/platform/container/down", request)
                .then((response) => response.data.response),
        }
    }
}
