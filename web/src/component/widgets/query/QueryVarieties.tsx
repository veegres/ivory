import {Box} from "@mui/material";

import {QueryVariety} from "../../../api/query/type";
import {SxPropsMap} from "../../../app/type";
import {QueryVarietyOptions} from "../../../app/utils";
import {InfoColorBox} from "../../view/box/InfoColorBox";

const SX: SxPropsMap = {
    box: {display: "flex", alignItems: "center", gap: 1, padding: "0px 5px"},
}

type Props = {
    varieties: QueryVariety[],
}

export function QueryVarieties(props: Props) {
    return (
        <Box sx={SX.box}>
            {props.varieties.map(v => {
                const {badge, label, color} = QueryVarietyOptions[v]
                return (
                    <InfoColorBox key={v} label={badge ?? "UNKNOWN"} title={label} bgColor={color} opacity={0.8}/>
                )
            })}
        </Box>
    )
}
