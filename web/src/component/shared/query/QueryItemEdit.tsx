import {QueryItemEditor} from "./QueryItemEditor";
import {Query, SxPropsMap} from "../../../app/types";
import {Box} from "@mui/material";
import {CopyIconButton, SaveIconButton} from "../../view/IconButtons";

const SX: SxPropsMap = {
    box: {display: "flex", gap: 1},
    editor: {flexGrow: 1},
    buttons: {display: "flex", flexDirection: "column", gap: 1}
}

type Props = {
    query: Query,
}

export function QueryItemEdit(props: Props) {
    return (
        <Box sx={SX.box}>
            <Box sx={SX.editor}>
                <QueryItemEditor value={props.query.default} editable={true} onUpdate={() => {}}/>
            </Box>
            <Box sx={SX.buttons}>
                <SaveIconButton placement={"left"} onClick={() => {}}/>
                <CopyIconButton placement={"left"} onClick={handleCopy}/>
            </Box>
        </Box>
    )

    function handleCopy() {
        return navigator.clipboard.writeText(props.query.default)
    }
}
