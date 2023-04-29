import {ToggleButton, ToggleButtonGroup} from "@mui/material";
import {DarkModeTwoTone, LightModeTwoTone} from "@mui/icons-material";
import {useAppearance} from "../../../provider/AppearanceProvider";
import {SxPropsMap} from "../../../type/common";

const SX: SxPropsMap = {
    button: {padding: "3px 8px"},
}

export function MenuThemeChanger() {
    const {state, setTheme} = useAppearance()

    return (
        <ToggleButtonGroup size={"small"} value={state.mode}>
            <ToggleButton sx={SX.button} value={"light"} onClick={() => setTheme("light")}>
                <LightModeTwoTone/>
            </ToggleButton>
            <ToggleButton sx={SX.button} value={"dark"} onClick={() => setTheme("dark")}>
                <DarkModeTwoTone/>
            </ToggleButton>
        </ToggleButtonGroup>
    )
}
