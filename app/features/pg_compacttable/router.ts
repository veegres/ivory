import {api} from "../api"
import {R} from "../management/type"
import {Bloat, BloatRequest} from "./type"

export const BloatApi = {
    list: {
        key: (cluster: string) => ["node", "pg_compacttable", "list", cluster],
        fn: (cluster: string) => api.get<R<Bloat[]>>(`/tool/bloat/cluster/${cluster}`)
            .then((response) => response.data.response),
    },
    logs: {
        key: (uuid: string) => ["node", "pg_compacttable", "logs", uuid],
        fn: (uuid: string) => api.get<string>(`/tool/bloat/${uuid}/logs`, {responseType: "text"})
            .then(({data}) => data === "" ? [] : data.split("\n")),
    },

    start: {
        key: () => ["node", "pg_compacttable", "job", "start"],
        fn: (ctr: BloatRequest) => api.post<R<Bloat>>("/tool/bloat/job/start", ctr)
            .then((response) => response.data.response),
    },
    stop: {
        key: () => ["node", "pg_compacttable", "job", "stop"],
        fn: (uuid: string) => api.post<R<Bloat>>(`/tool/bloat/job/${uuid}/stop`)
            .then((response) => response.data.response),
    },
    delete: {
        key: () => ["node", "pg_compacttable", "job", "delete"],
        fn: (uuid: string) => api.delete(`/tool/bloat/job/${uuid}/delete`)
            .then((response) => response.data.response),
    },

    stream: {
        key: () => ["node", "pg_compacttable", "job", "stream"],
        fn: (uuid: string) => new EventSource(`/api/tool/bloat/job/${uuid}/stream`)
    },
}