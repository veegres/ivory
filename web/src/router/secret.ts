import {useMutationAdapter} from "../hook/QueryCustom";
import {GeneralApi, InitialApi} from "../app/api";

export function useRouterSecretSet() {
    return useMutationAdapter({
        mutationFn: InitialApi.setSecret.fn,
        mutationKey: InitialApi.setSecret.key(),
        successKeys: [GeneralApi.info.key()],
    })
}
