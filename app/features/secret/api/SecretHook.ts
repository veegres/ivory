import {useMutationAdapter} from "../../../shared/hook/QueryCustom"
import {ManagementApi} from "../../management/api/ManagementRouter"
import {SecretApi} from "./SecretRouter"

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
