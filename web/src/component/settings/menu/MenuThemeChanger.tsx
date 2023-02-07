import {ToggleButtonGroup} from "@mui/lab";
import {ToggleButton} from "@mui/material";
import {DarkModeTwoTone, LightModeTwoTone} from "@mui/icons-material";
import {useTheme} from "../../../provider/ThemeProvider";

const SX = {
    button: {padding: "3px 8px"},
}
export function MenuThemeChanger() {
    const {mode, set} = useTheme()

    return (
        <ToggleButtonGroup size={"small"} value={mode}>
            <ToggleButton sx={SX.button} value={"light"} onClick={() => set("light")}>
                <LightModeTwoTone/>
            </ToggleButton>
            <ToggleButton sx={SX.button} value={"dark"} onClick={() => set("dark")}>
                <DarkModeTwoTone/>
            </ToggleButton>
        </ToggleButtonGroup>
    )
}
