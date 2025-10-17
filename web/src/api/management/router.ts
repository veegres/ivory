import {api} from "../api";
import {AppInfo, R, SecretUpdateRequest} from "./type";

export const ManagementApi = {
    info: {
        key: () => ["info"],
        fn: () => api.get<R<AppInfo>>(`/info`)
            .then((response) => response.data.response),
    },
    changeSecret: {
        key: () => ["secret", "change"],
        fn: (request: SecretUpdateRequest) => api.post<R<string>>(`/management/secret`, request)
            .then((response) => response.data.response),
    },
    erase: {
        key: () => ["erase"],
        fn: (type: "initial" | "management") => api.delete<R<string>>(`/${type}/erase`)
            .then((response) => response.data.response),
    },
}
