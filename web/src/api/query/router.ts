import {api} from "../api"
import {R} from "../management/type"
import {
    Query, QueryChart, QueryChartRequest,
    QueryConnection,
    QueryDatabasesRequest,
    QueryFields, QueryKillRequest,
    QueryRequest,
    QueryRunRequest, QuerySchemasRequest, QueryTablesRequest,
    QueryType
} from "./type"

export const QueryApi = {
    list: {
        key: (type?: QueryType) => ["query", "list", type],
        fn: (type?: QueryType) => api.get<R<Query[]>>("/query", {params: {type}})
            .then((response) => response.data.response),
    },
    update: {
        key: () => ["query", "update"],
        fn: ({id, query}: { id: string, query: QueryRequest }) => api.put<R<Query>>(`/query/${id}`, query)
            .then((response) => response.data.response),
    },
    create: {
        key: () => ["query", "create"],
        fn: (query: QueryRequest) => api.post<R<Query>>("/query", query)
            .then((response) => response.data.response),
    },
    delete: {
        key: () => ["query", "delete"],
        fn: (uuid: string) => api.delete<R<string>>(`/query/${uuid}`)
            .then((response) => response.data.response),
    },
    console: {
        key: (type: string) => ["query", "console", type],
        fn: (req: QueryRunRequest) => api.post<R<QueryFields>>("/query/execute/console", req)
            .then((response) => response.data.response),
    },
    activity: {
        key: () => QueryApi.console.key("activity"),
        fn: (req: QueryConnection) => api.post<R<QueryFields>>("/query/execute/activity", req)
            .then((response) => response.data.response),
    },
    databases: {
        key: (con: QueryConnection) => ["query", "databases", con],
        fn: (req: QueryDatabasesRequest) => api.post<R<string[]>>("/query/execute/databases", req)
            .then((response) => response.data.response),
    },
    schemas: {
        key: (con: QueryConnection) => ["query", "schemas", con],
        fn: (req: QuerySchemasRequest) => api.post<R<string[]>>("/query/execute/schemas", req)
            .then((response) => response.data.response),
    },
    tables: {
        key: (con: QueryConnection, schema?: string) => ["query", "tables", con, schema ?? "-"],
        fn: (req: QueryTablesRequest) => api.post<R<string[]>>("/query/execute/tables", req)
            .then((response) => response.data.response),
    },
    chart: {
        key: (req: QueryChartRequest) => ["query", "chart", req.type, req.connection],
        fn: (req: QueryChartRequest) => api.post<R<QueryChart>>("/query/execute/chart", req)
            .then((response) => response.data.response),
    },
    cancel: {
        key: () => ["query", "cancel"],
        fn: (req: QueryKillRequest) => api.post<R<string>>("/query/execute/cancel", req)
            .then((response) => response.data.response),
    },
    terminate: {
        key: () => ["query", "terminate"],
        fn: (req: QueryKillRequest) => api.post<R<string>>("/query/execute/terminate", req)
            .then((response) => response.data.response),
    },

    getLog: {
        key: (uuid: string) => ["query", "log", uuid],
        fn: (uuid: string) => api.get<R<QueryFields[]>>(`/query/log/${uuid}`)
            .then((response) => response.data.response),
    },
    deleteLog: {
        key: () => ["query", "log", "delete"],
        fn: (uuid: string) => api.delete<R<string>>(`/query/log/${uuid}`)
            .then((response) => response.data.response),
    },
}