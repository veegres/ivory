import {useMutationAdapter} from "../../hook/QueryCustom"
import {ManagementApi} from "../management/router"
import {AuthApi} from "./router"

export function useRouterLogin() {
    return useMutationAdapter({
        mutationFn: AuthApi.login.fn,
        mutationKey: AuthApi.login.key(),
        successKeys: [ManagementApi.info.key()],
    })
}

export function useRouterLogout() {
    return useMutationAdapter({
        mutationFn: AuthApi.logout.fn,
        mutationKey: AuthApi.logout.key(),
        successKeys: [ManagementApi.info.key()],
    })
}

export function useRouterConnect() {
    return useMutationAdapter({
        mutationFn: AuthApi.connect.fn,
        mutationKey: AuthApi.connect.key(),
    })
}
