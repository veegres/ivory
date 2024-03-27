import {useQuery} from "@tanstack/react-query";
import {ClusterApi} from "../app/api";

export function useRouterClusterList(tags?: string[]) {
    return useQuery({
        queryKey: ["cluster", "list"],
        queryFn: () => ClusterApi.list(tags)}
    )
}
