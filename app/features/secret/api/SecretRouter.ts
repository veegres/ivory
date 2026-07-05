import {api} from "../../Api"
import {R} from "../../management/api/ManagementType"
import {SecretSetRequest} from "./SecretType"

export const SecretApi = {
    set: {
        key: () => ["secret", "set"],
        fn: (request: SecretSetRequest) => api.post<R<string>>("/initial/secret", request)
            .then((response) => response.data.response),
    },
    skip: {
        key: () => ["secret", "skip"],
        fn: () => api.post<R<string>>("/initial/skip")
            .then((response) => response.data.response),
    }
}