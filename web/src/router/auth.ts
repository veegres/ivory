import {useMutationOptions} from "../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {GeneralApi} from "../app/api";

export function useRouterLogin(onSuccess: (data: any) => void) {
    const options = useMutationOptions([["info"]], onSuccess)
    return useMutation({mutationFn: GeneralApi.login, ...options})
}
