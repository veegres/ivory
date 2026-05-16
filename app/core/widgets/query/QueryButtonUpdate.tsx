import {useRouterQueryUpdate} from "../../../features/query/hook"
import {Request} from "../../../features/query/type"
import {SaveIconButton} from "../../../shared/component/button/IconButtons"

type Props = {
    id: string
    query: Request,
    onSuccess?: () => void,
}

export function QueryButtonUpdate(props: Props) {
    const {id, onSuccess} = props
    const {query, name, type} = props.query

    const update = useRouterQueryUpdate(type!, onSuccess)

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
