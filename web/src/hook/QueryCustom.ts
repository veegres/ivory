// TODO think how we can optimise it and update react-query by update response without refetch
import {QueryKey, useQueryClient} from "@tanstack/react-query";
import {useSnackbar} from "notistack";
import {getErrorMessage} from "../app/utils";

/**
 * Simplify handling `onSuccess` and `onError` requests for react-query client
 * providing common approach with request refetch and custom toast messages for
 * mutation requests
 *
 * @param queryKey
 * @param onSuccess it will be fired after refetch
 */
export function useMutationOptions(queryKey?: QueryKey, onSuccess?: () => void) {
    const { enqueueSnackbar } = useSnackbar()
    const queryClient = useQueryClient();

    return {
        queryClient,
        onSuccess: queryKey === undefined ? undefined : async () => {
            await queryClient.refetchQueries(queryKey)
            if (onSuccess) onSuccess()
        },
        onError: (error: any) => enqueueSnackbar(getErrorMessage(error), { variant: "error" }),
    }
}
