import {SecretSetRequest} from "./type";
import {api} from "../api";

import {R} from "../management/type";

export const SecretApi = {
    setSecret: {
        key: () => ["secret", "set"],
        fn: (request: SecretSetRequest) => api.post<R<string>>(`/initial/secret`, request)
            .then((response) => response.data.response),
    },
}