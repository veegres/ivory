import {api} from "../api"
import {AuthConfigObject} from "../config/type"
import {R} from "../management/type"
import {AuthType} from "./type"

export const AuthApi = {
    login: {
        key: () => ["login"],
        fn: async (data: {type: AuthType, subject: any}) => {
            const url = `/${AuthType[data.type].toLowerCase()}/login`
            if (data.type === AuthType.OIDC) {
                window.location.href = api.defaults.baseURL + url
                return Promise.resolve(null)
            }
            const response = await api.post<R<any>>(url, data.subject)
            return response.data.response
        }
    },
    connect: {
        key: () => ["connect"],
        fn: (data: {type: AuthType, config: AuthConfigObject}) => api
            .post<R<any>>(`/${AuthType[data.type].toLowerCase()}/connect`, data.config)
            .then((response) => response.data.response),
    },
    logout: {
        key: () => ["logout"],
        fn: () => api.post<R<any>>("/logout")
            .then((response) => response.data.response),
    },
}