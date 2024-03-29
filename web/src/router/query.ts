import {useMutation, useQuery} from "@tanstack/react-query";
import {QueryApi} from "../app/api";
import {QueryChartRequest, QueryRunRequest, QueryType} from "../type/query";
import {useMutationOptions} from "../hook/QueryCustom";

export function useRouterQueryList(type: QueryType, enabled: boolean = true) {
    return useQuery({
        queryKey: QueryApi.list.key(type),
        queryFn: () => QueryApi.list.fn(type),
        enabled,
    })
}

export function useRouterQueryUpdate(type: QueryType, onSuccess?: () => void) {
    const options = useMutationOptions([QueryApi.list.key(type)], onSuccess)
    return useMutation({mutationFn: QueryApi.update, ...options})
}

export function useRouterQueryDelete(type: QueryType, onSuccess?: () => void) {
    const options = useMutationOptions([QueryApi.list.key(type)], onSuccess)
    return useMutation({mutationFn: QueryApi.delete, ...options})
}

export function useRouterQueryCreate(type: QueryType, onSuccess?: () => void) {
    const options = useMutationOptions([QueryApi.list.key(type)], onSuccess)
    return useMutation({mutationFn: QueryApi.create, ...options})
}

export function useRouterQueryRun(request: QueryRunRequest) {
    return useQuery({
        queryKey: QueryApi.run.key(request.queryUuid),
        queryFn: () => QueryApi.run.fn(request.queryUuid ? {...request, query: undefined} : request),
        enabled: true, retry: false, refetchOnWindowFocus: false,
    })
}

export function useRouterQueryCancel(uuid?: string) {
    const options = useMutationOptions([QueryApi.run.key(uuid)])
    return useMutation({mutationFn: QueryApi.cancel, ...options})
}

export function useRouterQueryTerminate(uuid?: string) {
    const options = useMutationOptions([QueryApi.run.key(uuid)])
    return useMutation({mutationFn: QueryApi.terminate, ...options})
}

export function useRouterQueryLog(uuid: string) {
    return useQuery({
        queryKey: QueryApi.getLog.key(uuid),
        queryFn: () => QueryApi.getLog.fn(uuid),
        enabled: true, retry: false, refetchOnWindowFocus: false,
    })
}

export function useRouterQueryLogDelete(uuid: string) {
    const options = useMutationOptions([QueryApi.getLog.key(uuid)])
    return useMutation({mutationFn: () => QueryApi.deleteLog(uuid), ...options})
}

export function useRouterQueryChart(request: QueryChartRequest, enabled: boolean = true) {
    return useQuery({
        queryKey: QueryApi.chart.key(request),
        queryFn: () => QueryApi.chart.fn(request),
        retry: false, enabled,
    })
}
