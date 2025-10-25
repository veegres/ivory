import {useMutationAdapter} from "../../hook/QueryCustom";
import {ManagementApi} from "../management/router";
import {ConfigApi} from "./router";

export function useRouterConfigSet() {
    return useMutationAdapter({
        mutationFn: ConfigApi.setAppConfig.fn,
        mutationKey: ConfigApi.setAppConfig.key(),
        successKeys: [ManagementApi.info.key()],
    })
}