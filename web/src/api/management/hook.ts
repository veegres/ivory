import {useQuery} from "@tanstack/react-query";
import {useMutationAdapter} from "../../hook/QueryCustom";
import {ManagementApi} from "./router";

export function useRouterInfo() {
    return useQuery({
        queryKey: ManagementApi.info.key(),
        queryFn: () => ManagementApi.info.fn(),
        refetchOnWindowFocus: "always",
    })
}

export function useRouterSecretChange() {
    return useMutationAdapter({
        mutationFn: ManagementApi.changeSecret.fn,
        mutationKey: ManagementApi.changeSecret.key(),
        successKeys: [ManagementApi.info.key()],
    })
}

export function useRouterEraseInitial(onSuccess: () => void) {
    return useMutationAdapter({
        mutationFn: () => ManagementApi.erase.fn("initial"),
        mutationKey: ManagementApi.erase.key(),
        successKeys: [ManagementApi.info.key()],
        onSuccess: onSuccess,
    })
}

export function useRouterEraseSafe(onSuccess: () => void) {
    return useMutationAdapter({
        mutationFn: () => ManagementApi.erase.fn("safe"),
        mutationKey: ManagementApi.erase.key(),
        successKeys: [ManagementApi.info.key()],
        onSuccess: onSuccess,
    })
}

