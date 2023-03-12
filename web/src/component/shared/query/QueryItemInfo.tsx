import {Box} from "@mui/material";
import {SxPropsMap} from "../../../type/common";
import {QueryEditor} from "./QueryEditor";
import {QueryInfo} from "./QueryInfo";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    paper: {padding: "10px", background: "rgba(145,145,145,0.1)", borderRadius: "10px"}
}

type Props = {
    description: string,
    query: string,
}

export function QueryItemInfo(props: Props) {
    const {query, description} = props
    return (
        <Box sx={SX.box}>
            <QueryInfo>{description}</QueryInfo>
            <QueryInfo>
                <QueryEditor value={query} editable={false}/>
            </QueryInfo>
        </Box>
    )
}
