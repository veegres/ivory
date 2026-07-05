import {DeleteIconButton} from "../../../shared/component/button/IconButtons"
import {useRouterQueryDelete} from "../api/hook"
import {DbPlugin, Type} from "../api/type"

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
