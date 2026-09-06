import {useMutationAdapter} from "../../../shared/hook/QueryCustom"
import {ManagementApi} from "../../management/api/ManagementRouter"
import {UserApi} from "../../user/api/UserRouter"
import {AuthApi} from "./AuthRouter"

export function useRouterLogin(onSuccess?: () => void) {
    return useMutationAdapter({
        mutationFn: AuthApi.login.fn,
        mutationKey: AuthApi.login.key(),
        successKeys: [ManagementApi.info.key()],
        onSuccess: () => onSuccess?.(),
    })
}

export function useRouterLogout() {
    return useMutationAdapter({
        mutationFn: AuthApi.logout.fn,
        mutationKey: AuthApi.logout.key(),
        successKeys: [ManagementApi.info.key(), UserApi.registration.verify.keyCommon()],
    })
}

export function useRouterConnect() {
    return useMutationAdapter({
        mutationFn: AuthApi.connect.fn,
        mutationKey: AuthApi.connect.key(),
    })
}
