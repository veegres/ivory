import {useMutationOptions} from "../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {InitialApi} from "../app/api";

export function useRouterSecretSet() {
    const options = useMutationOptions([["info"]])
    return useMutation({mutationFn: InitialApi.setSecret, ...options})
}
