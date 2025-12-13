import {useMutationAdapter} from "../../hook/QueryCustom"
import {ManagementApi} from "../management/router"
import {SecretApi} from "./router"

export function useRouterSecretSet() {
    return useMutationAdapter({
        mutationFn: SecretApi.set.fn,
        mutationKey: SecretApi.set.key(),
        successKeys: [ManagementApi.info.key()],
    })
}

export function useRouterSecretSkip() {
    return useMutationAdapter({
        mutationFn: SecretApi.skip.fn,
        mutationKey: SecretApi.skip.key(),
        successKeys: [ManagementApi.info.key()],
    })
}
