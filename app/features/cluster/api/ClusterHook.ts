import {useQuery} from "@tanstack/react-query"
import {useEffect} from "react"

import {useMutationAdapter} from "../../../shared/hook/QueryCustom"
import {useStore} from "../../../shared/provider/StoreProvider"
import {TagApi} from "../../tag/api/TagRouter"
import {ClusterApi} from "./ClusterRouter"

export function useRouterClusterList(enabled: boolean = true) {
    const tags = useStore(s => s.activeTags)
    const keeper = useStore(s => s.activeClusterKeeperPlugin)
    const database = undefined
    const tagsTmp = tags[0] === "ALL" ? undefined : tags
    // NOTE: this query is updated by custom logic with useEffect, without using queryKey change
    // we cannot add `enable: false`, because mutation hooks then couldn't update it by using QueryClient
    const response = useQuery({
        // eslint-disable-next-line @tanstack/query/exhaustive-deps
        queryKey: ClusterApi.list.key(),
        queryFn: () => ClusterApi.list.fn(tagsTmp, keeper, database),
        enabled: enabled,
    })

    // TODO this request executes twice when we change keeper type, we need to remove hack and make proper loading logic
    // NOTE: we don't use queryKey to update it, because it will create a separate request and cause new fetching
    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(() => {response.refetch().then()}, [tags, keeper])

    return response
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
    return useMutationAdapter({
        mutationFn: ClusterApi.delete.fn,
        mutationKey: ClusterApi.delete.key(),
        successKeys: [ClusterApi.list.key(), TagApi.list.key()]
    })
}

export function useRouterClusterUpdate(name: string, onSuccess?: () => void) {
    return useMutationAdapter({
        mutationFn: ClusterApi.update.fn,
        mutationKey: ClusterApi.update.key(),
        successKeys: [ClusterApi.list.key(), TagApi.list.key(), ClusterApi.overview.key(name)],
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

export function useRouterClusterDeploy(onSuccess?: (data: string[]) => void) {
    return useMutationAdapter({
        mutationFn: ClusterApi.deploy.fn,
        mutationKey: ClusterApi.deploy.key(),
        successKeys: [ClusterApi.list.key(), TagApi.list.key()],
        onSuccess: (_, data) => onSuccess ? onSuccess(data) : void 0,
    })
}

export function useRouterClusterFixAuto(name: string) {
    return useMutationAdapter({
        mutationFn: ClusterApi.fixAuto.fn,
        mutationKey: ClusterApi.fixAuto.key(),
        successKeys: [ClusterApi.list.key(), ClusterApi.overview.key(name)]
    })
}
