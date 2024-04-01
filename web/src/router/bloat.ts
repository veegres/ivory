import {BloatApi} from "../app/api";
import {useQuery} from "@tanstack/react-query";
import {useMutationAdapter} from "../hook/QueryCustom";
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
    return useMutationAdapter({
        mutationFn: BloatApi.start.fn,
        mutationKey: BloatApi.start.key(),
        onSuccess: (client, job) => {
            client.setQueryData<Bloat[]>(
                BloatApi.list.key(cluster),
                (jobs) => [job, ...(jobs ?? [])]
            )
        },
    })
}

export function useRouterBloatDelete(uuid: string, cluster: string) {
    return useMutationAdapter({
        mutationFn: BloatApi.delete.fn,
        mutationKey: BloatApi.delete.key(),
        onSuccess: (client) => {
            client.setQueryData<Bloat[]>(
                BloatApi.list.key(cluster),
                (jobs) => jobs?.filter(v => v.uuid !== uuid)
            )
        }
    })
}

export function useRouterBloatStop() {
    return useMutationAdapter({mutationFn: BloatApi.stop.fn})
}
