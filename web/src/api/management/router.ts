import {api} from "../api";
import {AppInfo, R, SecretUpdateRequest} from "./type";
import {Login} from "../auth/type";

export const ManagementApi = {
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
    changeSecret: {
        key: () => ["secret", "change"],
        fn: (request: SecretUpdateRequest) => api.post<R<string>>(`/safe/secret`, request)
            .then((response) => response.data.response),
    },
    erase: {
        key: () => ["erase"],
        fn: (type: "initial" | "safe") => api.delete<R<string>>(`/${type}/erase`)
            .then((response) => response.data.response),
    },
}
