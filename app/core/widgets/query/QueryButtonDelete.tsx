import {useRouterQueryDelete} from "../../../features/query/hook"
import {Type} from "../../../features/query/type"
import {DeleteIconButton} from "../../../shared/component/button/IconButtons"

type Props = {
    id: string
    type: Type,
    onSuccess?: () => void,
}

export function QueryButtonDelete(props: Props) {
    const {id, type, onSuccess} = props

    const remove = useRouterQueryDelete(type, onSuccess)

    return (
        <DeleteIconButton loading={remove.isPending} onClick={handleClick}/>
    )

    function handleClick() {
        remove.mutate(id)
    }
}
