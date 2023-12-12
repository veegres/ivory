import {QueryVarietyOptions} from "../../../app/utils";
import {InfoColorBox} from "../../view/box/InfoColorBox";
import {QueryVariety} from "../../../type/query";
import {SxPropsMap} from "../../../type/common";
import {Box} from "@mui/material";

const SX: SxPropsMap = {
    box: {display: "flex", alignItems: "center", gap: 1, height: "100%"},
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
