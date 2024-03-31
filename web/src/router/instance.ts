import {useQuery} from "@tanstack/react-query";
import {InstanceApi} from "../app/api";
import {InstanceRequest} from "../type/instance";
import {useMutationAdapter} from "../hook/QueryCustom";
import {Sidecar} from "../type/general";

// TODO #414 it is accept request as a function right now, because otherwise it gets incorrect state.
//  It is considered as a hack to previous solution, idea is to get rid of InstanceDetection on
//  the frontend, because it is to complicated and generates a lot of renders with useQuery()
export function useRouterInstanceOverview(cluster: string, request: () => InstanceRequest) {
    return useQuery({
        queryKey: InstanceApi.overview.key(cluster),
        queryFn: () => InstanceApi.overview.fn(request()),
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
        mutationFn: InstanceApi.updateConfig.fn,
        successKeys: [InstanceApi.config.key(sidecar)],
        onSuccess: onSuccess,
    })
}

export function useRouterInstanceSwitchoverDelete(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.deleteSwitchover.fn,
        successKeys: [InstanceApi.overview.key(cluster)]
    })
}

export function useRouterInstanceRestartDelete(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.deleteRestart.fn,
        successKeys: [InstanceApi.overview.key(cluster)]
    })
}

export function useRouterInstanceRestart(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.restart.fn,
        successKeys: [InstanceApi.overview.key(cluster)]
    })
}

export function useRouterInstanceReload(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.reload.fn,
        successKeys: [InstanceApi.overview.key(cluster)]
    })
}

export function useRouterInstanceReinit(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.reinitialize.fn,
        successKeys: [InstanceApi.overview.key(cluster)]
    })
}

export function useRouterInstanceSwitchover(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.switchover.fn,
        successKeys: [InstanceApi.overview.key(cluster)]
    })
}

export function useRouterInstanceFailover(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.failover.fn,
        successKeys: [InstanceApi.overview.key(cluster)]
    })
}

export function useRouterInstanceActivate(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.activate.fn,
        successKeys: [InstanceApi.overview.key(cluster)]
    })
}

export function useRouterInstancePause(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.pause.fn,
        successKeys: [InstanceApi.overview.key(cluster)]
    })
}
