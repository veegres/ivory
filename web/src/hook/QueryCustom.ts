import {
    MutationOptions,
    QueryClient,
    QueryKey,
    useMutation,
    useQueryClient
} from "@tanstack/react-query";
import {getErrorMessage} from "../app/utils";
import {AxiosError} from "axios";
import {useSnackbar} from "../provider/SnackbarProvider";

interface MutationAdapterOptions<TData, TError, TVariables, TContext>
    extends Omit<MutationOptions<TData, TError, TVariables, TContext>, "onSuccess"> {
    successKeys?: QueryKey[],
    onSuccess?: (client: QueryClient, data: TData) => void,
}

/**
 * Simplify handling `onSuccess` and `onError` requests for useMutation hook,
 * providing common approach with request refetch and custom toast messages for
 * mutation error
 *
 * @param options accept mutation options
 * @param options.successKeys provide your useQuery keys to refetch info on success
 * @param options.onSuccess is a callback for success function request
 */
export function useMutationAdapter<TData = unknown, TError = AxiosError, TVariables = void, TContext = unknown>(
    options: MutationAdapterOptions<TData, TError, TVariables, TContext>
) {
    const {successKeys, mutationKey, onSuccess} = options

    const queryClient = useQueryClient();
    const snackbar = useSnackbar()

    return useMutation({
        ...options,
        onSuccess: handleSuccess,
        onError: handleError,
    })

    async function handleSuccess(data: any) {
        if (mutationKey) snackbar(`${mutationKey.join(" ").toUpperCase()} - DONE`, "success")
        if (successKeys) {
            for (const key of successKeys) {
                // NOTE: this doesn't refetch disabled queries
                await queryClient.refetchQueries({queryKey: key, exact: true})
            }
        }
        if (onSuccess) onSuccess(queryClient, data)
    }

    async function handleError(error: any) {
        if (mutationKey) snackbar(`${mutationKey.join(" ").toUpperCase()} - ERROR, ${getErrorMessage(error)}`, "success")
        else snackbar(getErrorMessage(error), "error")
    }
}
