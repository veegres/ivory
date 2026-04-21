import {api} from "../api"
import {R} from "../management/type"
import {Bloat, BloatRequest} from "./type"

export const BloatApi = {
    list: {
        key: (cluster: string) => ["node", "bloat", "list", cluster],
        fn: (cluster: string) => api.get<R<Bloat[]>>(`/cli/bloat/cluster/${cluster}`)
            .then((response) => response.data.response),
    },
    logs: {
        key: (uuid: string) => ["node", "bloat", "logs", uuid],
        fn: (uuid: string) => api.get<string>(`/cli/bloat/${uuid}/logs`, {responseType: "text"})
            .then(({data}) => data === "" ? [] : data.split("\n")),
    },

    start: {
        key: () => ["node", "bloat", "job", "start"],
        fn: (ctr: BloatRequest) => api.post<R<Bloat>>("/cli/bloat/job/start", ctr)
            .then((response) => response.data.response),
    },
    stop: {
        key: () => ["node", "bloat", "job", "stop"],
        fn: (uuid: string) => api.post<R<Bloat>>(`/cli/bloat/job/${uuid}/stop`)
            .then((response) => response.data.response),
    },
    delete: {
        key: () => ["node", "bloat", "job", "delete"],
        fn: (uuid: string) => api.delete(`/cli/bloat/job/${uuid}/delete`)
            .then((response) => response.data.response),
    },

    stream: {
        key: () => ["node", "bloat", "job", "stream"],
        fn: (uuid: string) => new EventSource(`/api/cli/bloat/job/${uuid}/stream`)
    },
}