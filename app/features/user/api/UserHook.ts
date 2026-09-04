import {useQuery} from "@tanstack/react-query"

import {useMutationAdapter} from "../../../shared/hook/QueryCustom"
import {UserApi} from "./UserRouter"
import {UserRegistered, UserRegistration} from "./UserType"

export function useRouterUserList() {
    return useQuery({
        queryKey: UserApi.list.key(),
        queryFn: () => UserApi.list.fn(),
    })
}

export function useRouterUserCreate(onSuccess?: (registered: UserRegistered) => void) {
    return useMutationAdapter({
        mutationFn: UserApi.create.fn,
        mutationKey: UserApi.create.key(),
        successKeys: [UserApi.list.key()],
        onSuccess: (_, registered) => onSuccess?.(registered),
    })
}

export function useRouterUserUpdate(onSuccess?: (registered: UserRegistered) => void) {
    return useMutationAdapter({
        mutationFn: UserApi.update.fn,
        mutationKey: UserApi.update.key(),
        successKeys: [UserApi.list.key()],
        onSuccess: (_, registered) => onSuccess?.(registered),
    })
}

export function useRouterUserDelete() {
    return useMutationAdapter({
        mutationFn: UserApi.delete.fn,
        mutationKey: UserApi.delete.key(),
        successKeys: [UserApi.list.key()],
    })
}

export function useRouterUserPassword(onSuccess?: () => void) {
    return useMutationAdapter({
        mutationFn: UserApi.password.fn,
        mutationKey: UserApi.password.key(),
        onSuccess: () => onSuccess?.(),
    })
}

export function useRouterUserPasswordReset(onSuccess?: (registration: UserRegistration) => void) {
    return useMutationAdapter({
        mutationFn: UserApi.passwordReset.issue.fn,
        mutationKey: UserApi.passwordReset.issue.key(),
        successKeys: [UserApi.list.key()],
        onSuccess: (_, registration) => onSuccess?.(registration),
    })
}

export function useRouterUserPasswordResetRevoke() {
    return useMutationAdapter({
        mutationFn: UserApi.passwordReset.revoke.fn,
        mutationKey: UserApi.passwordReset.revoke.key(),
        successKeys: [UserApi.list.key()],
    })
}

export function useRouterUserRegistrationVerify(token: string) {
    return useQuery({
        queryKey: UserApi.registration.verify.key(token),
        queryFn: () => UserApi.registration.verify.fn(token),
        retry: false,
    })
}

export function useRouterUserRegistrationPassword(onSuccess?: () => void) {
    return useMutationAdapter({
        mutationFn: UserApi.registration.password.fn,
        mutationKey: UserApi.registration.password.key(),
        onSuccess: () => onSuccess?.(),
    })
}
