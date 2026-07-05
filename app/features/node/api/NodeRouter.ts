import {api} from "../../Api"
import {R} from "../../management/api/ManagementType"
import {
    KeeperOneRequest,
    KeeperOneResponse,
    KeeperPlugin,
    PlatformActionRequest,
    PlatformDeployOptions,
    PlatformDeployOptionsRequest,
    PlatformLogsRequest,
    PlatformMetricsResponse,
    PlatformProcessesResponse,
    PlatformUpRequest,
    PlatformVaultConnection,
} from "./NodeType"


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
    processes: {
        key: (host: string) => ["node", "platform", "processes", host],
        fn: (request: PlatformVaultConnection) => api.get<R<PlatformProcessesResponse>>("/node/platform/processes", {params: {request: JSON.stringify(request)}})
            .then((response) => response.data.response),
    },
    deployment: {
        deployOptions: {
            key: (plugin: KeeperPlugin) => ["node", "platform", "container", "deploy-options", plugin],
            fn: (request: PlatformDeployOptionsRequest) => api.get<R<PlatformDeployOptions>>("/node/platform/container/deploy-options", {params: {request: JSON.stringify(request)}})
                .then((response) => response.data.response),
        },
        list: {
            key: (request: PlatformVaultConnection) => ["node", "platform", "container", "list", request.host, request.port, request.vaultId],
            fn: (request: PlatformVaultConnection) => api.get<R<string[]>>("/node/platform/container", {params: {request: JSON.stringify(request)}})
                .then((response) => response.data.response),
        },
        logs: {
            url: (request: PlatformLogsRequest) => `/api/node/platform/container/logs?request=${encodeURIComponent(JSON.stringify(request))}`,
        },
        metrics: {
            key: (request: PlatformActionRequest) => ["node", "platform", "container", "metrics", request.connection.host, request.name],
            fn: (request: PlatformActionRequest) => api.get<R<PlatformMetricsResponse>>("/node/platform/container/metrics", {params: {request: JSON.stringify(request)}})
                .then((response) => response.data.response),
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
