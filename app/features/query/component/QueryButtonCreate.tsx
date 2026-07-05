import {SaveIconButton} from "../../../shared/component/button/IconButtons"
import {useRouterQueryCreate} from "../api/QueryHook"
import {Request} from "../api/QueryType"

type Props = {
    query: Request,
    onSuccess?: () => void,
}

export function QueryButtonCreate(props: Props) {
    const {query, name, type, plugin} = props.query

    const create = useRouterQueryCreate(type, plugin, props.onSuccess)

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
