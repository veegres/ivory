import {Switch} from "@mui/material";
import {useAppearance} from "../../../provider/AppearanceProvider";

export function MenuRefetchChanger() {
    const {state, toggleRefetchOnWindowsRefocus} = useAppearance()

    return (
        <Switch checked={state.refetchOnWindowsFocus} onClick={toggleRefetchOnWindowsRefocus}/>
    )
}
