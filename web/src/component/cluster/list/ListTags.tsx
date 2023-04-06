import {ToggleButtonScrollable} from "../../view/ToggleButtonScrollable";
import {useQuery} from "@tanstack/react-query";
import {tagApi} from "../../../app/api";
import {useStore} from "../../../provider/StoreProvider";
import {Box} from "@mui/material";
import {SxPropsMap} from "../../../type/common";

const SX: SxPropsMap = {
    tags: {position: "relative", height: 0, top: "-37px"},
}

export function ListTags() {
    const { setTags } = useStore()
    const query = useQuery(["tag/list"], tagApi.list)
    const { data } = query

    return (
        <Box sx={SX.tags}>
            <ToggleButtonScrollable tags={data ?? []} onUpdate={setTags}/>
        </Box>
    )
}
