import {useRouterQueryCreate} from "../../../api/query/hook";
import {QueryRequest} from "../../../api/query/type";
import {SaveIconButton} from "../../view/button/IconButtons";

type Props = {
    query: QueryRequest,
    onSuccess?: () => void,
}

export function QueryButtonCreate(props: Props) {
    const {query, name, type} = props.query

    const create = useRouterQueryCreate(type, props.onSuccess)

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
