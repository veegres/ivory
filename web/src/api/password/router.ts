import {Password, PasswordMap, PasswordType} from "./type";
import {api} from "../api";
import {R} from "../type";

export const PasswordApi = {
    list: {
        key: (type?: PasswordType) => ["password", type],
        fn: (type?: PasswordType) => api.get<R<PasswordMap>>(`/password`, {params: {type}})
            .then((response) => response.data.response),
    },
    create: {
        key: () => ["password", "create"],
        fn: (credential: Password) => api.post<R<{ key: string, credential: Password }>>(`/password`, credential)
            .then((response) => response.data.response),
    },
    update: {
        key: () => ["password", "update"],
        fn: ({uuid, credential}: { uuid: string, credential: Password }) => api.patch<R<Password>>(`/password/${uuid}`, credential)
            .then((response) => response.data.response),
    },
    delete: {
        key: () => ["password", "delete"],
        fn: (uuid: string) => api.delete(`/password/${uuid}`)
            .then((response) => response.data.response),
    },
}