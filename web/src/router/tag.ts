import {useQuery} from "@tanstack/react-query";
import {TagApi} from "../app/api";

export function useRouterTagList() {
    return useQuery({
        queryKey: ["tag", "list"],
        queryFn: TagApi.list
    })
}
