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
        fn: (permission: string) =>
            api.post("/permission/request", {permission})
                .then((response) => response.data.response),
    },
    approve: {
        key: () => ["permission", "approve"],
        fn: ({username, permission}: { username: string, permission: string }) =>
            api.post(`/permission/users/${username}/approve`, {permission})
                .then((response) => response.data.response),
    },
    reject: {
        key: () => ["permission", "reject"],
        fn: ({username, permission}: { username: string, permission: string }) =>
            api.post(`/permission/users/${username}/reject`, {permission})
                .then((response) => response.data.response),
    },
    delete: {
        key: () => ["permission", "delete"],
        fn: (username: string) => api.delete(`/permission/users/${username}`)
            .then((response) => response.data.response),
    },
}
