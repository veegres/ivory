import {useQuery} from "@tanstack/react-query"

import {useMutationAdapter} from "../../../shared/hook/QueryCustom"
import {useStore} from "../../../shared/provider/StoreProvider"
import {TagApi} from "../../tag/api/TagRouter"
import {ClusterApi} from "./ClusterRouter"

export function useRouterClusterList(enabled: boolean = true) {
    const tags = useStore(s => s.activeTags)
    const keeper = useStore(s => s.activeClusterKeeperPlugin)
    const tagsTmp = tags[0] === "ALL" ? undefined : tags
    return useQuery({
        queryKey: ClusterApi.list.key(tagsTmp, keeper),
        queryFn: () => ClusterApi.list.fn(tagsTmp, keeper),
        enabled: enabled,
    })
}

export function useRouterClusterOverview(name?: string, enabled: boolean = true) {
    const activeCluster = useStore(s => s.activeCluster)
    const manualKeeper = useStore(s => s.manualKeeper)
    const [host, port] = activeCluster?.name === name ? manualKeeper?.split(":") ?? [] : []
    return useQuery({
        queryKey: ClusterApi.overview.key(name ?? "disabled", host, port),
        queryFn: () => ClusterApi.overview.fn(name ?? "disabled", host, port),
        enabled: !!name && enabled, retry: false,
    })
}

export function useRouterClusterDelete() {
    return useMutationAdapter({
        mutationFn: ClusterApi.delete.fn,
        mutationKey: ClusterApi.delete.key(),
        successKeys: [ClusterApi.list.keyCommon(), TagApi.list.key()]
    })
}

export function useRouterClusterUpdate(name: string, onSuccess?: () => void) {
    return useMutationAdapter({
        mutationFn: ClusterApi.update.fn,
        mutationKey: ClusterApi.update.key(),
        successKeys: [ClusterApi.list.keyCommon(), TagApi.list.key(), ClusterApi.overview.keyCommon(name)],
        onSuccess: onSuccess,
    })
}

export function useRouterClusterCreate(onSuccess?: () => void) {
    return useMutationAdapter({
        mutationFn: ClusterApi.create.fn,
        mutationKey: ClusterApi.create.key(),
        successKeys: [ClusterApi.list.keyCommon(), TagApi.list.key()],
        onSuccess: onSuccess,
    })
}

export function useRouterClusterCreateAuto(onSuccess?: () => void) {
    return useMutationAdapter({
        mutationFn: ClusterApi.detect.fn,
        mutationKey: ClusterApi.detect.key(),
        successKeys: [ClusterApi.list.keyCommon(), TagApi.list.key()],
        onSuccess: onSuccess,
    })
}

export function useRouterClusterDeploy(onSuccess?: (data: string[]) => void) {
    return useMutationAdapter({
        mutationFn: ClusterApi.deploy.fn,
        mutationKey: ClusterApi.deploy.key(),
        successKeys: [ClusterApi.list.keyCommon(), TagApi.list.key()],
        onSuccess: (_, data) => onSuccess ? onSuccess(data) : void 0,
    })
}

export function useRouterClusterFix(name: string) {
    return useMutationAdapter({
        mutationFn: ClusterApi.fix.fn,
        mutationKey: ClusterApi.fix.key(),
        successKeys: [ClusterApi.list.keyCommon(), ClusterApi.overview.keyCommon(name)]
    })
}
