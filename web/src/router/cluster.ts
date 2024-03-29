import {useMutation, useQuery} from "@tanstack/react-query";
import {ClusterApi, TagApi} from "../app/api";
import {useMutationOptions} from "../hook/QueryCustom";

export function useRouterClusterList(tags?: string[]) {
    return useQuery({
        queryKey: ClusterApi.list.key(),
        queryFn: () => ClusterApi.list.fn(tags)}
    )
}

export function useRouterClusterDelete() {
    const options = useMutationOptions([ClusterApi.list.key(), TagApi.list.key()])
    return useMutation({mutationFn: ClusterApi.delete, ...options})
}

export function useRouterClusterUpdate(onSuccess?: () => void) {
    const options = useMutationOptions([ClusterApi.list.key(), TagApi.list.key()], onSuccess)
    return useMutation({mutationFn: ClusterApi.update, ...options})
}

export function useRouterClusterCreateAuto(onSuccess?: () => void) {
    const options = useMutationOptions([ClusterApi.list.key(), TagApi.list.key()], onSuccess)
    return useMutation({mutationFn: ClusterApi.createAuto, ...options})
}
