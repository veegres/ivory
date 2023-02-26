import {useMutationOptions} from "../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {infoApi} from "../../../app/api";
import {useStore} from "../../../provider/StoreProvider";
import {LoadingButton} from "@mui/lab";

export function MenuEraseButton() {
    const { clear } = useStore()
    const cleanOptions = useMutationOptions([["info"]], clear)
    const clean = useMutation(infoApi.erase, cleanOptions)

    return (
        <LoadingButton
            size={"small"}
            color={"error"}
            variant={"outlined"}
            loading={clean.isLoading}
            onClick={() => clean.mutate()}
        >
            Erase
        </LoadingButton>
    )
}
