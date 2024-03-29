import {useMutation, useQuery} from "@tanstack/react-query";
import {InstanceApi} from "../app/api";
import {InstanceRequest} from "../type/instance";
import {useMutationOptions} from "../hook/QueryCustom";
import {Sidecar} from "../type/general";

export function useRouterInstanceOverview(cluster: string, request: InstanceRequest) {
    return useQuery({
        queryKey: InstanceApi.overview.key(cluster),
        queryFn: () => InstanceApi.overview.fn(request),
        retry: false,
    })
}

export function useRouterInstanceConfig(request: InstanceRequest, enabled: boolean) {
    return useQuery({
        queryKey: InstanceApi.config.key(request.sidecar),
        queryFn: () => InstanceApi.config.fn(request),
        enabled,
    })
}

export function useRouterInstanceConfigUpdate(sidecar: Sidecar, onSuccess: () => void) {
    const options = useMutationOptions([InstanceApi.config.key(sidecar)], onSuccess)
    return useMutation({mutationFn: InstanceApi.updateConfig, ...options})
}

export function useRouterInstanceSwitchoverDelete(cluster: string) {
    const options = useMutationOptions([InstanceApi.overview.key(cluster)])
    return useMutation({mutationFn: InstanceApi.deleteSwitchover, ...options})
}

export function useRouterInstanceRestartDelete(cluster: string) {
    const options = useMutationOptions([InstanceApi.overview.key(cluster)])
    return useMutation({mutationFn: InstanceApi.deleteRestart, ...options})
}

export function useRouterInstanceRestart(cluster: string) {
    const options = useMutationOptions([InstanceApi.overview.key(cluster)])
    return useMutation({mutationFn: InstanceApi.restart, ...options})
}

export function useRouterInstanceReload(cluster: string) {
    const options = useMutationOptions([InstanceApi.overview.key(cluster)])
    return useMutation({mutationFn: InstanceApi.reload, ...options})
}

export function useRouterInstanceReinit(cluster: string) {
    const options = useMutationOptions([InstanceApi.overview.key(cluster)])
    return useMutation({mutationFn: InstanceApi.reinitialize, ...options})
}

export function useRouterInstanceSwitchover(cluster: string) {
    const options = useMutationOptions([InstanceApi.overview.key(cluster)])
    return useMutation({mutationFn: InstanceApi.switchover, ...options})
}

export function useRouterInstanceFailover(cluster: string) {
    const options = useMutationOptions([InstanceApi.overview.key(cluster)])
    return useMutation({mutationFn: InstanceApi.failover, ...options})
}

export function useRouterInstanceActivate(cluster: string) {
    const options = useMutationOptions([InstanceApi.overview.key(cluster)])
    return useMutation({mutationFn: InstanceApi.activate, ...options})
}

export function useRouterInstancePause(cluster: string) {
    const options = useMutationOptions([InstanceApi.overview.key(cluster)])
    return useMutation({mutationFn: InstanceApi.pause, ...options})
}
