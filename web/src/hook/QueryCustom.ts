import {QueryKey, useQueryClient} from "@tanstack/react-query";
import {useSnackbar} from "notistack";
import {getErrorMessage} from "../app/utils";

// TODO think how we can optimise it and update react-query by updating state without refetch
/**
 * Simplify handling `onSuccess` and `onError` requests for react-query client
 * providing common approach with request refetch and custom toast messages for
 * mutation requests
 *
 * @param beforeKeys refetch QueryKeys before success
 * @param onSuccess it will be fired after refetch
 * @param afterKeys refetch QueryKeys after success
 */
export function useMutationOptions(beforeKeys?: QueryKey[], onSuccess?: (data: any) => void, afterKeys?: QueryKey[]) {
    const {enqueueSnackbar} = useSnackbar()
    const queryClient = useQueryClient();

    return {
        queryClient,
        onSuccess: handleSuccess,
        onError: (error: any) => enqueueSnackbar(getErrorMessage(error), {variant: "error"}),
    }

    async function handleSuccess(data: any) {
        if (beforeKeys) {
            for (const key of beforeKeys) {
                await queryClient.refetchQueries(key)
            }
        }
        if (onSuccess) onSuccess(data)
        if (afterKeys) {
            for (const key of afterKeys) {
                await queryClient.refetchQueries(key)
            }
        }
    }
}
