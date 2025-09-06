import {useQuery} from "@tanstack/react-query";
import {InstanceRequest, Sidecar} from "./type";
import {useMutationAdapter} from "../../hook/QueryCustom";
import {InstanceApi} from "./router";

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
        // eslint-disable-next-line @tanstack/query/exhaustive-deps
        queryKey: InstanceApi.config.key(request.sidecar),
        queryFn: () => InstanceApi.config.fn(request),
        enabled,
    })
}

export function useRouterInstanceConfigUpdate(sidecar: Sidecar, onSuccess: () => void) {
    return useMutationAdapter({
        mutationFn: InstanceApi.updateConfig.fn,
        mutationKey: InstanceApi.updateConfig.key(),
        successKeys: [InstanceApi.config.key(sidecar)],
        onSuccess: onSuccess,
    })
}

export function useRouterInstanceSwitchoverDelete(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.deleteSwitchover.fn,
        mutationKey: InstanceApi.deleteSwitchover.key(),
        successKeys: [InstanceApi.overview.key(cluster)]
    })
}

export function useRouterInstanceRestartDelete(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.deleteRestart.fn,
        mutationKey: InstanceApi.deleteRestart.key(),
        successKeys: [InstanceApi.overview.key(cluster)]
    })
}

export function useRouterInstanceRestart(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.restart.fn,
        mutationKey: InstanceApi.restart.key(),
        successKeys: [InstanceApi.overview.key(cluster)]
    })
}

export function useRouterInstanceReload(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.reload.fn,
        mutationKey: InstanceApi.reload.key(),
        successKeys: [InstanceApi.overview.key(cluster)]
    })
}

export function useRouterInstanceReinit(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.reinitialize.fn,
        mutationKey: InstanceApi.reinitialize.key(),
        successKeys: [InstanceApi.overview.key(cluster)]
    })
}

export function useRouterInstanceSwitchover(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.switchover.fn,
        mutationKey: InstanceApi.switchover.key(),
        successKeys: [InstanceApi.overview.key(cluster)]
    })
}

export function useRouterInstanceFailover(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.failover.fn,
        mutationKey: InstanceApi.failover.key(),
        successKeys: [InstanceApi.overview.key(cluster)]
    })
}

export function useRouterInstanceActivate(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.activate.fn,
        mutationKey: InstanceApi.activate.key(),
        successKeys: [InstanceApi.overview.key(cluster)]
    })
}

export function useRouterInstancePause(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.pause.fn,
        mutationKey: InstanceApi.pause.key(),
        successKeys: [InstanceApi.overview.key(cluster)]
    })
}
