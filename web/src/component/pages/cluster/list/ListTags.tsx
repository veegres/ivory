import {Box} from "@mui/material"

import {SxPropsMap} from "../../../../app/type"
import {useStore, useStoreAction} from "../../../../provider/StoreProvider"
import {ToggleButtonScrollable} from "../../../view/scrolling/ToggleButtonScrollable"

const SX: SxPropsMap = {
    tags: {position: "relative", height: 0, top: "-37px"},
}

type Props = {
    tags: string[],
}

export function ListTags(props: Props) {
    const {tags} = props
    const warnings = useStore(s => s.warnings)
    const activeTags = useStore(s => s.activeTags)
    const {setTags} = useStoreAction
    const warningsCount = Object.values(warnings).filter(it => it).length

    return (
        <Box sx={SX.tags}>
            <ToggleButtonScrollable
                tags={tags}
                selected={activeTags}
                onUpdate={setTags}
                warnings={warningsCount}
            />
        </Box>
    )
}
