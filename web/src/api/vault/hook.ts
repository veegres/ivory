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

export function useRouterVaultDelete(type?: VaultType) {
    return useMutationAdapter({
        mutationFn: VaultApi.delete.fn,
        mutationKey: VaultApi.delete.key(),
        successKeys: [VaultApi.list.key(type)],
    })
}

export function useRouterVaultUpdate(type?: VaultType, onSuccess?: () => void) {
    return useMutationAdapter({
        mutationFn: VaultApi.update.fn,
        mutationKey: VaultApi.update.key(),
        successKeys: [VaultApi.list.key(type)],
        onSuccess: onSuccess,
    })
}

export function useRouterVaultCreate(type?: VaultType, onSuccess?: () => void) {
    return useMutationAdapter({
        mutationFn: VaultApi.create.fn,
        mutationKey: VaultApi.create.key(),
        successKeys: [VaultApi.list.key(type)],
        onSuccess: onSuccess,
    })
}
