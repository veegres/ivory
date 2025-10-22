import {api} from "../api";
import {NewAppConfig} from "./type";
import {R} from "../management/type";

export const ConfigApi = {
    setAppConfig: {
        key: () => ["config", "set"],
        fn: (request: NewAppConfig) => api.post<R<string>>(`/config`, request)
            .then((response) => response.data.response),
    },
}