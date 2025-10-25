import {useQuery} from "@tanstack/react-query";

import {useMutationAdapter} from "../../hook/QueryCustom";
import {TagApi} from "../tag/router";
import {ClusterApi} from "./router";

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
