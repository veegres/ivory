import {QueryType} from "../../../type/query";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {queryApi} from "../../../app/api";
import {DeleteIconButton} from "../../view/button/IconButtons";

type Props = {
    id: string
    type: QueryType,
    onSuccess?: () => void,
}

export function QueryButtonDelete(props: Props) {
    const {id, type, onSuccess} = props

    const removeOptions = useMutationOptions([["query", "map", type]], onSuccess)
    const remove = useMutation({mutationFn: queryApi.delete, ...removeOptions})

    return (
        <DeleteIconButton loading={remove.isPending} onClick={handleClick}/>
    )

    function handleClick() {
        remove.mutate(id)
    }
}
