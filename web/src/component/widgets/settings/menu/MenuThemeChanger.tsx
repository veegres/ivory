import {ToggleButton, ToggleButtonGroup, Tooltip} from "@mui/material";
import {BrightnessMediumTwoTone, DarkModeTwoTone, LightModeTwoTone} from "@mui/icons-material";
import {Mode, useSettings} from "../../../../provider/AppProvider";

import {SxPropsMap} from "../../../../app/type";

const SX: SxPropsMap = {
    button: {padding: "3px 8px"},
}

export function MenuThemeChanger() {
    const {state, setTheme} = useSettings()

    return (
        <ToggleButtonGroup size={"small"} value={state.mode}>
            <ToggleButton sx={SX.button} value={Mode.LIGHT} onClick={() => setTheme(Mode.LIGHT)}>
                <Tooltip title={Mode.LIGHT.toUpperCase()} placement={"top"}>
                    <LightModeTwoTone/>
                </Tooltip>
            </ToggleButton>
            <ToggleButton sx={SX.button} value={Mode.SYSTEM} onClick={() => setTheme(Mode.SYSTEM)}>
                <Tooltip title={Mode.SYSTEM.toUpperCase()} placement={"top"}>
                    <BrightnessMediumTwoTone/>
                </Tooltip>
            </ToggleButton>
            <ToggleButton sx={SX.button} value={Mode.DARK} onClick={() => setTheme(Mode.DARK)}>
                <Tooltip title={Mode.DARK.toUpperCase()} placement={"top"}>
                    <DarkModeTwoTone/>
                </Tooltip>
            </ToggleButton>
        </ToggleButtonGroup>
    )
}
