import {useQuery} from "@tanstack/react-query"

import {useMutationAdapter} from "../../../shared/hook/QueryCustom"
import {useRouterClusterListKey} from "../../cluster/api/ClusterHook"
import {ManagementApi} from "./ManagementRouter"

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
        mutationFn: ManagementApi.free.fn,
        mutationKey: ManagementApi.free.key(),
    })
}

export function useRouterExport() {
    return useMutationAdapter({
        mutationFn: ManagementApi.export.fn,
        mutationKey: ManagementApi.export.key(),
    })
}

export function useRouterImport() {
    const clusterListKeys = useRouterClusterListKey()
    return useMutationAdapter({
        mutationFn: ManagementApi.import.fn,
        mutationKey: ManagementApi.import.key(),
        // NOTE: raw prefix invalidates all query lists regardless of type/plugin filters
        successKeys: [clusterListKeys, ["query", "list"]],
    })
}
