import {useQuery} from "@tanstack/react-query"

import {useMutationAdapter} from "../../hook/QueryCustom"
import {ManagementApi} from "../management/router"
import {PermissionApi} from "./router"

export function useRouterPermissions() {
    return useQuery({
        queryKey: PermissionApi.list.key(),
        queryFn: () => PermissionApi.list.fn()
    })
}

export function useRouterPermissionRequest() {
    return useMutationAdapter({
        mutationFn: PermissionApi.request.fn,
        mutationKey: PermissionApi.request.key(),
        successKeys: [ManagementApi.info.key()],
    })
}

export function useRouterPermissionApprove() {
    return useMutationAdapter({
        mutationFn: PermissionApi.approve.fn,
        mutationKey: PermissionApi.approve.key(),
        successKeys: [PermissionApi.list.key()],
    })
}

export function useRouterPermissionReject() {
    return useMutationAdapter({
        mutationFn: PermissionApi.reject.fn,
        mutationKey: PermissionApi.reject.key(),
        successKeys: [PermissionApi.list.key()],
    })
}

export function useRouterPermissionDelete() {
    return useMutationAdapter({
        mutationFn: PermissionApi.delete.fn,
        mutationKey: PermissionApi.delete.key(),
        successKeys: [PermissionApi.list.key()],
    })
}
