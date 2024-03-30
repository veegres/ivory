import {useMutationAdapter} from "../hook/QueryCustom";
import {GeneralApi} from "../app/api";

export function useRouterLogin(onSuccess: (data: any) => void) {
    return useMutationAdapter({
        mutationFn: GeneralApi.login,
        successKeys: [GeneralApi.info.key()],
        onSuccess: (_, data) => onSuccess(data),
    })
}
