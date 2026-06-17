import {api} from "../../api"
import {R} from "../../management/api/type"

export const TagApi = {
    list: {
        key: () => ["tag", "list"],
        fn: () => api.get<R<string[]>>("/tag")
            .then((response) => response.data.response),
    },
}