import {useMutationAdapter} from "../../hook/QueryCustom"
import {ManagementApi} from "../management/router"
import {SecretApi} from "./router"

export function useRouterSecretSet() {
    return useMutationAdapter({
        mutationFn: SecretApi.setSecret.fn,
        mutationKey: SecretApi.setSecret.key(),
        successKeys: [ManagementApi.info.key()],
    })
}
