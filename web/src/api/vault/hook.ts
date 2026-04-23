import {useQuery} from "@tanstack/react-query"

import {useMutationAdapter} from "../../hook/QueryCustom"
import {VaultApi} from "./router"
import {VaultType} from "./type"

export function useRouterVault(type?: VaultType) {
    return useQuery({
        queryKey: VaultApi.list.key(type),
        queryFn: () => VaultApi.list.fn(type)
    })
}

export function useRouterVaultDelete() {
    return useMutationAdapter({
        mutationFn: VaultApi.delete.fn,
        mutationKey: VaultApi.delete.key(),
        successKeys: [VaultApi.list.key()],
    })
}

export function useRouterVaultUpdate(onSuccess?: () => void) {
    return useMutationAdapter({
        mutationFn: VaultApi.update.fn,
        mutationKey: VaultApi.update.key(),
        successKeys: [VaultApi.list.key()],
        onSuccess: onSuccess,
    })
}

export function useRouterVaultCreate(onSuccess?: () => void) {
    return useMutationAdapter({
        mutationFn: VaultApi.create.fn,
        mutationKey: VaultApi.create.key(),
        successKeys: [VaultApi.list.key()],
        onSuccess: onSuccess,
    })
}
