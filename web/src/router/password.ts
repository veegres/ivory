import {useMutation, useQuery} from "@tanstack/react-query";
import {PasswordApi} from "../app/api";
import {PasswordType} from "../type/password";
import {useMutationOptions} from "../hook/QueryCustom";

export function useRouterPassword(type?: PasswordType) {
    return useQuery({
        queryKey: ["password", type],
        queryFn: () => PasswordApi.list(type)
    })
}

export function useRouterPasswordDelete() {
    const options = useMutationOptions([["password"]])
    return useMutation({mutationFn: PasswordApi.delete, ...options})
}

export function useRouterPasswordUpdate(onSuccess?: () => void) {
    const options = useMutationOptions([["password"]], onSuccess)
    return useMutation({mutationFn: PasswordApi.update, ...options})
}

export function useRouterPasswordCreate(onSuccess?: () => void) {
    const options = useMutationOptions([["password"]], onSuccess)
    return useMutation({mutationFn: PasswordApi.create, ...options})
}
