import {useQuery} from "@tanstack/react-query"

import {useMutationAdapter} from "../../hook/QueryCustom"
import {useStore} from "../../provider/StoreProvider"
import {TagApi} from "../tag/router"
import {ClusterApi} from "./router"

export function useRouterClusterList(tags: string[]) {
    const tagsFn = tags[0] === "ALL" ? undefined : tags
    // NOTE: this query is updated by custom logic with useEffect, without using queryKey change
    // we cannot add `enable: false`, because mutation hooks then couldn't update it by using QueryClient
    return useQuery({
        // eslint-disable-next-line @tanstack/query/exhaustive-deps
        queryKey: ClusterApi.list.key(),
        queryFn: () => ClusterApi.list.fn(tagsFn),
    })
}

export function useRouterClusterOverview(name?: string) {
    const activeCluster = useStore(s => s.activeCluster)
    const instance = activeCluster?.cluster.name === name ? activeCluster?.detectBy : undefined
    return useQuery({
        queryKey: ClusterApi.overview.key(name, instance?.sidecar),
        queryFn: () => ClusterApi.overview.fn(name ?? "disabled", instance?.sidecar),
        enabled: !!name, retry: false,
    })
}

export function useRouterClusterDelete() {
    return useMutationAdapter({
        mutationFn: ClusterApi.delete.fn,
        mutationKey: ClusterApi.delete.key(),
        successKeys: [ClusterApi.list.key(), TagApi.list.key()]
    })
}

export function useRouterClusterUpdate(onSuccess?: () => void) {
    return useMutationAdapter({
        mutationFn: ClusterApi.update.fn,
        mutationKey: ClusterApi.update.key(),
        successKeys: [ClusterApi.list.key(), TagApi.list.key()],
        onSuccess: onSuccess,
    })
}

export function useRouterClusterCreateAuto(onSuccess?: () => void) {
    return useMutationAdapter({
        mutationFn: ClusterApi.createAuto.fn,
        mutationKey: ClusterApi.createAuto.key(),
        successKeys: [ClusterApi.list.key(), TagApi.list.key()],
        onSuccess: onSuccess,
    })
}

export function useRouterClusterFixAuto(cluster: string) {
    return useMutationAdapter({
        mutationFn: ClusterApi.fixAuto.fn,
        mutationKey: ClusterApi.fixAuto.key(),
        successKeys: [ClusterApi.list.key(), ClusterApi.overview.key(cluster)],
    })
}
