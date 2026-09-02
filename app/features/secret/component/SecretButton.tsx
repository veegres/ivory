import {Button} from "@mui/material"

import {useRouterSecretSet} from "../api/SecretHook"

type Props = {
    keyWord: string,
}

export function SecretButton(props: Props) {
    const {keyWord} = props
    const secret = useRouterSecretSet()

    return (
        <Button
            variant={"contained"}
            loading={secret.isPending}
            onClick={() => secret.mutate({key: keyWord})}
        >
            Set
        </Button>
    )
}
