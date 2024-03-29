import {BloatApi} from "../app/api";
import {useMutation, useQuery} from "@tanstack/react-query";
import {useMutationOptions} from "../hook/QueryCustom";
import {Bloat} from "../type/bloat";

export function useRouterBloatLogs(uuid: string, enabled: boolean) {
    return useQuery({
        queryKey: ["instance", "bloat", "logs", uuid],
        queryFn: () => BloatApi.logs(uuid),
        refetchOnReconnect: false,
        refetchOnWindowFocus: false,
        enabled,
    })
}

export function useRouterBloatList(cluster: string, enabled: boolean) {
    return useQuery({
        initialData: [],
        queryKey: ["instance", "bloat", "list", cluster],
        queryFn: () => BloatApi.list(cluster),
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
                ["instance", "bloat", "list", cluster],
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
                ["instance", "bloat", "list", cluster],
                (jobs) => jobs?.filter(v => v.uuid !== uuid)
            )
        }
    })
}

export function useRouterBloatStop() {
    const options = useMutationOptions()
    return useMutation({mutationFn: BloatApi.stop, ...options})
}
