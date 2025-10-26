import {useQueryClient} from "@tanstack/react-query"
import {AxiosError, HttpStatusCode} from "axios"
import {ReactNode, useEffect} from "react"

import {api} from "../api/api"
import {ManagementApi} from "../api/management/router"

export function AuthProvider(props: {children: ReactNode}) {
    const queryClient = useQueryClient()

    useEffect(handleEffectAxiosInterceptor, [queryClient])

    return <>{props.children}</>

    // NOTE: with React.StrictMode we will have two of them because it makes rerender twice
    // only for development env
    function handleEffectAxiosInterceptor() {
        const id = api.interceptors.response.use(
            (response) => response,
            (error: AxiosError) => {
                if (error.response?.status === HttpStatusCode.Unauthorized || error.response?.status === HttpStatusCode.Forbidden) {
                    queryClient.refetchQueries({queryKey: ManagementApi.info.key()}).then()
                }
                return Promise.reject(error)
            },
        )
        return () => api.interceptors.response.eject(id)
    }
}
