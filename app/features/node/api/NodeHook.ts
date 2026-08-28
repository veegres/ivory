import {useQuery} from "@tanstack/react-query"

import {useMutationAdapter} from "../../../shared/hook/QueryCustom"
import {useStream} from "../../../shared/hook/Stream"
import {useStore} from "../../../shared/provider/StoreProvider"
import {ClusterApi} from "../../cluster/api/ClusterRouter"
import {NodeApi} from "./NodeRouter"
import {
    KeeperOneRequest,
    KeeperPlugin,
    PlatformActionRequest,
    PlatformLogsRequest,
    PlatformVaultConnection,
} from "./NodeType"

export function useRouterNodePlatformLogs(request: PlatformLogsRequest, enabled: boolean) {
    const {loading, response, reconnect} = useStream(NodeApi.platform.logs.url(request), {enabled})
    return {isFetching: loading, data: response, reconnect}
}

export function useRouterNodePlatformContainerLogs(request: PlatformLogsRequest) {
    const {loading, response, reconnect} = useStream(NodeApi.container.logs.url(request))
    return {isFetching: loading, data: response, reconnect}
}

export function useRouterNodePlatformList(request: PlatformVaultConnection) {
    return useQuery({
        queryKey: NodeApi.container.list.key(request),
        queryFn: () => NodeApi.container.list.fn(request),
    })
}

export function useRouterNodeConfig(request?: KeeperOneRequest) {
    const req = request ?? {host: "", port: 1, plugin: KeeperPlugin.PATRONI_POSTGRES}
    // eslint-disable-next-line @tanstack/query/exhaustive-deps
    return useQuery({
        queryKey: NodeApi.keeper.config.key(req.host, req.port),
        queryFn: () => NodeApi.keeper.config.fn(req),
        enabled: !!request,
    })
}

export function useRouterNodeConfigUpdate(request: KeeperOneRequest | undefined, onSuccess: () => void) {
    const req = request ?? {host: "", port: 1, plugin: KeeperPlugin.PATRONI_POSTGRES}
    return useMutationAdapter({
        mutationFn: NodeApi.keeper.updateConfig.fn,
        mutationKey: NodeApi.keeper.updateConfig.key(),
        successKeys: [NodeApi.keeper.config.key(req.host, req.port)],
        onSuccess: onSuccess,
    })
}

export function useRouterNodeSwitchoverDelete(cluster: string) {
    return useMutationAdapter({
        mutationFn: NodeApi.keeper.deleteSwitchover.fn,
        mutationKey: NodeApi.keeper.deleteSwitchover.key(),
        successKeys: [ClusterApi.overview.keyCommon(cluster)]
    })
}

export function useRouterNodeRestartDelete(cluster: string) {
    return useMutationAdapter({
        mutationFn: NodeApi.keeper.deleteRestart.fn,
        mutationKey: NodeApi.keeper.deleteRestart.key(),
        successKeys: [ClusterApi.overview.keyCommon(cluster)]
    })
}

export function useRouterNodeRestart(cluster: string) {
    return useMutationAdapter({
        mutationFn: NodeApi.keeper.restart.fn,
        mutationKey: NodeApi.keeper.restart.key(),
        successKeys: [ClusterApi.overview.keyCommon(cluster)]
    })
}

export function useRouterNodeReload(cluster: string) {
    return useMutationAdapter({
        mutationFn: NodeApi.keeper.reload.fn,
        mutationKey: NodeApi.keeper.reload.key(),
        successKeys: [ClusterApi.overview.keyCommon(cluster)]
    })
}

export function useRouterNodeReinit(cluster: string) {
    return useMutationAdapter({
        mutationFn: NodeApi.keeper.reinitialize.fn,
        mutationKey: NodeApi.keeper.reinitialize.key(),
        successKeys: [ClusterApi.overview.keyCommon(cluster)]
    })
}

export function useRouterNodeSwitchover(cluster: string) {
    return useMutationAdapter({
        mutationFn: NodeApi.keeper.switchover.fn,
        mutationKey: NodeApi.keeper.switchover.key(),
        successKeys: [ClusterApi.overview.keyCommon(cluster)]
    })
}

export function useRouterNodeFailover(cluster: string) {
    return useMutationAdapter({
        mutationFn: NodeApi.keeper.failover.fn,
        mutationKey: NodeApi.keeper.failover.key(),
        successKeys: [ClusterApi.overview.keyCommon(cluster)]
    })
}

