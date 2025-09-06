import {useQuery} from "@tanstack/react-query";
import {useMutationAdapter} from "../../hook/QueryCustom";
import {GeneralApi, InitialApi, SafeApi} from "./router";


export function useRouterInfo() {
    return useQuery({
        queryKey: GeneralApi.info.key(),
        queryFn: () => GeneralApi.info.fn(),
        refetchOnWindowFocus: "always",
    })
}

export function useRouterSecretChange() {
    return useMutationAdapter({
        mutationFn: SafeApi.changeSecret.fn,
        mutationKey: SafeApi.changeSecret.key(),
        successKeys: [GeneralApi.info.key()],
    })
}

export function useRouterEraseInitial(onSuccess: () => void) {
    return useMutationAdapter({
        mutationFn: InitialApi.erase.fn,
        mutationKey: InitialApi.erase.key(),
        successKeys: [GeneralApi.info.key()],
        onSuccess: onSuccess,
    })
}

export function useRouterEraseSafe(onSuccess: () => void) {
    return useMutationAdapter({
        mutationFn: SafeApi.erase.fn,
        mutationKey: SafeApi.erase.key(),
        successKeys: [GeneralApi.info.key()],
        onSuccess: onSuccess,
    })
}

export function useRouterConfigSet() {
    return useMutationAdapter({
        mutationFn: GeneralApi.setConfig.fn,
        mutationKey: GeneralApi.setConfig.key(),
        successKeys: [GeneralApi.info.key()],
    })
}
