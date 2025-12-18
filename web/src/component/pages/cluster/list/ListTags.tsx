import {Box, InputBase, ToggleButton} from "@mui/material"

import {useRouterTagList} from "../../../../api/tag/hook"
import {SxPropsMap} from "../../../../app/type"
import {useStore, useStoreAction} from "../../../../provider/StoreProvider"
import {ErrorSmart} from "../../../view/box/ErrorSmart"
import {ToggleButtonScrollable} from "../../../view/scrolling/ToggleButtonScrollable"

const SX: SxPropsMap = {
    tags: {position: "relative", height: 0, top: "-37px"},
    input: {padding: "0px", width: "100px", height: "14px", fontSize: "14px"},
}

export function ListTags() {
    const tags = useRouterTagList()
    const warnings = useStore(s => s.warnings)
    const search = useStore(s => s.searchCluster)
    const activeTags = useStore(s => s.activeTags)
    const {setTags, setSearchCluster} = useStoreAction

    const warningsCount = Object.values(warnings).filter(it => it).length
    if (tags.error) return <ErrorSmart error={tags.error}/>

    return (
        <Box sx={SX.tags}>
            <ToggleButtonScrollable
                tags={tags.data ?? []}
                selected={activeTags}
                onUpdate={setTags}
                renderActions={renderActions()}
            />
        </Box>
    )

    function renderActions() {
        return [
            <InputBase
                key={"search"}
                type={"text"}
                size={"small"}
                slotProps={{input: {sx: SX.input}}}
                placeholder={"Filter by name"}
                value={search}
                onChange={e => setSearchCluster(e.target.value)}
            />,
            <ToggleButton
                key={"warnings"}
                sx={SX.element}
                color={"warning"}
                size={"small"}
                selected={warningsCount > 0}
                disabled
                value={warnings}
            >
                {warningsCount}
            </ToggleButton>
        ]
    }
}
