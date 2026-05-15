import {useQuery} from "@tanstack/react-query"

import {useMutationAdapter} from "../../hook/QueryCustom"
import {ClusterApi} from "../cluster/router"
import {NodeConfig} from "../cluster/type"
import {Plugin} from "../keeper/type"
import {NodeApi} from "./router"
import {CloudConnection,ContainerLogsRequest, KeeperRequest} from "./type"

export function useRouterNodeOverview(request: KeeperRequest, enabled: boolean) {
    return useQuery({
        // eslint-disable-next-line @tanstack/query/exhaustive-deps
        queryKey: NodeApi.overview.key(request.host, request.port),
        queryFn: () => NodeApi.overview.fn(request),
        enabled,
    })
}

export function useRouterNodeConfig(request?: KeeperRequest) {
    const req = request ?? {host: "", port: 1, plugin: Plugin.PATRONI}
    return useQuery({
        // eslint-disable-next-line @tanstack/query/exhaustive-deps
        queryKey: NodeApi.config.key(req.host, req.port),
        queryFn: () => NodeApi.config.fn(req),
        enabled: !!request,
    })
}

export function useRouterNodeConfigUpdate(config: NodeConfig, onSuccess: () => void) {
    return useMutationAdapter({
        mutationFn: NodeApi.updateConfig.fn,
        mutationKey: NodeApi.updateConfig.key(),
        successKeys: [NodeApi.config.key(config.host, config.keeperPort)],
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

export function useRouterNodeMetrics(c: CloudConnection, refetchInterval?: number) {
    return useQuery({
        // eslint-disable-next-line @tanstack/query/exhaustive-deps
        queryKey: NodeApi.metrics.key(c.host),
        queryFn: () => NodeApi.metrics.fn(c),
        refetchInterval,
        retry: false,
    })
}

export function useRouterNodeContainerList(c: CloudConnection, enabled: boolean) {
    return useQuery({
        // eslint-disable-next-line @tanstack/query/exhaustive-deps
        queryKey: NodeApi.container.list.key(c.host),
        queryFn: () => NodeApi.container.list.fn(c),
        enabled,
    })
}

export function useRouterNodeContainerLogs(request: ContainerLogsRequest, enabled: boolean) {
    return useQuery({
        // eslint-disable-next-line @tanstack/query/exhaustive-deps
        queryKey: NodeApi.container.logs.key(request.connection.host, request.container),
        queryFn: () => NodeApi.container.logs.fn(request),
        enabled,
    })
}

export function useRouterNodeContainerDeploy() {
    return useMutationAdapter({
        mutationFn: NodeApi.container.deploy.fn,
        mutationKey: NodeApi.container.deploy.key(),
    })
}

export function useRouterNodeContainerStop() {
    return useMutationAdapter({
        mutationFn: NodeApi.container.stop.fn,
        mutationKey: NodeApi.container.stop.key(),
    })
}

export function useRouterNodeContainerRun() {
    return useMutationAdapter({
        mutationFn: NodeApi.container.run.fn,
        mutationKey: NodeApi.container.run.key(),
    })
}

export function useRouterNodeContainerDelete() {
    return useMutationAdapter({
        mutationFn: NodeApi.container.delete.fn,
        mutationKey: NodeApi.container.delete.key(),
    })
}
