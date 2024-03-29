import {useMutationOptions} from "../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {GeneralApi, InitialApi} from "../app/api";

export function useRouterSecretSet() {
    const options = useMutationOptions([GeneralApi.info.key()])
    return useMutation({mutationFn: InitialApi.setSecret, ...options})
}
