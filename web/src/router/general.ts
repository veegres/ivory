import {useQuery} from "@tanstack/react-query";
import {GeneralApi} from "../app/api";

export function useRouterInfo() {
    return useQuery({
        queryKey: ["info"],
        queryFn: GeneralApi.info,
        refetchOnWindowFocus: "always",
    })
}
