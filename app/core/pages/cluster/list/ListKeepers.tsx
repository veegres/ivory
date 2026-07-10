import {ToggleButton} from "@mui/material"

import {HiddenScrolling} from "../../../../shared/component/scrolling/HiddenScrolling"
import {EnumOptions, SxPropsMap} from "../../../../shared/helper/HelperType"
import {KeeperPluginOptions} from "../../../../shared/helper/HelperUtils"
import {useStore, useStoreAction} from "../../../../shared/provider/StoreProvider"

const SX: SxPropsMap = {
    element: {padding: "3px 7px", borderRadius: "3px", lineHeight: "1"},
}

export function ListKeepers() {
    const keeper = useStore(s => s.activeClusterKeeperPlugin)
    const {setClusterKeeperPlugin} = useStoreAction

    return (
        <HiddenScrolling arrowHeight={"23px"} position={"center"}>
            {Object.entries(KeeperPluginOptions).map(([k, v]) => renderButtons(k, v))}
        </HiddenScrolling>
    )

    function renderButtons(key: string, value: EnumOptions) {
        return (
            <ToggleButton
                sx={SX.element}
                color={"info"}
                size={"small"}
                key={value.key}
                selected={keeper === key}
                value={key}
                onClick={(_, v) => setClusterKeeperPlugin(v)}
            >
                {value.label}
            </ToggleButton>
        )
    }
}