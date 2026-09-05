import {ToggleButton, ToggleButtonGroup} from "@mui/material"

import {Feature} from "../../../features/Feature"
import {ManageAccess} from "../../../features/management/component/ManageAccess"
import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    group: {gap: 1, ".MuiToggleButtonGroup-grouped": {border: 1, borderColor: "divider", borderRadius: 1, padding: "0px 10px"}},
}

export interface Tabs {
    [key: number]: {label: string, feature?: Feature},
}

type Props = {
    tabs: Tabs,
    tab: number,
    setTab: (index: number) => void,
    fullWidth?: boolean,
}

export function TabsButton(props: Props) {
    const {tabs, tab, setTab, fullWidth = true} = props
    return (
        <ToggleButtonGroup
            sx={SX.group}
            value={tab}
            fullWidth={fullWidth}
            exclusive={true}
            onChange={(_, value) => setTab(value ?? tab)}
        >
            {Object.entries(tabs).map(([key, tab]) => !tab.feature ? renderButton(key, tab.label) : (
                <ManageAccess key={key} feature={tab.feature}>{renderButton(key, tab.label)}</ManageAccess>
            ))}
        </ToggleButtonGroup>
    )

    function renderButton(key: string, tab: string) {
        return (
            <ToggleButton key={key} value={Number(key)}>{tab}</ToggleButton>
        )
    }
}
