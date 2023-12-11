import {Box} from "@mui/material";
import {SxPropsMap} from "../../../type/common";
import {QueryEditor} from "./QueryEditor";
import {QueryDescription} from "./QueryDescription";
import {QueryVariety} from "../../../type/query";
import {QueryVarietyOptions} from "../../../app/utils";
import {InfoColorBox} from "../../view/box/InfoColorBox";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    paper: {padding: "10px", background: "rgba(145,145,145,0.1)", borderRadius: "10px"},
    varieties: {display: "grid", gridAutoColumns: "minmax(0, 1fr)", gridAutoFlow: "column", gap: 2},
}

type Props = {
    description: string,
    query: string,
    varieties?: QueryVariety[],
}

export function QueryItemInfo(props: Props) {
    const {query, description, varieties} = props
    return (
        <Box sx={SX.box}>
            {varieties && (
                <QueryDescription>
                    <Box sx={SX.varieties}>
                        {varieties.map(v => {
                            const {label, color} = QueryVarietyOptions[v]
                            return (
                                <InfoColorBox key={v} label={label} bgColor={color} opacity={0.8}/>
                            )
                        })}
                    </Box>
                </QueryDescription>
            )}
            {description && <QueryDescription>{description}</QueryDescription>}
            <QueryDescription>
                <QueryEditor value={query} editable={false}/>
            </QueryDescription>
        </Box>
    )
}
