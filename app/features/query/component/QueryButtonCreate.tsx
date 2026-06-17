import {SaveIconButton} from "../../../shared/component/button/IconButtons"
import {useRouterQueryCreate} from "../api/hook"
import {Request} from "../api/type"

type Props = {
    query: Request,
    onSuccess?: () => void,
}

export function QueryButtonCreate(props: Props) {
    const {query, name, type} = props.query

    const create = useRouterQueryCreate(type!, props.onSuccess)

    return (
        <SaveIconButton
            loading={create.isPending}
            disabled={!name || !query}
            color={"primary"}
            onClick={handleClick}
        />
    )

    function handleClick() {
        create.mutate(props.query)
    }
}
