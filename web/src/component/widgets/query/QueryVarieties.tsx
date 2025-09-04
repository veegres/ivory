import {QueryVarietyOptions} from "../../../app/utils";
import {InfoColorBox} from "../../view/box/InfoColorBox";
import {QueryVariety} from "../../../type/query";
import {SxPropsMap} from "../../../type/general";
import {Box} from "@mui/material";

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
