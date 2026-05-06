import {ToggleButton, ToggleButtonGroup} from "@mui/material"

import {SxPropsMap} from "../../../app/type"

const SX: SxPropsMap = {
    group: {gap: 1, ".MuiToggleButtonGroup-grouped": {border: 1, borderColor: "divider", borderRadius: 1, padding: 0.5}},
}

export interface Tabs {
    [key: number]: {label: string}
}

type Props = {
    tabs: Tabs,
    tab: number,
    setTab: (index: number) => void,
}

export function TabsButton(props: Props) {
    const {tabs, tab, setTab} = props
    return (
        <ToggleButtonGroup
            sx={SX.group}
            value={tab}
            fullWidth
            exclusive
            onChange={(_, value) => setTab(value ?? tab)}
        >
            {Object.entries(tabs).map(([key, value]) => (
                <ToggleButton key={key} value={Number(key)}>{value.label}</ToggleButton>
            ))}
        </ToggleButtonGroup>
    )
}
