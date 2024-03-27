import {useQuery} from "@tanstack/react-query";
import {CertApi} from "../app/api";
import {CertType} from "../type/cert";

export function useRouterCertList(type: CertType) {
    return useQuery({
        queryKey: ["certs", type],
        queryFn: () => CertApi.list(type)
    })
}
