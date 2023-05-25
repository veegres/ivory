import {useMutationOptions} from "../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {initialApi} from "../../../app/api";
import {LoadingButton} from "@mui/lab";

type Props = {
    keyWord: string,
    refWord: string,
}

export function SecretButton(props: Props) {
    const {keyWord, refWord} = props
    const setReqOptions = useMutationOptions([["info"]])
    const setReq = useMutation(initialApi.setSecret, setReqOptions)

    return (
        <LoadingButton
            variant={"contained"}
            loading={setReq.isLoading}
            onClick={() => setReq.mutate({ref: refWord, key: keyWord})}
        >
            Done
        </LoadingButton>
    )
}
