import {useMutation, useQuery} from "@tanstack/react-query";
import {InstanceApi} from "../app/api";
import {InstanceRequest} from "../type/instance";
import {useMutationOptions} from "../hook/QueryCustom";
import {Sidecar} from "../type/general";

export function useRouterInstanceOverview(cluster: string, request: InstanceRequest) {
    return useQuery({
        queryKey: ["instance", "overview", cluster],
        queryFn: () => InstanceApi.overview(request),
        retry: false,
    })
}

export function useRouterInstanceConfig(request: InstanceRequest, enabled: boolean) {
    return useQuery({
        queryKey: ["instance", "config", request.sidecar.host, request.sidecar.port],
        queryFn: () => InstanceApi.config(request),
        enabled,
    })
}

export function useRouterInstanceConfigUpdate(sidecar: Sidecar, onSuccess: () => void) {
    const options = useMutationOptions([["instance", "config", sidecar.host, sidecar.port]], onSuccess)
    return useMutation({mutationFn: InstanceApi.updateConfig, ...options})
}

export function useRouterInstanceSwitchoverDelete(cluster: string) {
    const options = useMutationOptions([["instance", "overview", cluster]])
    return useMutation({mutationFn: InstanceApi.deleteSwitchover, ...options})
}

export function useRouterInstanceRestartDelete(cluster: string) {
    const options = useMutationOptions([["instance", "overview", cluster]])
    return useMutation({mutationFn: InstanceApi.deleteRestart, ...options})
}

export function useRouterInstanceRestart(cluster: string) {
    const options = useMutationOptions([["instance", "overview", cluster]])
    return useMutation({mutationFn: InstanceApi.restart, ...options})
}

export function useRouterInstanceReload(cluster: string) {
    const options = useMutationOptions([["instance", "overview", cluster]])
    return useMutation({mutationFn: InstanceApi.reload, ...options})
}

export function useRouterInstanceReinit(cluster: string) {
    const options = useMutationOptions([["instance", "overview", cluster]])
    return useMutation({mutationFn: InstanceApi.reinitialize, ...options})
}

export function useRouterInstanceSwitchover(cluster: string) {
    const options = useMutationOptions([["instance", "overview", cluster]])
    return useMutation({mutationFn: InstanceApi.switchover, ...options})
}

export function useRouterInstanceFailover(cluster: string) {
    const options = useMutationOptions([["instance", "overview", cluster]])
    return useMutation({mutationFn: InstanceApi.failover, ...options})
}

export function useRouterInstanceActivate(cluster: string) {
    const options = useMutationOptions([["instance", "overview", cluster]])
    return useMutation({mutationFn: InstanceApi.activate, ...options})
}

export function useRouterInstancePause(cluster: string) {
    const options = useMutationOptions([["instance", "overview", cluster]])
    return useMutation({mutationFn: InstanceApi.pause, ...options})
}
