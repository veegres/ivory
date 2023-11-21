import {AutocompleteTags} from "../../view/autocomplete/AutocompleteTags";
import {useQuery} from "@tanstack/react-query";
import {tagApi} from "../../../app/api";

type Props = {
    selected?: string[],
    onUpdate: (tags: string[]) => void,
}

export function OptionsTags(props: Props) {
    const query = useQuery({queryKey: ["tag/list"], queryFn: tagApi.list})
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
