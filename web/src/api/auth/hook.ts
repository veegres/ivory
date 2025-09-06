import {useMutationAdapter} from "../../hook/QueryCustom";
import {GeneralApi} from "../management/router";

export function useRouterLogin(onSuccess: (data: any) => void) {
    return useMutationAdapter({
        mutationFn: GeneralApi.login.fn,
        mutationKey: GeneralApi.login.key(),
        successKeys: [GeneralApi.info.key()],
        onSuccess: (_, data) => onSuccess(data),
    })
}
