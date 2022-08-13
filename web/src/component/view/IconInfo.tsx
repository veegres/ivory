import {Box} from "@mui/material";
import {cloneElement, ReactElement} from "react";
import {useTheme} from "../../provider/ThemeProvider";
import {InfoBox} from "./InfoBox";

const SX = {
    row: { display: "flex", gap: 1, justifyContent: "space-between" },
}

type Item = { icon: ReactElement, name: string, active: boolean, color?: string }
type Props = {
    items: Item[]
}

export function IconInfo(props: Props) {
    const { info, mode } = useTheme()
    const { items } = props
    return (
        <InfoBox tooltip={renderTitle(items)}>
            {items.map((item, index) => {
                const defaultColor = item.color ?? "default"
                const color = item.active ? defaultColor : info?.palette.text.disabled
                return cloneElement(item.icon, {key: index, sx: {color}})
            })}
        </InfoBox>
    )

    function renderTitle(items: Item[]) {
        return (
            <Box>
                {items.map((item, index) => (
                    <Box key={index} sx={SX.row}>
                        <Box>{item.name}:</Box>
                        {renderAnswer(item.active)}
                    </Box>
                ))}
            </Box>
        )
    }

    function renderAnswer(active?: boolean) {
        const answer = active ? "YES" : "NO"
        const color = active ? info?.palette.success[mode] : info?.palette.error[mode]
        return (
            <Box sx={{ color }}>{answer}</Box>
        )
    }
}

