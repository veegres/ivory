import {useMutation, useQuery} from "@tanstack/react-query";
import {GeneralApi, InitialApi, SafeApi} from "../app/api";
import {useMutationOptions} from "../hook/QueryCustom";

export function useRouterInfo() {
    return useQuery({
        queryKey: GeneralApi.info.key(),
        queryFn: () => GeneralApi.info.fn(),
        refetchOnWindowFocus: "always",
    })
}

export function useRouterSecretChange() {
    const options = useMutationOptions([GeneralApi.info.key()])
    return useMutation({mutationFn: SafeApi.changeSecret, ...options})
}

export function useRouterEraseInitial(onSuccess: () => void) {
    const options = useMutationOptions([GeneralApi.info.key()], onSuccess)
    return useMutation({mutationFn: InitialApi.erase, ...options})
}

export function useRouterEraseSafe(onSuccess: () => void) {
    const options = useMutationOptions([GeneralApi.info.key()], onSuccess)
    return useMutation({mutationFn: SafeApi.erase, ...options})
}

export function useRouterConfigSet() {
    const options = useMutationOptions([GeneralApi.info.key()])
    return useMutation({mutationFn: GeneralApi.setConfig, ...options})
}
