import {AutocompleteTags} from "../../view/AutocompleteTags";



export function OverviewSettingsTags() {
    const onUpdate = (tags: string[]) => console.log(tags)
    const tags = ["test"];

    return (
        <AutocompleteTags tags={tags} loading={false} onUpdate={onUpdate}/>
    )
}
