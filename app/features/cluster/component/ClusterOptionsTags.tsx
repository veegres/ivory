import {AutocompleteTags} from "../../../shared/component/autocomplete/AutocompleteTags"
import {useRouterTagList} from "../../tag/api/TagHook"

type Props = {
    selected?: string[],
    onUpdate: (tags: string[]) => void,
}

export function ClusterOptionsTags(props: Props) {
    const query = useRouterTagList()
    const {data, isPending} = query
    return (
        <AutocompleteTags
            tags={data ?? []}
            selected={props.selected ?? []}
            loading={isPending}
            onUpdate={props.onUpdate}
        />
    )
}
