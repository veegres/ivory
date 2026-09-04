import {useQuery} from "@tanstack/react-query"

import {useMutationAdapter} from "../../../shared/hook/QueryCustom"
import {ClusterApi} from "../../cluster/api/ClusterRouter"
import {QueryApi} from "../../query/api/QueryRouter"
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

// NOTE: erasing is only ever asked for by somebody signed in - whoever has lost
// the secret word reinstalls instead, which is why there is no public erase any
// more
export function useRouterErase(onSuccess: () => void) {
    return useMutationAdapter({
        mutationFn: ManagementApi.erase.fn,
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
    return useMutationAdapter({
        mutationFn: ManagementApi.import.fn,
        mutationKey: ManagementApi.import.key(),
        successKeys: [ClusterApi.list.keyCommon(), QueryApi.list.keyCommon()],
    })
}
