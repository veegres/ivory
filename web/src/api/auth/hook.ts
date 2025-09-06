import {useMutationAdapter} from "../../hook/QueryCustom";
import {ManagementApi} from "../management/router";

export function useRouterLogin(onSuccess: (data: any) => void) {
    return useMutationAdapter({
        mutationFn: ManagementApi.login.fn,
        mutationKey: ManagementApi.login.key(),
        successKeys: [ManagementApi.info.key()],
        onSuccess: (_, data) => onSuccess(data),
    })
}
