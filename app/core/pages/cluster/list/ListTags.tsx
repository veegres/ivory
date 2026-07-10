import {Box, InputBase, ToggleButton, Tooltip} from "@mui/material"

import {useRouterTagList} from "../../../../features/tag/api/TagHook"
import {ErrorSmart} from "../../../../shared/component/box/ErrorSmart"
import {ToggleButtonScrollable} from "../../../../shared/component/scrolling/ToggleButtonScrollable"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {useStore, useStoreAction} from "../../../../shared/provider/StoreProvider"

const SX: SxPropsMap = {
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
        <ToggleButtonScrollable
            items={tags.data ?? []}
            selected={activeTags}
            onUpdate={setTags}
            renderActions={renderActions()}
        />
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
            <Tooltip key={"warnings"} title={"Detected Problems"} placement={"top"}>
                <Box component={"span"}>
                    <ToggleButton
                        key={"warnings"}
                        sx={SX.element}
                        color={"warning"}
                        size={"small"}
                        selected={warningsCount > 0}
                        disabled={true}
                        value={warnings}
                    >
                        {warningsCount}
                    </ToggleButton>
                </Box>
            </Tooltip>
        ]
    }
}
