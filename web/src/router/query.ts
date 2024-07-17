import {useQuery} from "@tanstack/react-query";
import {QueryApi} from "../app/api";
import {QueryChartRequest, QueryConnection, QueryRunRequest, QueryType} from "../type/query";
import {useMutationAdapter} from "../hook/QueryCustom";

export function useRouterQueryList(type: QueryType, enabled: boolean = true) {
    return useQuery({
        queryKey: QueryApi.list.key(type),
        queryFn: () => QueryApi.list.fn(type),
        enabled,
    })
}

export function useRouterQueryUpdate(type: QueryType, onSuccess?: () => void) {
    return useMutationAdapter({
        mutationFn: QueryApi.update.fn,
        mutationKey: QueryApi.update.key(),
        successKeys: [QueryApi.list.key(type)],
        onSuccess: onSuccess,
    })
}

export function useRouterQueryDelete(type: QueryType, onSuccess?: () => void) {
    return useMutationAdapter({
        mutationFn: QueryApi.delete.fn,
        mutationKey: QueryApi.delete.key(),
        successKeys: [QueryApi.list.key(type)],
        onSuccess: onSuccess,
    })
}

export function useRouterQueryCreate(type: QueryType, onSuccess?: () => void) {
    return useMutationAdapter({
        mutationFn: QueryApi.create.fn,
        mutationKey: QueryApi.create.key(),
        successKeys: [QueryApi.list.key(type)],
        onSuccess: onSuccess,
    })
}

export function useRouterQueryRun(request: QueryRunRequest) {
    return useQuery({
        // eslint-disable-next-line @tanstack/query/exhaustive-deps
        queryKey: QueryApi.run.key(request.queryUuid),
        queryFn: () => QueryApi.run.fn(request.queryUuid ? {...request, query: undefined} : request),
        retry: false, refetchOnWindowFocus: false, structuralSharing: false,
    })
}

export function useRouterQueryChart(request: QueryChartRequest, enabled: boolean = false) {
    return useQuery({
        queryKey: QueryApi.chart.key(request),
        queryFn: () => QueryApi.chart.fn(request),
        retry: false, enabled, refetchOnWindowFocus: false,
    })
}

export function useRouterQueryDatabase(connection: QueryConnection, params: any, enabled: boolean = true) {
    return useQuery({
        // eslint-disable-next-line @tanstack/query/exhaustive-deps
        queryKey: [...QueryApi.databases.key(connection), params.name],
        queryFn: () => QueryApi.databases.fn({connection, ...params}),
        retry: false, enabled, refetchOnWindowFocus: false,
    })
}

export function useRouterQuerySchemas(connection: QueryConnection, params: any, enabled: boolean = true) {
    return useQuery({
        // eslint-disable-next-line @tanstack/query/exhaustive-deps
        queryKey: [...QueryApi.schemas.key(connection), params.name],
        queryFn: () => QueryApi.schemas.fn({connection, ...params}),
        retry: false, enabled, refetchOnWindowFocus: false,
    })
}

export function useRouterQueryTables(connection: QueryConnection, params: any, enabled: boolean = true) {
    return useQuery({
        // eslint-disable-next-line @tanstack/query/exhaustive-deps
        queryKey: [...QueryApi.tables.key(connection), params.name],
        queryFn: () => QueryApi.tables.fn({connection, ...params}),
        retry: false, enabled, refetchOnWindowFocus: false,
    })
}

export function useRouterQueryCancel(uuid?: string) {
    return useMutationAdapter({
        mutationFn: QueryApi.cancel.fn,
        mutationKey: QueryApi.cancel.key(),
        successKeys: [QueryApi.run.key(uuid)],
    })
}

export function useRouterQueryTerminate(uuid?: string) {
    return useMutationAdapter({
        mutationFn: QueryApi.terminate.fn,
        mutationKey: QueryApi.terminate.key(),
        successKeys: [QueryApi.run.key(uuid)],
    })
}

export function useRouterQueryLog(uuid: string) {
    return useQuery({
        queryKey: QueryApi.getLog.key(uuid),
        queryFn: () => QueryApi.getLog.fn(uuid),
        enabled: true, retry: false, refetchOnWindowFocus: false,
    })
}

export function useRouterQueryLogDelete(uuid: string) {
    return useMutationAdapter({
        mutationFn: QueryApi.deleteLog.fn,
        mutationKey: QueryApi.deleteLog.key(),
        successKeys: [QueryApi.getLog.key(uuid)],
    })
}
