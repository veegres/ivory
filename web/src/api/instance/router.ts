import {api} from "../api"
import {R} from "../management/type"
import {InstanceRequest, Sidecar} from "./type"


export const InstanceApi = {
    config: {
        key: (sidecar: Sidecar) => ["instance", "config", sidecar.host, sidecar.port],
        fn: (request: InstanceRequest) => api.get<R<any>>("/instance/config", {params: {request: JSON.stringify(request)}})
            .then((response) => response.data.response),
    },
    updateConfig: {
        key: () => ["instance", "config", "update"],
        fn: (request: InstanceRequest) => api.patch<R<any>>("/instance/config", request)
            .then((response) => response.data.response),
    },
    switchover: {
        key: () => ["instance", "switchover"],
        fn: (request: InstanceRequest) => api.post<R<string>>("/instance/switchover", request)
            .then((response) => response.data.response),
    },
    deleteSwitchover: {
        key: () => ["instance", "switchover", "delete"],
        fn: (request: InstanceRequest) => api.delete<R<string>>("/instance/switchover", {params: {request: JSON.stringify(request)}})
            .then((response) => response.data.response),
    },
    reinitialize: {
        key: () => ["instance", "reinitialize"],
        fn: (request: InstanceRequest) => api.post<R<string>>("/instance/reinitialize", request)
            .then((response) => response.data.response),
    },
    restart: {
        key: () => ["instance", "restart"],
        fn: (request: InstanceRequest) => api.post<R<string>>("/instance/restart", request)
            .then((response) => response.data.response),
    },
    deleteRestart: {
        key: () => ["instance", "restart", "delete"],
        fn: (request: InstanceRequest) => api.delete<R<string>>("/instance/restart", {params: {request: JSON.stringify(request)}})
            .then((response) => response.data.response),
    },
    reload: {
        key: () => ["instance", "reload"],
        fn: (request: InstanceRequest) => api.post<R<string>>("/instance/reload", request)
            .then((response) => response.data.response),
    },
    failover: {
        key: () => ["instance", "failover"],
        fn: (request: InstanceRequest) => api.post<R<string>>("/instance/failover", request)
            .then((response) => response.data.response),
    },
    activate: {
        key: () => ["instance", "activate"],
        fn: (request: InstanceRequest) => api.post<R<string>>("/instance/activate", request)
            .then((response) => response.data.response),
    },
    pause: {
        key: () => ["instance", "pause"],
        fn: (request: InstanceRequest) => api.post<R<string>>("/instance/pause", request)
            .then((response) => response.data.response),
    },
}