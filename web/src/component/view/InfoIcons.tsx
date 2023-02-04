import React, {cloneElement, ReactElement} from "react";
import {useTheme} from "../../provider/ThemeProvider";
import {InfoBox} from "./InfoBox";
import {InfoTitle} from "./InfoTitle";
import {Box} from "@mui/material";

const SX = {
    box: { display: "flex", alignItems: "center", justifyContent: "center" },
    badge: {
        position: "absolute", color: "white", fontSize: "8px", fontWeight: "bold",
        minWidth: "12px", height: "12px", display: "flex", justifyContent: "center", alignItems: "center",
        borderRadius: "50%", padding: "2px", textTransform: "uppercase",
    }
}

type Item = { icon: ReactElement, label: string, active: boolean, iconColor?: string, badge?: string }
type Props = {
    items: Item[]
}

export function InfoIcons(props: Props) {
    const { info, mode } = useTheme()
    const { items } = props
    const titleItems = items.map(item => ({
        ...item,
        value: item.active ? "Yes" : "No",
        bgColor: item.active ? info?.palette.success[mode] : info?.palette.error[mode]
    }))

    return (
        <InfoBox tooltip={<InfoTitle items={titleItems} />}>
            {items.map((item, index) => {
                const defaultColor = item.iconColor ?? "default"
                const color = item.active ? defaultColor : info?.palette.text.disabled!!
                const background = item.active ? "rgba(19,88,131,0.95)" : "rgba(66,38,38,0.95)"
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

