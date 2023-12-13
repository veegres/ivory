import {useMutationOptions} from "../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {queryApi} from "../../../app/api";
import {SaveIconButton} from "../../view/button/IconButtons";
import {QueryRequest} from "../../../type/query";

type Props = {
    id: string
    query: QueryRequest,
    onSuccess?: () => void,
}

export function QueryButtonUpdate(props: Props) {
    const {id, onSuccess} = props
    const {query, name, type} = props.query

    const updateOptions = useMutationOptions([["query", "map", type]], onSuccess)
    const update = useMutation({mutationFn: queryApi.update, ...updateOptions})

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
