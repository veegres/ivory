import {api} from "../api"
import {R} from "../management/type"
import {
    Chart,
    ChartRequest,
    Connection,
    ConsoleRequest,
    DatabasesRequest,
    DbResponse,
    KillRequest,
    Request,
    Response,
    SchemasRequest,
    TablesRequest,
    TemplateRequest,
    Type,
} from "./type"

export const QueryApi = {
    list: {
        key: (type?: Type) => ["query", "list", type],
        fn: (type?: Type) => api.get<R<Response[]>>("/query", {params: {type}})
            .then((response) => response.data.response),
    },
    update: {
        key: () => ["query", "update"],
        fn: ({id, query}: { id: string, query: Request }) => api.put<R<Response>>(`/query/${id}`, query)
            .then((response) => response.data.response),
    },
    create: {
        key: () => ["query", "create"],
        fn: (query: Request) => api.post<R<Response>>("/query", query)
            .then((response) => response.data.response),
    },
    delete: {
        key: () => ["query", "delete"],
        fn: (uuid: string) => api.delete<R<string>>(`/query/${uuid}`)
            .then((response) => response.data.response),
    },
    console: {
        key: () => ["query", "execute", "console"],
        fn: (req: ConsoleRequest) => api.post<R<DbResponse>>("/query/execute/console", req)
            .then((response) => response.data.response),
    },
    template: {
        key: (id?: string) => ["query", "execute", "template", id],
        fn: (req: TemplateRequest) => api.post<R<DbResponse>>("/query/execute/template", req)
            .then((response) => response.data.response),
    },
    activity: {
        key: () => ["query", "execute", "activity"],
        fn: (req: Connection) => api.post<R<DbResponse>>("/query/execute/activity", req)
            .then((response) => response.data.response),
    },
    databases: {
        key: (con: Connection) => ["query", "execute", "databases", con],
        fn: (req: DatabasesRequest) => api.post<R<string[]>>("/query/execute/databases", req)
            .then((response) => response.data.response),
    },
    schemas: {
        key: (con: Connection) => ["query", "execute", "schemas", con],
        fn: (req: SchemasRequest) => api.post<R<string[]>>("/query/execute/schemas", req)
            .then((response) => response.data.response),
    },
    tables: {
        key: (con: Connection, schema?: string) => ["query", "execute", "tables", con, schema ?? "-"],
        fn: (req: TablesRequest) => api.post<R<string[]>>("/query/execute/tables", req)
            .then((response) => response.data.response),
    },
    chart: {
        key: (req: ChartRequest) => ["query", "execute", "chart", req.type, req.connection],
        fn: (req: ChartRequest) => api.post<R<Chart>>("/query/execute/chart", req)
            .then((response) => response.data.response),
    },
    cancel: {
        key: () => ["query", "execute", "cancel"],
        fn: (req: KillRequest) => api.post<R<string>>("/query/execute/cancel", req)
            .then((response) => response.data.response),
    },
    terminate: {
        key: () => ["query", "execute", "terminate"],
        fn: (req: KillRequest) => api.post<R<string>>("/query/execute/terminate", req)
            .then((response) => response.data.response),
    },

    getLog: {
        key: (uuid: string) => ["query", "log", uuid],
        fn: (uuid: string) => api.get<R<DbResponse[]>>(`/query/log/${uuid}`)
            .then((response) => response.data.response),
    },
    deleteLog: {
        key: () => ["query", "log", "delete"],
        fn: (uuid: string) => api.delete<R<string>>(`/query/log/${uuid}`)
            .then((response) => response.data.response),
    },
}
