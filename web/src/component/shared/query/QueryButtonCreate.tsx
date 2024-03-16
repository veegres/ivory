import {QueryRequest} from "../../../type/query";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {QueryApi} from "../../../app/api";
import {SaveIconButton} from "../../view/button/IconButtons";

type Props = {
    query: QueryRequest,
    onSuccess?: () => void,
}

export function QueryButtonCreate(props: Props) {
    const {query, name, type} = props.query

    const createOptions = useMutationOptions([["query", "map", type]], props.onSuccess)
    const create = useMutation({mutationFn: QueryApi.create, ...createOptions})

    return (
        <SaveIconButton
            loading={create.isPending}
            disabled={!name || !query}
            color={"primary"}
            onClick={handleClick}
        />
    )

    function handleClick() {
        create.mutate(props.query)
    }
}
