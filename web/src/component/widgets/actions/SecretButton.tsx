import {useRouterSecretSet} from "../../../api/secret/hook";
import {Button} from "@mui/material";

type Props = {
    keyWord: string,
    refWord: string,
}

export function SecretButton(props: Props) {
    const {keyWord, refWord} = props
    const secret = useRouterSecretSet()

    return (
        <Button
            size={"small"}
            variant={"contained"}
            loading={secret.isPending}
            onClick={() => secret.mutate({ref: refWord, key: keyWord})}
        >
            Set
        </Button>
    )
}
