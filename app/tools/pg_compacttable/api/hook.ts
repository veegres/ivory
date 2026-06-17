import {useQuery} from "@tanstack/react-query"

import {useMutationAdapter} from "../../../shared/hook/QueryCustom"
import {PgCompactTableApi} from "./router"
import {PgCompactTable} from "./type"

export function useRouterPgCompactTableLogs(uuid: string, enabled: boolean) {
    return useQuery({
        queryKey: PgCompactTableApi.logs.key(uuid),
        queryFn: () => PgCompactTableApi.logs.fn(uuid),
        refetchOnReconnect: false,
        refetchOnWindowFocus: false,
        enabled,
    })
}

export function useRouterPgCompactTableList(cluster: string, enabled: boolean) {
    return useQuery({
        initialData: [],
        queryKey: PgCompactTableApi.list.key(cluster),
        queryFn: () => PgCompactTableApi.list.fn(cluster),
        enabled,
    })
}

export function useRouterPgCompactTableStart(cluster: string) {
    return useMutationAdapter({
        mutationFn: PgCompactTableApi.start.fn,
        mutationKey: PgCompactTableApi.start.key(),
        onSuccess: (client, job) => {
            client.setQueryData<PgCompactTable[]>(
                PgCompactTableApi.list.key(cluster),
                (jobs) => [job, ...(jobs ?? [])]
            )
        },
    })
}

export function useRouterPgCompactTableDelete(uuid: string, cluster: string) {
    return useMutationAdapter({
        mutationFn: PgCompactTableApi.delete.fn,
        mutationKey: PgCompactTableApi.delete.key(),
        onSuccess: (client) => {
            client.setQueryData<PgCompactTable[]>(
                PgCompactTableApi.list.key(cluster),
                (jobs) => jobs?.filter(v => v.uuid !== uuid)
            )
        }
    })
}

export function useRouterPgCompactTableStop() {
    return useMutationAdapter({mutationFn: PgCompactTableApi.stop.fn})
}
