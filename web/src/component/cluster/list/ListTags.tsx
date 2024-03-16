import {ToggleButtonScrollable} from "../../view/scrolling/ToggleButtonScrollable";
import {useQuery} from "@tanstack/react-query";
import {TagApi} from "../../../app/api";
import {useStore, useStoreAction} from "../../../provider/StoreProvider";
import {Box} from "@mui/material";
import {SxPropsMap} from "../../../type/common";

const SX: SxPropsMap = {
    tags: {position: "relative", height: 0, top: "-37px"},
}

export function ListTags() {
    const {warnings, activeTags} = useStore()
    const {setTags} = useStoreAction()
    const query = useQuery({queryKey: ["tag/list"], queryFn: TagApi.list})
    const {data} = query
    const warningsCount = Object.values(warnings).filter(it => it).length

    return (
        <Box sx={SX.tags}>
            <ToggleButtonScrollable
                tags={data ?? []}
                selected={activeTags}
                onUpdate={setTags}
                warnings={warningsCount}
            />
        </Box>
    )
}
