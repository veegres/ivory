import {useQuery} from "@tanstack/react-query"

import {useMutationAdapter} from "../../shared/hook/QueryCustom"
import {useStream} from "../../shared/hook/Stream"
import {useStore} from "../../shared/provider/StoreProvider"
import {ClusterApi} from "../cluster/router"
import {NodeConfig} from "../cluster/type"
import {Plugin} from "../keeper/type"
import {NodeApi} from "./router"
import {KeeperOneRequest, PlatformConnection, PlatformLogsRequest} from "./type"

export function useRouterNodePlatformLogs(request: PlatformLogsRequest) {
    const {loading, response} = useStream(NodeApi.deployment.logs.url(request))
    const data = response
        .filter(r => r.type === "log" || r.type === "server" || r.type === "browser")
        .map(r => r.message)
    return {isFetching: loading, data}
}

export function useRouterNodePlatformList(request: PlatformConnection) {
    return useQuery({
        queryKey: NodeApi.deployment.list.key(request),
        queryFn: () => NodeApi.deployment.list.fn(request),
    })
}

export function useRouterNodeConfig(request?: KeeperOneRequest) {
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

export function useRouterNodeMetrics(c: PlatformConnection, refetchInterval?: number) {
    return useQuery({
        // eslint-disable-next-line @tanstack/query/exhaustive-deps
        queryKey: NodeApi.metrics.key(c.host),
        queryFn: () => NodeApi.metrics.fn(c),
        refetchInterval,
        retry: false,
    })
}

export function useRouterNodePlatformStart(connection: PlatformConnection) {
    const activeCluster = useStore(s => s.activeCluster)
    const activeClusterKey = activeCluster ? ClusterApi.overview.key(activeCluster.name) : []
    return useMutationAdapter({
        mutationFn: NodeApi.deployment.start.fn,
        mutationKey: NodeApi.deployment.start.key(),
        successKeys: [NodeApi.deployment.list.key(connection), activeClusterKey],
    })
}

export function useRouterNodePlatformStop(connection: PlatformConnection) {
    const activeCluster = useStore(s => s.activeCluster)
    const activeClusterKey = activeCluster ? ClusterApi.overview.key(activeCluster.name) : []
    return useMutationAdapter({
        mutationFn: NodeApi.deployment.stop.fn,
        mutationKey: NodeApi.deployment.stop.key(),
        successKeys: [NodeApi.deployment.list.key(connection), activeClusterKey],
    })
}

export function useRouterNodePlatformDown(connection: PlatformConnection) {
    const activeCluster = useStore(s => s.activeCluster)
    const activeClusterKey = activeCluster ? ClusterApi.overview.key(activeCluster.name) : []
    return useMutationAdapter({
        mutationFn: NodeApi.deployment.down.fn,
        mutationKey: NodeApi.deployment.down.key(),
        successKeys: [NodeApi.deployment.list.key(connection), activeClusterKey],
    })
}
