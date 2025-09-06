import {useMutationAdapter} from "../../hook/QueryCustom";
import {GeneralApi, InitialApi} from "../management/router";

export function useRouterSecretSet() {
    return useMutationAdapter({
        mutationFn: InitialApi.setSecret.fn,
        mutationKey: InitialApi.setSecret.key(),
        successKeys: [GeneralApi.info.key()],
    })
}
