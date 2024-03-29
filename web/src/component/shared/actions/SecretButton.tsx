import {LoadingButton} from "@mui/lab";
import {useRouterSecretSet} from "../../../router/secret";

type Props = {
    keyWord: string,
    refWord: string,
}

export function SecretButton(props: Props) {
    const {keyWord, refWord} = props
    const secret = useRouterSecretSet()

    return (
        <LoadingButton
            variant={"contained"}
            loading={secret.isPending}
            onClick={() => secret.mutate({ref: refWord, key: keyWord})}
        >
            Set
        </LoadingButton>
    )
}
