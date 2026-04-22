import {api} from "../api"
import {R} from "../management/type"
import {Connection,NodeRequest} from "./type"


export const NodeApi = {
    overview: {
        key: (connection: Connection) => ["node", "db", "overview", connection.host, connection.keeperPort],
        fn: (request: NodeRequest) => api.get<R<any>>("/node/db/overview", {params: {request: JSON.stringify(request)}})
            .then((response) => response.data.response),
    },
    config: {
        key: (connection: Connection) => ["node", "db", "config", connection.host, connection.keeperPort],
        fn: (request: NodeRequest) => api.get<R<any>>("/node/db/config", {params: {request: JSON.stringify(request)}})
            .then((response) => response.data.response),
    },
    updateConfig: {
        key: () => ["node", "db", "config", "update"],
        fn: (request: NodeRequest) => api.patch<R<any>>("/node/db/config", request)
            .then((response) => response.data.response),
    },
    switchover: {
        key: () => ["node", "db", "switchover"],
        fn: (request: NodeRequest) => api.post<R<string>>("/node/db/switchover", request)
            .then((response) => response.data.response),
    },
    deleteSwitchover: {
        key: () => ["node", "db", "switchover", "delete"],
        fn: (request: NodeRequest) => api.delete<R<string>>("/node/db/switchover", {params: {request: JSON.stringify(request)}})
            .then((response) => response.data.response),
    },
    reinitialize: {
        key: () => ["node", "db", "reinitialize"],
        fn: (request: NodeRequest) => api.post<R<string>>("/node/db/reinitialize", request)
            .then((response) => response.data.response),
    },
    restart: {
        key: () => ["node", "db", "restart"],
        fn: (request: NodeRequest) => api.post<R<string>>("/node/db/restart", request)
            .then((response) => response.data.response),
    },
    deleteRestart: {
        key: () => ["node", "db", "restart", "delete"],
        fn: (request: NodeRequest) => api.delete<R<string>>("/node/db/restart", {params: {request: JSON.stringify(request)}})
            .then((response) => response.data.response),
    },
    reload: {
        key: () => ["node", "db", "reload"],
        fn: (request: NodeRequest) => api.post<R<string>>("/node/db/reload", request)
            .then((response) => response.data.response),
    },
    failover: {
        key: () => ["node", "db", "failover"],
        fn: (request: NodeRequest) => api.post<R<string>>("/node/db/failover", request)
            .then((response) => response.data.response),
    },
    activate: {
        key: () => ["node", "db", "activate"],
        fn: (request: NodeRequest) => api.post<R<string>>("/node/db/activate", request)
            .then((response) => response.data.response),
    },
    pause: {
        key: () => ["node", "db", "pause"],
        fn: (request: NodeRequest) => api.post<R<string>>("/node/db/pause", request)
            .then((response) => response.data.response),
    },
    metrics: {
        key: (vmId: string) => ["node", "vm", "metrics", vmId],
        fn: (request: {vmId: string}) => api.get<R<any>>("/node/vm/metrics", {params: {request: JSON.stringify(request)}})
            .then((response) => response.data.response),
    },
    docker: {
        list: {
            key: (vmId: string) => ["node", "vm", "docker", "list", vmId],
            fn: (request: {vmId: string}) => api.get<R<any>>("/node/vm/docker", {params: {request: JSON.stringify(request)}})
                .then((response) => response.data.response),
        },
        logs: {
            key: (vmId: string, container: string) => ["node", "vm", "docker", "logs", vmId, container],
            fn: (request: any) => api.get<R<any>>("/node/vm/docker/logs", {params: {request: JSON.stringify(request)}})
                .then((response) => response.data.response),
        },
        deploy: {
            key: () => ["node", "vm", "docker", "deploy"],
            fn: (request: any) => api.post<R<any>>("/node/vm/docker/deploy", request)
                .then((response) => response.data.response),
        },
        stop: {
            key: () => ["node", "vm", "docker", "stop"],
            fn: (request: any) => api.post<R<any>>("/node/vm/docker/stop", request)
                .then((response) => response.data.response),
        },
        run: {
            key: () => ["node", "vm", "docker", "run"],
            fn: (request: any) => api.post<R<any>>("/node/vm/docker/run", request)
                .then((response) => response.data.response),
        },
        delete: {
            key: () => ["node", "vm", "docker", "delete"],
            fn: (request: any) => api.post<R<any>>("/node/vm/docker/delete", request)
                .then((response) => response.data.response),
        }
    }
}
