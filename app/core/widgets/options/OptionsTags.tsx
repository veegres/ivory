import {useRouterTagList} from "../../../features/tag/hook"
import {AutocompleteTags} from "../../../shared/component/autocomplete/AutocompleteTags"

type Props = {
    selected?: string[],
    onUpdate: (tags: string[]) => void,
}

export function OptionsTags(props: Props) {
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
