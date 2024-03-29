import {useMutation, useQuery} from "@tanstack/react-query";
import {QueryApi} from "../app/api";
import {QueryChartRequest, QueryRunRequest, QueryType} from "../type/query";
import {useMutationOptions} from "../hook/QueryCustom";

export function useRouterQueryList(type: QueryType, enabled: boolean = true) {
    return useQuery({
        queryKey: ["query", "list", type],
        queryFn: () => QueryApi.list(type),
        enabled,
    })
}

export function useRouterQueryUpdate(type: QueryType, onSuccess?: () => void) {
    const options = useMutationOptions([["query", "list", type]], onSuccess)
    return useMutation({mutationFn: QueryApi.update, ...options})
}

export function useRouterQueryDelete(type: QueryType, onSuccess?: () => void) {
    const options = useMutationOptions([["query", "list", type]], onSuccess)
    return useMutation({mutationFn: QueryApi.delete, ...options})
}

export function useRouterQueryCreate(type: QueryType, onSuccess?: () => void) {
    const options = useMutationOptions([["query", "list", type]], onSuccess)
    return useMutation({mutationFn: QueryApi.create, ...options})
}

export function useRouterQueryRun(request: QueryRunRequest) {
    return useQuery({
        queryKey: ["query", "run", request.queryUuid ?? "standalone"],
        queryFn: () => QueryApi.run(request.queryUuid ? {...request, query: undefined} : request),
        enabled: true, retry: false, refetchOnWindowFocus: false,
    })
}

export function useRouterQueryCancel(uuid?: string) {
    const options = useMutationOptions([["query", "run", uuid ?? "standalone"]])
    return useMutation({mutationFn: QueryApi.cancel, ...options})
}

export function useRouterQueryTerminate(uuid?: string) {
    const options = useMutationOptions([["query", "run", uuid ?? "standalone"]])
    return useMutation({mutationFn: QueryApi.terminate, ...options})
}

export function useRouterQueryLog(uuid: string) {
    return useQuery({
        queryKey: ["query", "log", uuid],
        queryFn: () => QueryApi.getLog(uuid),
        enabled: true, retry: false, refetchOnWindowFocus: false,
    })
}

export function useRouterQueryLogDelete(uuid: string) {
    const options = useMutationOptions([["query", "log", uuid]])
    return useMutation({mutationFn: () => QueryApi.deleteLog(uuid), ...options})
}

export function useRouterQueryChart(request: QueryChartRequest, enabled: boolean = true) {
    return useQuery({
        queryKey: ["query", "chart", request.type, request.db.host, request.db.port, request.db.name, request.credentialId],
        queryFn: () => QueryApi.chart(request),
        retry: false, enabled,
    })
}
