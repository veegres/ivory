import {useMutation, useQuery} from "@tanstack/react-query";
import {CertApi} from "../app/api";
import {CertType} from "../type/cert";
import {useMutationOptions} from "../hook/QueryCustom";

export function useRouterCertList(type: CertType) {
    return useQuery({
        queryKey: ["certs", type],
        queryFn: () => CertApi.list(type)
    })
}

export function useRouterCertDelete() {
    const options = useMutationOptions([["certs"]])
    return useMutation({mutationFn: CertApi.delete, ...options})
}

export function useRouterCertUpload() {
    const options = useMutationOptions([["certs"]])
    return useMutation({mutationFn: CertApi.upload, ...options})
}

export function useRouterCertAdd(onSuccess?: () => void) {
    const options = useMutationOptions([["certs"]], onSuccess)
    return useMutation({mutationFn: CertApi.add, ...options})
}
