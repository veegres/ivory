import {useMutation, useQuery} from "@tanstack/react-query";
import {PasswordApi} from "../app/api";
import {PasswordType} from "../type/password";
import {useMutationOptions} from "../hook/QueryCustom";

export function useRouterPassword(type?: PasswordType) {
    return useQuery({
        queryKey: PasswordApi.list.key(type),
        queryFn: () => PasswordApi.list.fn(type)
    })
}

export function useRouterPasswordDelete() {
    const options = useMutationOptions([PasswordApi.list.key()])
    return useMutation({mutationFn: PasswordApi.delete, ...options})
}

export function useRouterPasswordUpdate(onSuccess?: () => void) {
    const options = useMutationOptions([PasswordApi.list.key()], onSuccess)
    return useMutation({mutationFn: PasswordApi.update, ...options})
}

export function useRouterPasswordCreate(onSuccess?: () => void) {
    const options = useMutationOptions([PasswordApi.list.key()], onSuccess)
    return useMutation({mutationFn: PasswordApi.create, ...options})
}
