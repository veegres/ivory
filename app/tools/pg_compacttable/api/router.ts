import {api} from "../../../features/api"
import {R} from "../../../features/management/api/type"
import {PgCompactTable, PgCompactTableRequest} from "./type"

export const PgCompactTableApi = {
    list: {
        key: (cluster: string) => ["node", "pg_compacttable", "list", cluster],
        fn: (cluster: string) => api.get<R<PgCompactTable[]>>(`/tool/bloat/cluster/${cluster}`)
            .then((response) => response.data.response),
    },
    logs: {
        key: (uuid: string) => ["node", "pg_compacttable", "logs", uuid],
        fn: (uuid: string) => api.get<string>(`/tool/bloat/${uuid}/logs`, {responseType: "text"})
            .then(({data}) => data === "" ? [] : data.split("\n")),
    },

    start: {
        key: () => ["node", "pg_compacttable", "job", "start"],
        fn: (ctr: PgCompactTableRequest) => api.post<R<PgCompactTable>>("/tool/bloat/job/start", ctr)
            .then((response) => response.data.response),
    },
    stop: {
        key: () => ["node", "pg_compacttable", "job", "stop"],
        fn: (uuid: string) => api.post<R<PgCompactTable>>(`/tool/bloat/job/${uuid}/stop`)
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