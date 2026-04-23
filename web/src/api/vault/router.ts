import {api} from "../api"
import {R} from "../management/type"
import {Vault, VaultMap, VaultType} from "./type"

export const VaultApi = {
    list: {
        key: (type?: VaultType) => ["vault", type],
        fn: (type?: VaultType) => api.get<R<VaultMap>>("/vault", {params: {type}})
            .then((response) => response.data.response),
    },
    create: {
        key: () => ["vault", "create"],
        fn: (vault: Vault) => api.post<R<{ key: string, vault: Vault }>>("/vault", vault)
            .then((response) => response.data.response),
    },
    update: {
        key: () => ["vault", "update"],
        fn: ({uuid, vault}: { uuid: string, vault: Vault }) => api.patch<R<Vault>>(`/vault/${uuid}`, vault)
            .then((response) => response.data.response),
    },
    delete: {
        key: () => ["vault", "delete"],
        fn: (uuid: string) => api.delete(`/vault/${uuid}`)
            .then((response) => response.data.response),
    },
}