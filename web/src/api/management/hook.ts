import {useQuery} from "@tanstack/react-query"

import {useMutationAdapter} from "../../hook/QueryCustom"
import {ManagementApi} from "./router"

export function useRouterInfo(enabled: boolean = true) {
    return useQuery({
        queryKey: ManagementApi.info.key(),
        queryFn: () => ManagementApi.info.fn(),
        refetchOnWindowFocus: "always", enabled,
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
        mutationFn: () => ManagementApi.erase.fn("management"),
        mutationKey: ManagementApi.erase.key(),
        successKeys: [ManagementApi.info.key()],
        onSuccess: onSuccess,
    })
}

export function useRouterFree() {
    return useMutationAdapter({
        mutationFn: () => ManagementApi.free.fn(),
        mutationKey: ManagementApi.free.key(),
        successKeys: [ManagementApi.info.key()],
    })
}

