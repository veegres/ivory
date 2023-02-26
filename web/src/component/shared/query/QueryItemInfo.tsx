import {Box} from "@mui/material";
import {Query, SxPropsMap} from "../../../app/types";
import {QueryItemEditor} from "./QueryItemEditor";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    paper: {padding: "10px", background: "rgba(145,145,145,0.1)", borderRadius: "10px"}
}

type Props = {
    query: Query,
}

export function QueryItemInfo(props: Props) {
    return (
        <Box sx={SX.box}>
            <Box sx={SX.paper}>{props.query.description}</Box>
            <Box sx={SX.paper}>
                <QueryItemEditor value={props.query.default} editable={false}/>
            </Box>
        </Box>
    )
}
