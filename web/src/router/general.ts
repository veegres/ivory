import {useQuery} from "@tanstack/react-query";
import {GeneralApi, InitialApi, SafeApi} from "../app/api";
import {useMutationAdapter} from "../hook/QueryCustom";

export function useRouterInfo() {
    return useQuery({
        queryKey: GeneralApi.info.key(),
        queryFn: () => GeneralApi.info.fn(),
        refetchOnWindowFocus: "always",
    })
}

export function useRouterSecretChange() {
    return useMutationAdapter({
        mutationFn: SafeApi.changeSecret,
        successKeys: [GeneralApi.info.key()],
    })
}

export function useRouterEraseInitial(onSuccess: () => void) {
    return useMutationAdapter({
        mutationFn: InitialApi.erase,
        successKeys: [GeneralApi.info.key()],
        onSuccess: onSuccess,
    })
}

export function useRouterEraseSafe(onSuccess: () => void) {
    return useMutationAdapter({
        mutationFn: SafeApi.erase,
        successKeys: [GeneralApi.info.key()],
        onSuccess: onSuccess,
    })
}

export function useRouterConfigSet() {
    return useMutationAdapter({
        mutationFn: GeneralApi.setConfig,
        successKeys: [GeneralApi.info.key()],
    })
}
