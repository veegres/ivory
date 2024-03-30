import {
    MutationOptions,
    QueryClient,
    QueryKey,
    useMutation,
    useQueryClient
} from "@tanstack/react-query";
import {useSnackbar} from "notistack";
import {getErrorMessage} from "../app/utils";

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
export function useMutationAdapter<TData, TError, TVariables, TContext>(
    options: MutationAdapterOptions<TData, TError, TVariables, TContext>
) {
    const {mutationFn, successKeys, onSuccess} = options
    const queryClient = useQueryClient();
    const {enqueueSnackbar} = useSnackbar()

    return useMutation({
        mutationFn: mutationFn,
        onSuccess: handleSuccess,
        onError: handleError,
    })

    async function handleSuccess(data: any) {
        if (successKeys) {
            for (const key of successKeys) {
                await queryClient.refetchQueries({queryKey: key})
            }
        }
        if (onSuccess) {
            onSuccess(queryClient, data)
        }
    }

    async function handleError(error: any) {
        enqueueSnackbar(getErrorMessage(error), {variant: "error"})
    }
}
