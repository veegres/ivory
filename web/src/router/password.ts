import {useQuery} from "@tanstack/react-query";
import {PasswordApi} from "../app/api";
import {PasswordType} from "../type/password";

export function useRouterPassword(type?: PasswordType) {
    return useQuery({
        queryKey: ["password", type],
        queryFn: () => PasswordApi.list(type)
    })
}
