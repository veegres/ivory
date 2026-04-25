import {useQuery} from "@tanstack/react-query"

import {useMutationAdapter} from "../../hook/QueryCustom"
import {ClusterApi} from "../cluster/router"
import {NodeApi} from "./router"
import {Connection, DockerLogsRequest, NodeRequest, SshRequest} from "./type"

export function useRouterNodeOverview(request: NodeRequest, enabled: boolean) {
    return useQuery({
        // eslint-disable-next-line @tanstack/query/exhaustive-deps
        queryKey: NodeApi.overview.key(request.connection),
        queryFn: () => NodeApi.overview.fn(request),
        enabled,
    })
}

export function useRouterNodeConfig(request: NodeRequest, enabled: boolean) {
    return useQuery({
        // eslint-disable-next-line @tanstack/query/exhaustive-deps
        queryKey: NodeApi.config.key(request.connection),
        queryFn: () => NodeApi.config.fn(request),
        enabled,
    })
}

export function useRouterNodeConfigUpdate(connection: Connection, onSuccess: () => void) {
    return useMutationAdapter({
        mutationFn: NodeApi.updateConfig.fn,
        mutationKey: NodeApi.updateConfig.key(),
        successKeys: [NodeApi.config.key(connection)],
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

export function useRouterNodeMetrics(request: SshRequest, enabled: boolean) {
    return useQuery({
        // eslint-disable-next-line @tanstack/query/exhaustive-deps
        queryKey: NodeApi.metrics.key(request.connection),
        queryFn: () => NodeApi.metrics.fn(request),
        enabled,
    })
}

export function useRouterNodeDockerList(request: SshRequest, enabled: boolean) {
    return useQuery({
        // eslint-disable-next-line @tanstack/query/exhaustive-deps
        queryKey: NodeApi.docker.list.key(request.connection),
        queryFn: () => NodeApi.docker.list.fn(request),
        enabled,
    })
}

export function useRouterNodeDockerLogs(request: DockerLogsRequest, enabled: boolean) {
    return useQuery({
        // eslint-disable-next-line @tanstack/query/exhaustive-deps
        queryKey: NodeApi.docker.logs.key(request.connection, request.container),
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
