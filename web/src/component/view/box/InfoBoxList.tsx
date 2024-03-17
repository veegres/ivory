import {cloneElement, ReactElement} from "react";
import {Box} from "@mui/material";
import {SxPropsMap} from "../../../type/common";
import {InfoBox} from "./InfoBox";
import {InfoColorBoxList} from "./InfoColorBoxList";
import {green, grey} from "@mui/material/colors";

const SX: SxPropsMap = {
    box: {position: "relative", display: "flex", alignItems: "center", justifyContent: "center"},
    badge: {
        position: "absolute", color: "white", fontSize: "8px", fontWeight: "bold",
        minWidth: "12px", height: "12px", display: "flex", justifyContent: "center", alignItems: "center",
        borderRadius: "50%", padding: "2px", textTransform: "uppercase",
    }
}

type Item = { icon: ReactElement, label: string, active: boolean, iconColor?: string, badge?: string }
type Props = {
    items: Item[],
    label?: string,
}

export function InfoBoxList(props: Props) {
    const {items, label} = props
    const titleItems = items.map(item => ({
        label: item.label,
        bgColor: item.active ? green[600] : grey[800],
    }))

    return (
        <InfoBox tooltip={<InfoColorBoxList items={titleItems} label={label}/>}>
            {items.map((item, index) => {
                const defaultColor = item.iconColor ?? "default"
                const color = item.active ? defaultColor : "text.disabled"
                const background = item.active ? green[600] : grey[800]
                return (
                    <Box key={index} sx={SX.box}>
                        {item.badge && <Box sx={{...SX.badge, background}}>{item.badge}</Box>}
                        {cloneElement(item.icon, {sx: {color}})}
                    </Box>
                )
            })}
        </InfoBox>
    )
}

