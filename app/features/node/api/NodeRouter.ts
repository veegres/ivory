import {api} from "../../Api"
import {R} from "../../management/api/ManagementType"
import {
    KeeperDeployRequest,
    KeeperOneRequest,
    KeeperOneResponse,
    PlatformActionRequest,
    PlatformInfoResponse,
    PlatformLogsRequest,
    PlatformMetricsResponse,
    PlatformProcessesResponse,
    PlatformUpRequest,
    PlatformVaultConnection,
} from "./NodeType"

// NodeApi mirrors the server's node.Router route groups: keeper (HA
// management, /node/keeper/...), system (node-level, /node/platform/system/...),
// container (lifecycle, /node/platform/container/...) and, nested inside
// container like the server's containerKeeperGroup, container.keeper (deploy
// actions, /node/platform/container/keeper/...).
export const NodeApi = {
    keeper: {
        overview: {
            key: (h: string, p?: number) => ["node", "keeper", "overview", h, p],
            fn: (request: KeeperOneRequest) => api.get<R<KeeperOneResponse>>("/node/keeper/overview", {params: {request: JSON.stringify(request)}})
                .then((response) => response.data.response),
        },
        config: {
            key: (h: string, p?: number) => ["node", "keeper", "config", h, p],
            fn: (request: KeeperOneRequest) => api.get<R<KeeperOneResponse>>("/node/keeper/config", {params: {request: JSON.stringify(request)}})
                .then((response) => response.data.response),
        },
        updateConfig: {
            key: () => ["node", "keeper", "config", "update"],
            fn: (request: KeeperOneRequest) => api.patch<R<KeeperOneResponse>>("/node/keeper/config", request)
                .then((response) => response.data.response),
        },
        switchover: {
            key: () => ["node", "keeper", "switchover"],
            fn: (request: KeeperOneRequest) => api.post<R<string>>("/node/keeper/switchover", request)
                .then((response) => response.data.response),
        },
        deleteSwitchover: {
            key: () => ["node", "keeper", "switchover", "delete"],
            fn: (request: KeeperOneRequest) => api.delete<R<string>>("/node/keeper/switchover", {params: {request: JSON.stringify(request)}})
                .then((response) => response.data.response),
        },
        reinitialize: {
            key: () => ["node", "keeper", "reinitialize"],
            fn: (request: KeeperOneRequest) => api.post<R<string>>("/node/keeper/reinitialize", request)
                .then((response) => response.data.response),
        },
        restart: {
            key: () => ["node", "keeper", "restart"],
            fn: (request: KeeperOneRequest) => api.post<R<string>>("/node/keeper/restart", request)
                .then((response) => response.data.response),
        },
        deleteRestart: {
            key: () => ["node", "keeper", "restart", "delete"],
            fn: (request: KeeperOneRequest) => api.delete<R<string>>("/node/keeper/restart", {params: {request: JSON.stringify(request)}})
                .then((response) => response.data.response),
        },
        reload: {
            key: () => ["node", "keeper", "reload"],
            fn: (request: KeeperOneRequest) => api.post<R<string>>("/node/keeper/reload", request)
                .then((response) => response.data.response),
        },
        failover: {
            key: () => ["node", "keeper", "failover"],
            fn: (request: KeeperOneRequest) => api.post<R<string>>("/node/keeper/failover", request)
                .then((response) => response.data.response),
        },
        activate: {
            key: () => ["node", "keeper", "activate"],
            fn: (request: KeeperOneRequest) => api.post<R<string>>("/node/keeper/activate", request)
                .then((response) => response.data.response),
        },
        pause: {
            key: () => ["node", "keeper", "pause"],
            fn: (request: KeeperOneRequest) => api.post<R<string>>("/node/keeper/pause", request)
                .then((response) => response.data.response),
        },
    },
    system: {
        metrics: {
            key: (host: string) => ["node", "platform", "system", "metrics", host],
            fn: (request: PlatformVaultConnection) => api.get<R<PlatformMetricsResponse>>("/node/platform/system/metrics", {params: {request: JSON.stringify(request)}})
                .then((response) => response.data.response),
        },
        logs: {
            url: (request: PlatformLogsRequest) => `/api/node/platform/system/logs?request=${encodeURIComponent(JSON.stringify(request))}`,
        },
        processes: {
            key: (host: string) => ["node", "platform", "system", "processes", host],
            fn: (request: PlatformVaultConnection) => api.get<R<PlatformProcessesResponse>>("/node/platform/system/processes", {params: {request: JSON.stringify(request)}})
                .then((response) => response.data.response),
        },
        info: {
            key: (host: string) => ["node", "platform", "system", "info", host],
            fn: (request: PlatformVaultConnection) => api.get<R<PlatformInfoResponse>>("/node/platform/system/info", {params: {request: JSON.stringify(request)}})
                .then((response) => response.data.response),
        },
    },
    container: {
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
        },
        // node.container.keeper.deploy is the single-node counterpart of
        // cluster.deploy.
        keeper: {
            deploy: {
                key: () => ["node", "platform", "container", "keeper", "deploy"],
                fn: (request: KeeperDeployRequest) => api.post<R<string[]>>("/node/platform/container/keeper/deploy", request)
                    .then((response) => response.data.response),
            },
        },
    },
}
