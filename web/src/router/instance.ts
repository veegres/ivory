import {useQuery} from "@tanstack/react-query";
import {InstanceApi} from "../app/api";
import {InstanceRequest} from "../type/instance";

export function useRouterInstanceOverview(cluster: string, request: InstanceRequest) {
    return useQuery({
        queryKey: ["instance", "overview", cluster],
        queryFn: () => InstanceApi.overview(request),
        retry: false,
    })
}

export function useRouterInstanceConfig(request: InstanceRequest, enabled: boolean) {
    return useQuery({
        queryKey: ["instance", "config", request.sidecar.host, request.sidecar.port],
        queryFn: () => InstanceApi.config(request),
        enabled,
    })
}
