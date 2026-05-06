import {useQuery} from "@tanstack/react-query"

import {useMutationAdapter} from "../../hook/QueryCustom"
import {ClusterApi} from "../cluster/router"
import {NodeConfig} from "../cluster/type"
import {Plugin} from "../keeper/type"
import {NodeApi} from "./router"
import {DockerLogsRequest, KeeperRequest, SshConnection} from "./type"

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

export function useRouterNodeMetrics(c: SshConnection, refetchInterval?: number) {
    return useQuery({
        // eslint-disable-next-line @tanstack/query/exhaustive-deps
        queryKey: NodeApi.metrics.key(c.host),
        queryFn: () => NodeApi.metrics.fn(c),
        refetchInterval,
        retry: false,
    })
}

export function useRouterNodeDockerList(c: SshConnection, enabled: boolean) {
    return useQuery({
        // eslint-disable-next-line @tanstack/query/exhaustive-deps
        queryKey: NodeApi.docker.list.key(c.host),
        queryFn: () => NodeApi.docker.list.fn(c),
        enabled,
    })
}

export function useRouterNodeDockerLogs(request: DockerLogsRequest, enabled: boolean) {
    return useQuery({
        // eslint-disable-next-line @tanstack/query/exhaustive-deps
        queryKey: NodeApi.docker.logs.key(request.connection.host, request.container),
        queryFn: () => NodeApi.docker.logs.fn(request),
        enabled,
    })
}

export function useRouterNodeDockerDeploy() {
    return useMutationAdapter({
        mutationFn: NodeApi.docker.deploy.fn,
        mutationKey: NodeApi.docker.deploy.key(),
    })
}

export function useRouterNodeDockerStop() {
    return useMutationAdapter({
        mutationFn: NodeApi.docker.stop.fn,
        mutationKey: NodeApi.docker.stop.key(),
    })
}

export function useRouterNodeDockerRun() {
    return useMutationAdapter({
        mutationFn: NodeApi.docker.run.fn,
        mutationKey: NodeApi.docker.run.key(),
    })
}

export function useRouterNodeDockerDelete() {
    return useMutationAdapter({
        mutationFn: NodeApi.docker.delete.fn,
        mutationKey: NodeApi.docker.delete.key(),
    })
}
