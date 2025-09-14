import {ToggleButtonScrollable} from "../../../view/scrolling/ToggleButtonScrollable";
import {useStore, useStoreAction} from "../../../../provider/StoreProvider";
import {Box} from "@mui/material";
import {useRouterTagList} from "../../../../api/tag/hook";
import {SxPropsMap} from "../../../../app/type";

const SX: SxPropsMap = {
    tags: {position: "relative", height: 0, top: "-37px"},
}

export function ListTags() {
    const warnings = useStore(s => s.warnings)
    const activeTags = useStore(s => s.activeTags)
    const {setTags} = useStoreAction
    const query = useRouterTagList()
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
