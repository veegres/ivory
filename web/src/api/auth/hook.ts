import {useMutationAdapter} from "../../hook/QueryCustom";
import {AuthApi} from "./router";
import {ManagementApi} from "../management/router";
import {AuthType} from "./type";

export function useRouterLogin(type: AuthType) {
    return useMutationAdapter({
        mutationFn: AuthApi.login.fn(type),
        mutationKey: AuthApi.login.key(type),
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
