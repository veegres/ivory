import {api} from "../api"
import {R} from "../management/type"
import {UserPermissions} from "./type"

export const PermissionApi = {
    list: {
        key: () => ["permission", "list"],
        fn: () => api.get<R<UserPermissions[]>>("/permission/users")
            .then((response) => response.data.response),
    },
    request: {
        key: () => ["permission", "request"],
        fn: (permissions: string[]) =>
            api.post("/permission/request", {permissions})
                .then((response) => response.data.response),
    },
    approve: {
        key: () => ["permission", "approve"],
        fn: ({username, permissions}: { username: string, permissions: string[] }) =>
            api.post(`/permission/users/${username}/approve`, {permissions})
                .then((response) => response.data.response),
    },
    reject: {
        key: () => ["permission", "reject"],
        fn: ({username, permissions}: { username: string, permissions: string[] }) =>
            api.post(`/permission/users/${username}/reject`, {permissions})
                .then((response) => response.data.response),
    },
    delete: {
        key: () => ["permission", "delete"],
        fn: (username: string) => api.delete(`/permission/users/${username}`)
            .then((response) => response.data.response),
    },
}
