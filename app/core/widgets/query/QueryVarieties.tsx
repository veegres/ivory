import {Box} from "@mui/material"

import {VarietyType} from "../../../features/query/type"
import {InfoColorBox} from "../../../shared/component/box/InfoColorBox"
import {SxPropsMap} from "../../../shared/helper/type"
import {QueryVarietyOptions} from "../../../shared/helper/utils"

const SX: SxPropsMap = {
    box: {display: "flex", alignItems: "center", gap: "4px", padding: "0px 5px", height: "100%"},
}

type Props = {
    varieties: VarietyType[],
}

export function QueryVarieties(props: Props) {
    return (
        <Box sx={SX.box}>
            {props.varieties.map(v => {
                const {badge, label, color} = QueryVarietyOptions[v]
                return (
                    <InfoColorBox
                        key={v}
                        label={badge ?? "UNKNOWN"}
                        title={label}
                        bgColor={color}
                        opacity={0.8}
                    />
                )
            })}
        </Box>
    )
}
