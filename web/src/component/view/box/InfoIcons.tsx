import {cloneElement, ReactElement} from "react";
import {useAppearance} from "../../../provider/AppearanceProvider";
import {Box} from "@mui/material";
import {SxPropsMap} from "../../../type/common";
import {InfoBox} from "./InfoBox";
import {InfoTitle} from "./InfoTitle";

const SX: SxPropsMap = {
    box: {position: "relative", display: "flex", alignItems: "center", justifyContent: "center"},
    badge: {
        position: "absolute", color: "white", fontSize: "8px", fontWeight: "bold",
        minWidth: "12px", height: "12px", display: "flex", justifyContent: "center", alignItems: "center",
        borderRadius: "50%", padding: "2px", textTransform: "uppercase",
    }
}

const colors = {
    active: {dark: "rgba(14,80,11,0.95)", light: "rgba(16,117,19,0.95)"},
    disabled: {dark: "rgba(66,38,38,0.95)", light: "rgba(110,61,61,0.95)"}
}

type Item = { icon: ReactElement, label: string, active: boolean, iconColor?: string, badge?: string }
type Props = {
    items: Item[]
}

export function InfoIcons(props: Props) {
    const {theme} = useAppearance()
    const {items} = props
    const titleItems = items.map(item => ({
        label: item.label,
        bgColor: item.active ? colors.active[theme] : colors.disabled[theme]
    }))

    return (
        <InfoBox tooltip={<InfoTitle items={titleItems}/>}>
            {items.map((item, index) => {
                const defaultColor = item.iconColor ?? "default"
                const color = item.active ? defaultColor : "text.disabled"
                const background = item.active ? colors.active[theme] : colors.disabled[theme]
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

