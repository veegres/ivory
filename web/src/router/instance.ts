import {useQuery} from "@tanstack/react-query";
import {InstanceApi} from "../app/api";
import {InstanceRequest} from "../type/instance";
import {useMutationAdapter} from "../hook/QueryCustom";
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
    return useMutationAdapter({
        mutationFn: InstanceApi.updateConfig,
        successKeys: [InstanceApi.config.key(sidecar)],
        onSuccess: onSuccess,
    })
}

export function useRouterInstanceSwitchoverDelete(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.deleteSwitchover,
        successKeys: [InstanceApi.overview.key(cluster)]
    })
}

export function useRouterInstanceRestartDelete(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.deleteRestart,
        successKeys: [InstanceApi.overview.key(cluster)]
    })
}

export function useRouterInstanceRestart(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.restart,
        successKeys: [InstanceApi.overview.key(cluster)]
    })
}

export function useRouterInstanceReload(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.reload,
        successKeys: [InstanceApi.overview.key(cluster)]
    })
}

export function useRouterInstanceReinit(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.reinitialize,
        successKeys: [InstanceApi.overview.key(cluster)]
    })
}

export function useRouterInstanceSwitchover(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.switchover,
        successKeys: [InstanceApi.overview.key(cluster)]
    })
}

export function useRouterInstanceFailover(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.failover,
        successKeys: [InstanceApi.overview.key(cluster)]
    })
}

export function useRouterInstanceActivate(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.activate,
        successKeys: [InstanceApi.overview.key(cluster)]
    })
}

export function useRouterInstancePause(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.pause,
        successKeys: [InstanceApi.overview.key(cluster)]
    })
}
