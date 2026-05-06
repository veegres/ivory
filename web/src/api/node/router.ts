import {api} from "../api"
import {R} from "../management/type"
import {DockerLogsRequest, DockerRequest, KeeperRequest, MetricsResponse, SshConnection} from "./type"


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
        key: (host: string) => ["node", "ssh", "metrics", host],
        fn: (request: SshConnection) => api.get<R<MetricsResponse>>("/node/ssh/metrics", {params: {request: JSON.stringify(request)}})
            .then((response) => response.data.response),
    },
    docker: {
        list: {
            key: (host: string) => ["node", "ssh", "docker", "list", host],
            fn: (request: SshConnection) => api.get<R<any>>("/node/ssh/docker", {params: {request: JSON.stringify(request)}})
                .then((response) => response.data.response),
        },
        logs: {
            key: (host: string, container: string) => ["node", "ssh", "docker", "logs", host, container],
            fn: (request: DockerLogsRequest) => api.get<R<any>>("/node/ssh/docker/logs", {params: {request: JSON.stringify(request)}})
                .then((response) => response.data.response),
        },
        deploy: {
            key: () => ["node", "ssh", "docker", "deploy"],
            fn: (request: DockerRequest) => api.post<R<any>>("/node/ssh/docker/deploy", request)
                .then((response) => response.data.response),
        },
        stop: {
            key: () => ["node", "ssh", "docker", "stop"],
            fn: (request: DockerRequest) => api.post<R<any>>("/node/ssh/docker/stop", request)
                .then((response) => response.data.response),
        },
        run: {
            key: () => ["node", "ssh", "docker", "run"],
            fn: (request: DockerRequest) => api.post<R<any>>("/node/ssh/docker/run", request)
                .then((response) => response.data.response),
        },
        delete: {
            key: () => ["node", "ssh", "docker", "delete"],
            fn: (request: DockerRequest) => api.post<R<any>>("/node/ssh/docker/delete", request)
                .then((response) => response.data.response),
        }
    }
}
