import {useQuery} from "@tanstack/react-query";
import {CertApi} from "../app/api";
import {CertType} from "../type/cert";
import {useMutationAdapter} from "../hook/QueryCustom";

export function useRouterCertList(type: CertType) {
    return useQuery({
        queryKey: CertApi.list.key(type),
        queryFn: () => CertApi.list.fn(type)
    })
}

export function useRouterCertDelete(type: CertType) {
    return useMutationAdapter({
        mutationFn: CertApi.delete.fn,
        successKeys: [CertApi.list.key(type)]
    })
}

export function useRouterCertUpload(type: CertType) {
    return useMutationAdapter({
        mutationFn: CertApi.upload.fn,
        successKeys: [CertApi.list.key(type)]
    })
}

export function useRouterCertAdd(type: CertType, onSuccess?: () => void) {
    return useMutationAdapter({
        mutationFn: CertApi.add.fn,
        successKeys: [CertApi.list.key(type)],
        onSuccess: onSuccess,
    })
}
