import {Switch} from "@mui/material"

import {useSettings} from "../../../../provider/AppProvider"

export function MenuRefetchChanger() {
    const {state, toggleRefetchOnWindowsRefocus} = useSettings()

    return (
        <Switch checked={state.refetchOnWindowsFocus} onClick={toggleRefetchOnWindowsRefocus}/>
    )
}
