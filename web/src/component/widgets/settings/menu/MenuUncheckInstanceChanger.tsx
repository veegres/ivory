import {Switch} from "@mui/material";
import {useSettings} from "../../../../provider/SettingsProvider";

export function MenuUncheckInstanceChanger() {
    const {state, toggleUncheckInstanceBlockOnClusterChange} = useSettings()

    return (
        <Switch checked={state.uncheckInstanceBlockOnClusterChange} onClick={toggleUncheckInstanceBlockOnClusterChange}/>
    )
}
