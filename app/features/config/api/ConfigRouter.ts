import {api} from "../../Api"
import {R} from "../../management/api/ManagementType"
import {NewAppConfig} from "./ConfigType"

export const ConfigApi = {
    setAppConfig: {
        key: () => ["config", "set"],
        fn: (request: NewAppConfig) => api.post<R<string>>("/config", request)
            .then((response) => response.data.response),
    },
}