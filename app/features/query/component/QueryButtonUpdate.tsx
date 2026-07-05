import {SaveIconButton} from "../../../shared/component/button/IconButtons"
import {useRouterQueryUpdate} from "../api/QueryHook"
import {Request} from "../api/QueryType"

type Props = {
    id: string
    query: Request,
    onSuccess?: () => void,
}

export function QueryButtonUpdate(props: Props) {
    const {id, onSuccess} = props
    const {query, name, type, plugin} = props.query

    const update = useRouterQueryUpdate(type, plugin, onSuccess)

    return (
        <SaveIconButton
            loading={update.isPending}
            disabled={!name || !query}
            color={"primary"}
            onClick={handleClick}
        />
    )

    function handleClick() {
        update.mutate({id, query: props.query})
    }
}
