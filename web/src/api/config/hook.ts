import {useMutationAdapter} from "../../hook/QueryCustom";
import {ConfigApi} from "./router";
import {ManagementApi} from "../management/router";

export function useRouterConfigSet() {
    return useMutationAdapter({
        mutationFn: ConfigApi.setAppConfig.fn,
        mutationKey: ConfigApi.setAppConfig.key(),
        successKeys: [ManagementApi.info.key()],
    })
}