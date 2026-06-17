import {api} from "../../api"
import {R} from "../../management/api/type"
import {NewAppConfig} from "./type"

export const ConfigApi = {
    setAppConfig: {
        key: () => ["config", "set"],
        fn: (request: NewAppConfig) => api.post<R<string>>("/config", request)
            .then((response) => response.data.response),
    },
}