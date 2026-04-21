import {api} from "../api"
import {R} from "../management/type"
import {NodeRequest, Sidecar} from "./type"


export const NodeApi = {
    config: {
        key: (sidecar: Sidecar) => ["node", "config", sidecar.host, sidecar.port],
        fn: (request: NodeRequest) => api.get<R<any>>("/node/config", {params: {request: JSON.stringify(request)}})
            .then((response) => response.data.response),
    },
    updateConfig: {
        key: () => ["node", "config", "update"],
        fn: (request: NodeRequest) => api.patch<R<any>>("/node/config", request)
            .then((response) => response.data.response),
    },
    switchover: {
        key: () => ["node", "switchover"],
        fn: (request: NodeRequest) => api.post<R<string>>("/node/switchover", request)
            .then((response) => response.data.response),
    },
    deleteSwitchover: {
        key: () => ["node", "switchover", "delete"],
        fn: (request: NodeRequest) => api.delete<R<string>>("/node/switchover", {params: {request: JSON.stringify(request)}})
            .then((response) => response.data.response),
    },
    reinitialize: {
        key: () => ["node", "reinitialize"],
        fn: (request: NodeRequest) => api.post<R<string>>("/node/reinitialize", request)
            .then((response) => response.data.response),
    },
    restart: {
        key: () => ["node", "restart"],
        fn: (request: NodeRequest) => api.post<R<string>>("/node/restart", request)
            .then((response) => response.data.response),
    },
    deleteRestart: {
        key: () => ["node", "restart", "delete"],
        fn: (request: NodeRequest) => api.delete<R<string>>("/node/restart", {params: {request: JSON.stringify(request)}})
            .then((response) => response.data.response),
    },
    reload: {
        key: () => ["node", "reload"],
        fn: (request: NodeRequest) => api.post<R<string>>("/node/reload", request)
            .then((response) => response.data.response),
    },
    failover: {
        key: () => ["node", "failover"],
        fn: (request: NodeRequest) => api.post<R<string>>("/node/failover", request)
            .then((response) => response.data.response),
    },
    activate: {
        key: () => ["node", "activate"],
        fn: (request: NodeRequest) => api.post<R<string>>("/node/activate", request)
            .then((response) => response.data.response),
    },
    pause: {
        key: () => ["node", "pause"],
        fn: (request: NodeRequest) => api.post<R<string>>("/node/pause", request)
            .then((response) => response.data.response),
    },
}