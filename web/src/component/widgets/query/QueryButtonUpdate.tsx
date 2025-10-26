import {useRouterQueryUpdate} from "../../../api/query/hook"
import {QueryRequest} from "../../../api/query/type"
import {SaveIconButton} from "../../view/button/IconButtons"

type Props = {
    id: string
    query: QueryRequest,
    onSuccess?: () => void,
}

export function QueryButtonUpdate(props: Props) {
    const {id, onSuccess} = props
    const {query, name, type} = props.query

    const update = useRouterQueryUpdate(type, onSuccess)

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
