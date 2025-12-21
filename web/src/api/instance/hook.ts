import {useQuery} from "@tanstack/react-query"

import {useMutationAdapter} from "../../hook/QueryCustom"
import {ClusterApi} from "../cluster/router"
import {InstanceApi} from "./router"
import {InstanceRequest, Sidecar} from "./type"

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
        successKeys: [ClusterApi.overview.key(cluster)]
    })
}

export function useRouterInstanceRestartDelete(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.deleteRestart.fn,
        mutationKey: InstanceApi.deleteRestart.key(),
        successKeys: [ClusterApi.overview.key(cluster)]
    })
}

export function useRouterInstanceRestart(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.restart.fn,
        mutationKey: InstanceApi.restart.key(),
        successKeys: [ClusterApi.overview.key(cluster)]
    })
}

export function useRouterInstanceReload(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.reload.fn,
        mutationKey: InstanceApi.reload.key(),
        successKeys: [ClusterApi.overview.key(cluster)]
    })
}

export function useRouterInstanceReinit(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.reinitialize.fn,
        mutationKey: InstanceApi.reinitialize.key(),
        successKeys: [ClusterApi.overview.key(cluster)]
    })
}

export function useRouterInstanceSwitchover(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.switchover.fn,
        mutationKey: InstanceApi.switchover.key(),
        successKeys: [ClusterApi.overview.key(cluster)]
    })
}

export function useRouterInstanceFailover(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.failover.fn,
        mutationKey: InstanceApi.failover.key(),
        successKeys: [ClusterApi.overview.key(cluster)]
    })
}

export function useRouterInstanceActivate(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.activate.fn,
        mutationKey: InstanceApi.activate.key(),
        successKeys: [ClusterApi.overview.key(cluster)]
    })
}

export function useRouterInstancePause(cluster: string) {
    return useMutationAdapter({
        mutationFn: InstanceApi.pause.fn,
        mutationKey: InstanceApi.pause.key(),
        successKeys: [ClusterApi.overview.key(cluster)]
    })
}
