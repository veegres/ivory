import {useQuery} from "@tanstack/react-query";
import {QueryApi} from "../app/api";
import {QueryChartRequest, QueryRunRequest, QueryType} from "../type/query";

export function useRouterQueryRun(request: QueryRunRequest) {
    return useQuery({
        queryKey: ["query", "run", request.queryUuid ?? "standalone"],
        queryFn: () => QueryApi.run(request.queryUuid ? {...request, query: undefined} : request),
        enabled: true, retry: false, refetchOnWindowFocus: false,
    })
}

export function useRouterQueryMap(type: QueryType, enabled: boolean = true) {
    return useQuery({
        queryKey: ["query", "map", type],
        queryFn: () => QueryApi.list(type),
        enabled,
    })
}

export function useRouterQueryLog(uuid: string) {
    return useQuery({
        queryKey: ["query", "log", uuid],
        queryFn: () => QueryApi.getLog(uuid),
        enabled: true, retry: false, refetchOnWindowFocus: false,
    })
}

export function useRouterQueryChart(request: QueryChartRequest, enabled: boolean = true) {
    return useQuery({
        queryKey: ["query", "chart", request.type, request.db.host, request.db.port, request.db.name, request.credentialId],
        queryFn: () => QueryApi.chart(request),
        retry: false, enabled,
    })
}
