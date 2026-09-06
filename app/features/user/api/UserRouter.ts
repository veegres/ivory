import {api} from "../../Api"
import {R} from "../../management/api/ManagementType"
import {
    User,
    UserCreateRequest,
    UserPasswordUpdateRequest,
    UserRegistered,
    UserRegistration,
    UserRegistrationPasswordRequest,
    UserRegistrationPayload,
    UserUpdateRequest
} from "./UserType"

export const UserApi = {
    list: {
        key: () => ["user", "list"],
        fn: () => api.get<R<User[]>>("/user")
            .then((response) => response.data.response),
    },
    create: {
        key: () => ["user", "create"],
        fn: (request: UserCreateRequest) => api.post<R<UserRegistered>>("/user", request)
            .then((response) => response.data.response),
    },
    update: {
        key: () => ["user", "update"],
        fn: (request: {username: string, body: UserUpdateRequest}) =>
            api.put<R<UserRegistered>>(`/user/${encodeURIComponent(request.username)}`, request.body)
                .then((response) => response.data.response),
    },
    delete: {
        key: () => ["user", "delete"],
        fn: (username: string) => api.delete<R<string>>(`/user/${encodeURIComponent(username)}`)
            .then((response) => response.data.response),
    },
    password: {
        key: () => ["user", "password"],
        fn: (request: UserPasswordUpdateRequest) => api.post<R<string>>("/user/password", request)
            .then((response) => response.data.response),
    },
    // NOTE: resetting somebody else's password takes their account over, which is
    // why it is its own route behind its own permission - what it hands out is a
    // registration link like the one a new user gets
    passwordReset: {
        issue: {
            key: () => ["user", "password", "reset", "issue"],
            fn: (username: string) =>
                api.post<R<UserRegistration>>(`/user/${encodeURIComponent(username)}/password/reset`)
                    .then((response) => response.data.response),
        },
        revoke: {
            key: () => ["user", "password", "reset", "revoke"],
            fn: (username: string) =>
                api.delete<R<string>>(`/user/${encodeURIComponent(username)}/password/reset`)
                    .then((response) => response.data.response),
        },
    },
    registration: {
        verify: {
            key: (token: string) => ["user", "registration", "verify", token],
            keyCommon: () => ["user", "registration", "verify"],
            fn: (token: string) => api.post<R<UserRegistrationPayload>>("/user/registration/verify", {token})
                .then((response) => response.data.response),
        },
        password: {
            key: () => ["user", "registration", "password"],
            fn: (request: UserRegistrationPasswordRequest) =>
                api.post<R<User>>("/user/registration/password", request)
                    .then((response) => response.data.response),
        },
    },
}
