import {api} from "../api";
import {R} from "../type";
import {AppConfig, AppInfo, Login} from "./type";
import {SecretSetRequest, SecretUpdateRequest} from "../secret/type";

export const GeneralApi = {
    info: {
        key: () => ["info"],
        fn: () => api.get<R<AppInfo>>(`/info`)
            .then((response) => response.data.response),
    },
    login: {
        key: () => ["login"],
        fn: (req: Login) => api.post<R<any>>(`/login`, req)
            .then((response) => response.data.response),
    },
    setConfig: {
        key: () => ["config", "set"],
        fn: (request: AppConfig) => api.post<R<string>>(`/initial/config`, request)
            .then((response) => response.data.response),
    }
}

export const InitialApi = {
    setSecret: {
        key: () => ["secret", "set"],
        fn: (request: SecretSetRequest) => api.post<R<string>>(`/initial/secret`, request)
            .then((response) => response.data.response),
    },
    erase: {
        key: () => ["erase"],
        fn: () => api.delete<R<string>>(`/initial/erase`)
            .then((response) => response.data.response),
    },
}

export const SafeApi = {
    changeSecret: {
        key: () => ["secret", "change"],
        fn: (request: SecretUpdateRequest) => api.post<R<string>>(`/safe/secret`, request)
            .then((response) => response.data.response),
    },
    erase: {
        key: () => ["erase"],
        fn: () => api.delete<R<string>>(`/safe/erase`)
            .then((response) => response.data.response),
    },
}