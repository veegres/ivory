import {DeleteIconButton} from "../../../shared/component/button/IconButtons"
import {useRouterQueryDelete} from "../api/QueryHook"
import {DbPlugin, Type} from "../api/QueryType"

type Props = {
    id: string
    type: Type,
    plugin: DbPlugin,
    onSuccess?: () => void,
}

export function QueryButtonDelete(props: Props) {
    const {id, type, plugin, onSuccess} = props

    const remove = useRouterQueryDelete(type, plugin, onSuccess)

    return (
        <DeleteIconButton loading={remove.isPending} onClick={handleClick}/>
    )

    function handleClick() {
        remove.mutate(id)
    }
}
