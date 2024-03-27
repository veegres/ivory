import {AutocompleteTags} from "../../view/autocomplete/AutocompleteTags";
import {useRouterTagList} from "../../../router/tag";

type Props = {
    selected?: string[],
    onUpdate: (tags: string[]) => void,
}

export function OptionsTags(props: Props) {
    const query = useRouterTagList()
    const {data, isPending} = query
    const tags = data ?? [];
    return (
        <AutocompleteTags
            tags={tags}
            selected={props.selected ?? []}
            loading={isPending}
            onUpdate={props.onUpdate}
        />
    )
}
