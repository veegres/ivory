import {useQuery} from "@tanstack/react-query";
import {PasswordApi} from "../app/api";
import {PasswordType} from "../type/password";
import {useMutationAdapter} from "../hook/QueryCustom";

export function useRouterPassword(type?: PasswordType) {
    return useQuery({
        queryKey: PasswordApi.list.key(type),
        queryFn: () => PasswordApi.list.fn(type)
    })
}

export function useRouterPasswordDelete() {
    return useMutationAdapter({
        mutationFn: PasswordApi.delete,
        successKeys: [PasswordApi.list.key()],
    })
}

export function useRouterPasswordUpdate(onSuccess?: () => void) {
    return useMutationAdapter({
        mutationFn: PasswordApi.update,
        successKeys: [PasswordApi.list.key()],
        onSuccess: onSuccess,
    })
}

export function useRouterPasswordCreate(onSuccess?: () => void) {
    return useMutationAdapter({
        mutationFn: PasswordApi.create,
        successKeys: [PasswordApi.list.key()],
        onSuccess: onSuccess,
    })
}
