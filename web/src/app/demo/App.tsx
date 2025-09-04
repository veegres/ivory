import {Box, Button} from "@mui/material";
import {SxPropsMap} from "../../type/general";
import {MenuThemeChanger} from "../../component/widgets/settings/menu/MenuThemeChanger";
import {useState} from "react";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", justifyContent: "space-evenly", alignItems: "center", width: "100vw", height: "100vh"},
    actions: {display: "flex", alignItems: "center", gap: 3},
    info: {color: "text.disabled", fontStyle: "italic", textAlign: "center"},
}

/**
 * This component is designed for testing purposes. It can ease/help you to write view components,
 * please, after your test clean this component back. Normally it shouldn't be updates at least of some
 * common parts that could help others to test something. So, please, before commiting any changes
 * remember to revert changes in this package.
 */
export function App() {
    const [boolean, setBoolean] = useState(false)
    const [num, setNum] = useState(0)

    return (
        <Box sx={SX.box}>
            <Box sx={SX.info} width={"700px"}>
                This component is designed for testing purposes. It can ease/help you to write view components
                please after your test clean this component back. Normally it shouldn't be updates at least of some
                common parts that could help others to test something. So, please, before commiting any changes
                remember to revert changes in this package.
            </Box>
            <Box sx={SX.actions}>
                <Button size={"small"} variant={"outlined"} onClick={() => setNum(num + 1)}>{num.toString()}</Button>
                <Button size={"small"} variant={"outlined"} onClick={() => setBoolean(!boolean)}>{boolean.toString()}</Button>
                <MenuThemeChanger/>
            </Box>
            <Box width={"800px"} height={"700px"}>
                {/* put your component here */}
            </Box>
        </Box>
    )
}
