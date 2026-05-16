import {useRouterQueryCreate} from "../../../features/query/hook"
import {Request} from "../../../features/query/type"
import {SaveIconButton} from "../../../shared/component/button/IconButtons"

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
