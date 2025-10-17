import {api} from "../api";
import {AuthType} from "../config/type";
import {R} from "../management/type";

export const AuthApi = {
    login: {
        key: (type: AuthType) => [AuthType[type], "login"],
        fn: (type: AuthType) => async (req: any) => {
            const url = `/${AuthType[type].toLowerCase()}/login`
            if (type === AuthType.OIDC) {
                window.location.href = `/api` + url
                return Promise.resolve(null)
            }
            const response = await api.post<R<any>>(url, req);
            return response.data.response;
        }
    },
    logout: {
        key: () => ["logout"],
        fn: () => api.post<R<any>>(`/logout`)
            .then((response) => response.data.response),
    },
}