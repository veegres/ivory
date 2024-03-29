import {useMutation, useQuery} from "@tanstack/react-query";
import {ClusterApi} from "../app/api";
import {useMutationOptions} from "../hook/QueryCustom";

export function useRouterClusterList(tags?: string[]) {
    return useQuery({
        queryKey: ["cluster", "list"],
        queryFn: () => ClusterApi.list(tags)}
    )
}

export function useRouterClusterDelete() {
    const options = useMutationOptions([["cluster", "list"], ["tag", "list"]])
    return useMutation({mutationFn: ClusterApi.delete, ...options})
}

export function useRouterClusterUpdate(onSuccess?: () => void) {
    const options = useMutationOptions([["cluster", "list"], ["tag", "list"]], onSuccess)
    return useMutation({mutationFn: ClusterApi.update, ...options})
}

export function useRouterClusterCreateAuto(onSuccess?: () => void) {
    const options = useMutationOptions([["cluster", "list"], ["tag", "list"]], onSuccess)
    return useMutation({mutationFn: ClusterApi.createAuto, ...options})
}
