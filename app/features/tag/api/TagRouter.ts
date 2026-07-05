import {api} from "../../Api"
import {R} from "../../management/api/ManagementType"

export const TagApi = {
    list: {
        key: () => ["tag", "list"],
        fn: () => api.get<R<string[]>>("/tag")
            .then((response) => response.data.response),
    },
}