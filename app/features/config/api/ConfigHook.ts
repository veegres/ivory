import {useMutationAdapter} from "../../../shared/hook/QueryCustom"
import {ManagementApi} from "../../management/api/ManagementRouter"
import {ConfigApi} from "./ConfigRouter"

export function useRouterConfigSet() {
    return useMutationAdapter({
        mutationFn: ConfigApi.setAppConfig.fn,
        mutationKey: ConfigApi.setAppConfig.key(),
        successKeys: [ManagementApi.info.key()],
    })
}