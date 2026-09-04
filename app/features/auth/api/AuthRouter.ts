import {api} from "../../Api"
import {AuthConfigObject} from "../../config/api/ConfigType"
import {R} from "../../management/api/ManagementType"
import {UserAuthType} from "../../user/api/UserType"

export const AuthApi = {
    login: {
        key: () => ["login"],
        fn: async (data: {type: UserAuthType, subject: any}) => {
            const url = `/${data.type}/login`
            if (data.type === UserAuthType.OIDC) {
                window.location.href = api.defaults.baseURL + url
                return Promise.resolve(null)
            }
            const response = await api.post<R<any>>(url, data.subject)
            return response.data.response
        }
    },
    connect: {
        key: () => ["connect"],
        fn: (data: {type: UserAuthType, config: AuthConfigObject}) => api
            .post<R<any>>(`/${data.type}/connect`, data.config)
            .then((response) => response.data.response),
    },
    logout: {
        key: () => ["logout"],
        fn: () => api.post<R<any>>("/logout")
            .then((response) => response.data.response),
    },
}