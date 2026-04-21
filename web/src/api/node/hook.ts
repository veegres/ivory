import {useQuery} from "@tanstack/react-query"

import {useMutationAdapter} from "../../hook/QueryCustom"
import {ClusterApi} from "../cluster/router"
import {NodeApi} from "./router"
import {NodeRequest, Sidecar} from "./type"

export function useRouterNodeConfig(request: NodeRequest, enabled: boolean) {
    return useQuery({
        // eslint-disable-next-line @tanstack/query/exhaustive-deps
        queryKey: NodeApi.config.key(request.sidecar),
        queryFn: () => NodeApi.config.fn(request),
        enabled,
    })
}

export function useRouterNodeConfigUpdate(sidecar: Sidecar, onSuccess: () => void) {
    return useMutationAdapter({
        mutationFn: NodeApi.updateConfig.fn,
        mutationKey: NodeApi.updateConfig.key(),
        successKeys: [NodeApi.config.key(sidecar)],
        onSuccess: onSuccess,
    })
}

export function useRouterNodeSwitchoverDelete(cluster: string) {
    return useMutationAdapter({
        mutationFn: NodeApi.deleteSwitchover.fn,
        mutationKey: NodeApi.deleteSwitchover.key(),
        successKeys: [ClusterApi.overview.key(cluster)]
    })
}

export function useRouterNodeRestartDelete(cluster: string) {
    return useMutationAdapter({
        mutationFn: NodeApi.deleteRestart.fn,
        mutationKey: NodeApi.deleteRestart.key(),
        successKeys: [ClusterApi.overview.key(cluster)]
    })
}

export function useRouterNodeRestart(cluster: string) {
    return useMutationAdapter({
        mutationFn: NodeApi.restart.fn,
        mutationKey: NodeApi.restart.key(),
        successKeys: [ClusterApi.overview.key(cluster)]
    })
}

export function useRouterNodeReload(cluster: string) {
    return useMutationAdapter({
        mutationFn: NodeApi.reload.fn,
        mutationKey: NodeApi.reload.key(),
        successKeys: [ClusterApi.overview.key(cluster)]
    })
}

export function useRouterNodeReinit(cluster: string) {
    return useMutationAdapter({
        mutationFn: NodeApi.reinitialize.fn,
        mutationKey: NodeApi.reinitialize.key(),
        successKeys: [ClusterApi.overview.key(cluster)]
    })
}

export function useRouterNodeSwitchover(cluster: string) {
    return useMutationAdapter({
        mutationFn: NodeApi.switchover.fn,
        mutationKey: NodeApi.switchover.key(),
        successKeys: [ClusterApi.overview.key(cluster)]
    })
}

export function useRouterNodeFailover(cluster: string) {
    return useMutationAdapter({
        mutationFn: NodeApi.failover.fn,
        mutationKey: NodeApi.failover.key(),
        successKeys: [ClusterApi.overview.key(cluster)]
    })
}

export function useRouterNodeActivate(cluster: string) {
    return useMutationAdapter({
        mutationFn: NodeApi.activate.fn,
        mutationKey: NodeApi.activate.key(),
        successKeys: [ClusterApi.overview.key(cluster)]
    })
}

export function useRouterNodePause(cluster: string) {
    return useMutationAdapter({
        mutationFn: NodeApi.pause.fn,
        mutationKey: NodeApi.pause.key(),
        successKeys: [ClusterApi.overview.key(cluster)]
    })
}
