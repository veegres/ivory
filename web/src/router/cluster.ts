import {useQuery} from "@tanstack/react-query";
import {ClusterApi, TagApi} from "../app/api";
import {useMutationAdapter} from "../hook/QueryCustom";

export function useRouterClusterList(tags?: string[]) {
    return useQuery({
        queryKey: ClusterApi.list.key(),
        queryFn: () => ClusterApi.list.fn(tags)}
    )
}

export function useRouterClusterDelete() {
    return useMutationAdapter({
        mutationFn: ClusterApi.delete.fn,
        successKeys: [ClusterApi.list.key(), TagApi.list.key()]
    })
}

export function useRouterClusterUpdate(onSuccess?: () => void) {
    return useMutationAdapter({
        mutationFn: ClusterApi.update.fn,
        successKeys: [ClusterApi.list.key(), TagApi.list.key()],
        onSuccess: onSuccess,
    })
}

export function useRouterClusterCreateAuto(onSuccess?: () => void) {
    return useMutationAdapter({
        mutationFn: ClusterApi.createAuto.fn,
        successKeys: [ClusterApi.list.key(), TagApi.list.key()],
        onSuccess: onSuccess,
    })
}