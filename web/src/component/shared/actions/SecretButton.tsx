import {useMutationOptions} from "../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {InitialApi} from "../../../app/api";
import {LoadingButton} from "@mui/lab";

type Props = {
    keyWord: string,
    refWord: string,
}

export function SecretButton(props: Props) {
    const {keyWord, refWord} = props
    const setReqOptions = useMutationOptions([["info"]])
    const setReq = useMutation({mutationFn: InitialApi.setSecret, ...setReqOptions})

    return (
        <LoadingButton
            variant={"contained"}
            loading={setReq.isPending}
            onClick={() => setReq.mutate({ref: refWord, key: keyWord})}
        >
            Set
        </LoadingButton>
    )
}
