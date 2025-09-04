import {QueryType} from "../../../type/query";
import {DeleteIconButton} from "../../view/button/IconButtons";
import {useRouterQueryDelete} from "../../../router/query";

type Props = {
    id: string
    type: QueryType,
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
