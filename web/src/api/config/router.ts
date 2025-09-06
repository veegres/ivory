import {api} from "../api";
import {AppConfig} from "./type";
import {R} from "../management/type";

export const ConfigApi = {
    setAppConfig: {
        key: () => ["config", "set"],
        fn: (request: AppConfig) => api.post<R<string>>(`/initial/config`, request)
            .then((response) => response.data.response),
    },
}