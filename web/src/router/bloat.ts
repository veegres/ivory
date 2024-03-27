import {BloatApi} from "../app/api";
import {useQuery} from "@tanstack/react-query";

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
