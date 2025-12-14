import {Box} from "@mui/material"

import {useRouterTagList} from "../../../../api/tag/hook"
import {SxPropsMap} from "../../../../app/type"
import {useStore, useStoreAction} from "../../../../provider/StoreProvider"
import {ErrorSmart} from "../../../view/box/ErrorSmart"
import {ToggleButtonScrollable} from "../../../view/scrolling/ToggleButtonScrollable"

const SX: SxPropsMap = {
    tags: {position: "relative", height: 0, top: "-37px"},
}

export function ListTags() {
    const tags = useRouterTagList()
    const warnings = useStore(s => s.warnings)
    const activeTags = useStore(s => s.activeTags)
    const {setTags} = useStoreAction
    const warningsCount = Object.values(warnings).filter(it => it).length
    if (tags.error) return <ErrorSmart error={tags.error}/>

    return (
        <Box sx={SX.tags}>
            <ToggleButtonScrollable
                tags={tags.data ?? []}
                selected={activeTags}
                onUpdate={setTags}
                warnings={warningsCount}
            />
        </Box>
    )
}
