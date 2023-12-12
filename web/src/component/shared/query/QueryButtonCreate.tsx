import {QueryRequest} from "../../../type/query";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {queryApi} from "../../../app/api";
import {SaveIconButton} from "../../view/button/IconButtons";

type Props = {
    query: QueryRequest,
    onSuccess?: () => void,
}

export function QueryButtonCreate(props: Props) {
    const {query, name, type} = props.query

    const createOptions = useMutationOptions([["query", "map", type]], props.onSuccess)
    const create = useMutation({mutationFn: queryApi.create, ...createOptions})

    return (
        <SaveIconButton
            loading={create.isPending}
            disabled={!name || !query}
            color={"success"}
            onClick={handleClick}
        />
    )

    function handleClick() {
        create.mutate(props.query)
    }
}
