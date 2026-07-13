import {useQuery} from "@tanstack/react-query"

import {useMutationAdapter} from "../../../shared/hook/QueryCustom"
import {useStore} from "../../../shared/provider/StoreProvider"
import {KeeperPlugin} from "../../node/api/NodeType"
import {TagApi} from "../../tag/api/TagRouter"
import {ClusterApi} from "./ClusterRouter"

function useRouterClusterListSearchParam(): [KeeperPlugin, string[] | undefined] {
    const tags = useStore(s => s.activeTags)
    const keeper = useStore(s => s.activeClusterKeeperPlugin)
    const tagsTmp = tags[0] === "ALL" ? undefined : tags
    return [keeper, tagsTmp]
}

export function useRouterClusterListKey() {
    const [keeper, tags] = useRouterClusterListSearchParam()
    return ClusterApi.list.key(tags, keeper)
}

export function useRouterClusterList(enabled: boolean = true) {
    const [keeper, tags] = useRouterClusterListSearchParam()
    return useQuery({
        queryKey: ClusterApi.list.key(tags, keeper),
        queryFn: () => ClusterApi.list.fn(tags, keeper),
        enabled: enabled,
    })
}

export function useRouterClusterOverview(name?: string, enabled: boolean = true) {
    const activeCluster = useStore(s => s.activeCluster)
    const manualKeeper = useStore(s => s.manualKeeper)
    const [host, port] = activeCluster?.name === name ? manualKeeper?.split(":") ?? [] : []
    return useQuery({
        queryKey: ClusterApi.overview.key(name, host, port),
        queryFn: () => ClusterApi.overview.fn(name ?? "disabled", host, port),
        enabled: !!name && enabled, retry: false,
    })
}

export function useRouterClusterDelete() {
    const clusterListKeys = useRouterClusterListKey()
    return useMutationAdapter({
        mutationFn: ClusterApi.delete.fn,
        mutationKey: ClusterApi.delete.key(),
        successKeys: [clusterListKeys, TagApi.list.key()]
    })
}

export function useRouterClusterUpdate(name: string, onSuccess?: () => void) {
    const clusterListKeys = useRouterClusterListKey()
    return useMutationAdapter({
        mutationFn: ClusterApi.update.fn,
        mutationKey: ClusterApi.update.key(),
        successKeys: [clusterListKeys, TagApi.list.key(), ClusterApi.overview.key(name)],
        onSuccess: onSuccess,
    })
}

export function useRouterClusterCreateAuto(onSuccess?: () => void) {
    const clusterListKeys = useRouterClusterListKey()
    return useMutationAdapter({
        mutationFn: ClusterApi.createAuto.fn,
        mutationKey: ClusterApi.createAuto.key(),
        successKeys: [clusterListKeys, TagApi.list.key()],
        onSuccess: onSuccess,
    })
}

export function useRouterClusterDeploy(onSuccess?: (data: string[]) => void) {
    const clusterListKeys = useRouterClusterListKey()
    return useMutationAdapter({
        mutationFn: ClusterApi.deploy.fn,
        mutationKey: ClusterApi.deploy.key(),
        successKeys: [clusterListKeys, TagApi.list.key()],
        onSuccess: (_, data) => onSuccess ? onSuccess(data) : void 0,
    })
}

export function useRouterClusterFixAuto(name: string) {
    const clusterListKeys = useRouterClusterListKey()
    return useMutationAdapter({
        mutationFn: ClusterApi.fixAuto.fn,
        mutationKey: ClusterApi.fixAuto.key(),
        successKeys: [clusterListKeys, ClusterApi.overview.key(name)]
    })
}