export function useRouterNodeActivate(cluster: string) {
    return useMutationAdapter({
        mutationFn: NodeApi.keeper.activate.fn,
        mutationKey: NodeApi.keeper.activate.key(),
        successKeys: [ClusterApi.overview.keyCommon(cluster)]
    })
}

export function useRouterNodePause(cluster: string) {
    return useMutationAdapter({
        mutationFn: NodeApi.keeper.pause.fn,
        mutationKey: NodeApi.keeper.pause.key(),
        successKeys: [ClusterApi.overview.keyCommon(cluster)]
    })
}

// NOTE: refetchInterval is intentionally not passed here - it is controlled by
//  <Refresher/> via queryClient.setQueryDefaults for this query key. Passing an
//  explicit value (even undefined) here would always win over that default.
export function useRouterNodeMetrics(c: PlatformVaultConnection) {
    // eslint-disable-next-line @tanstack/query/exhaustive-deps
    return useQuery({
        queryKey: NodeApi.platform.metrics.key(c.host),
        queryFn: () => NodeApi.platform.metrics.fn(c),
        retry: false,
    })
}

export function useRouterNodePlatformContainerMetrics(request: PlatformActionRequest) {
    return useQuery({
        queryKey: NodeApi.container.metrics.key(request),
        queryFn: () => NodeApi.container.metrics.fn(request),
        retry: false,
    })
}

export function useRouterNodePlatformProcesses(c: PlatformVaultConnection) {
    // eslint-disable-next-line @tanstack/query/exhaustive-deps
    return useQuery({
        queryKey: NodeApi.platform.processes.key(c.host),
        queryFn: () => NodeApi.platform.processes.fn(c),
        retry: false,
    })
}

export function useRouterNodePlatformInfo(c: PlatformVaultConnection) {
    // eslint-disable-next-line @tanstack/query/exhaustive-deps
    return useQuery({
        queryKey: NodeApi.platform.info.key(c.host),
        queryFn: () => NodeApi.platform.info.fn(c),
        retry: false,
    })
}

export function useRouterNodePlatformStart(connection: PlatformVaultConnection) {
    const activeCluster = useStore(s => s.activeCluster)
    const activeClusterKey = activeCluster ? ClusterApi.overview.keyCommon(activeCluster.name) : []
    return useMutationAdapter({
        mutationFn: NodeApi.container.start.fn,
        mutationKey: NodeApi.container.start.key(),
        successKeys: [NodeApi.container.list.key(connection), activeClusterKey],
    })
}

export function useRouterNodePlatformStop(connection: PlatformVaultConnection) {
    const activeCluster = useStore(s => s.activeCluster)
    const activeClusterKey = activeCluster ? ClusterApi.overview.keyCommon(activeCluster.name) : []
    return useMutationAdapter({
        mutationFn: NodeApi.container.stop.fn,
        mutationKey: NodeApi.container.stop.key(),
        successKeys: [NodeApi.container.list.key(connection), activeClusterKey],
    })
}

export function useRouterNodePlatformRestart(connection: PlatformVaultConnection) {
    const activeCluster = useStore(s => s.activeCluster)
    const activeClusterKey = activeCluster ? ClusterApi.overview.keyCommon(activeCluster.name) : []
    return useMutationAdapter({
        mutationFn: NodeApi.container.restart.fn,
        mutationKey: NodeApi.container.restart.key(),
        successKeys: [NodeApi.container.list.key(connection), activeClusterKey],
    })
}

export function useRouterNodePlatformDown(connection: PlatformVaultConnection) {
    const activeCluster = useStore(s => s.activeCluster)
    const activeClusterKey = activeCluster ? ClusterApi.overview.keyCommon(activeCluster.name) : []
    return useMutationAdapter({
        mutationFn: NodeApi.container.down.fn,
        mutationKey: NodeApi.container.down.key(),
        successKeys: [NodeApi.container.list.key(connection), activeClusterKey],
    })
}

export function useRouterNodeKeeperDeploy(connection: PlatformVaultConnection) {
    const activeCluster = useStore(s => s.activeCluster)
    const activeClusterKey = activeCluster ? ClusterApi.overview.keyCommon(activeCluster.name) : []
    return useMutationAdapter({
        mutationFn: NodeApi.container.keeper.deploy.fn,
        mutationKey: NodeApi.container.keeper.deploy.key(),
        successKeys: [NodeApi.container.list.key(connection), activeClusterKey],
    })
}

export function useRouterNodeKeeperDeploySpec(plugin: KeeperPlugin) {
    return useQuery({
        queryKey: NodeApi.container.keeper.deploySpec.key(plugin),
        queryFn: () => NodeApi.container.keeper.deploySpec.fn({plugin}),
    })
}
