import {useMutationAdapter} from "../../hook/QueryCustom";
import {AuthType} from "../config/type";
import {AuthApi} from "./router";
import {ManagementApi} from "../management/router";

export function useRouterLogin(onSuccess: (data: any) => void, type: AuthType) {
    return useMutationAdapter({
        mutationFn: AuthApi.login.fn(type),
        mutationKey: AuthApi.login.key(type),
        successKeys: [ManagementApi.info.key()],
        onSuccess: (_, data) => onSuccess(data),
    })
}

export function useRouterLogout() {
    return useMutationAdapter({
        mutationFn: AuthApi.logout.fn,
        mutationKey: AuthApi.logout.key(),
        successKeys: [ManagementApi.info.key()],
    })
}
