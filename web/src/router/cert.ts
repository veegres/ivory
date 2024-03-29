import {useMutation, useQuery} from "@tanstack/react-query";
import {CertApi} from "../app/api";
import {CertType} from "../type/cert";
import {useMutationOptions} from "../hook/QueryCustom";

export function useRouterCertList(type: CertType) {
    return useQuery({
        queryKey: CertApi.list.key(type),
        queryFn: () => CertApi.list.fn(type)
    })
}

export function useRouterCertDelete(type: CertType) {
    const options = useMutationOptions([CertApi.list.key(type)])
    return useMutation({mutationFn: CertApi.delete, ...options})
}

export function useRouterCertUpload(type: CertType) {
    const options = useMutationOptions([CertApi.list.key(type)])
    return useMutation({mutationFn: CertApi.upload, ...options})
}

export function useRouterCertAdd(type: CertType, onSuccess?: () => void) {
    const options = useMutationOptions([CertApi.list.key(type)], onSuccess)
    return useMutation({mutationFn: CertApi.add, ...options})
}
