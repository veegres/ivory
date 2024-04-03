import {useQuery} from "@tanstack/react-query";
import {ClusterApi, TagApi} from "../app/api";
import {useMutationAdapter} from "../hook/QueryCustom";

export function useRouterClusterList(tags?: string[], subscribe: boolean = false) {
    return useQuery({
        queryKey: ClusterApi.list.key(),
        // NOTE: we need subscribe here, because we don't want to pass tags in Overview component
        // it will require to get tags from there to keep same behaviour. We want just subscribe by key.
        queryFn: subscribe ? undefined : () => ClusterApi.list.fn(tags),
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
