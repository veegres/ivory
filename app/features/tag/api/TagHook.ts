import {useQuery} from "@tanstack/react-query"

import {TagApi} from "./TagRouter"

export function useRouterTagList() {
    return useQuery({
        queryKey: TagApi.list.key(),
        queryFn: () => TagApi.list.fn(),
    })
}
