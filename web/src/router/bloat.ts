import {BloatApi} from "../app/api";
import {useMutation, useQuery} from "@tanstack/react-query";
import {useMutationOptions} from "../hook/QueryCustom";
import {Bloat} from "../type/bloat";

export function useRouterBloatLogs(uuid: string, enabled: boolean) {
    return useQuery({
        queryKey: BloatApi.logs.key(uuid),
        queryFn: () => BloatApi.logs.fn(uuid),
        refetchOnReconnect: false,
        refetchOnWindowFocus: false,
        enabled,
    })
}

export function useRouterBloatList(cluster: string, enabled: boolean) {
    return useQuery({
        initialData: [],
        queryKey: BloatApi.list.key(cluster),
        queryFn: () => BloatApi.list.fn(cluster),
        enabled,
    })
}

export function useRouterBloatStart(cluster: string) {
    const options = useMutationOptions()
    return useMutation({
        mutationFn: BloatApi.start,
        ...options,
        onSuccess: (job) => {
            options.queryClient.setQueryData<Bloat[]>(
                BloatApi.list.key(cluster),
                (jobs) => [job, ...(jobs ?? [])]
            )
        },
    })
}

export function useRouterBloatDelete(uuid: string, cluster: string) {
    const options = useMutationOptions()
    return useMutation({
        mutationFn: BloatApi.delete,
        ...options,
        onSuccess: () => {
            options.queryClient.setQueryData<Bloat[]>(
                BloatApi.list.key(cluster),
                (jobs) => jobs?.filter(v => v.uuid !== uuid)
            )
        }
    })
}

export function useRouterBloatStop() {
    const options = useMutationOptions()
    return useMutation({mutationFn: BloatApi.stop, ...options})
}
