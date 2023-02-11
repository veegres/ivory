import {ToggleButtonScrollable} from "../../view/ToggleButtonScrollable";
import {useQuery} from "@tanstack/react-query";
import {tagApi} from "../../../app/api";


export function ListTags() {
    const query = useQuery(["tag/list"], tagApi.list)
    const { data } = query

    return (
        <ToggleButtonScrollable tags={data ?? []} onUpdate={handleUpdate}/>
    )

    function handleUpdate() {

    }
}
