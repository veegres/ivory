import {DeleteIconButton} from "../../../shared/component/button/IconButtons"
import {useRouterQueryDelete} from "../api/hook"
import {Type} from "../api/type"

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
