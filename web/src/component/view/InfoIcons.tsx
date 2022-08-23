import {cloneElement, ReactElement} from "react";
import {useTheme} from "../../provider/ThemeProvider";
import {InfoBox} from "./InfoBox";
import {InfoTitle} from "./InfoTitle";

type Item = { icon: ReactElement, name: string, active: boolean, iconColor?: string }
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
                const color = item.active ? defaultColor : info?.palette.text.disabled
                return cloneElement(item.icon, {key: index, sx: {color}})
            })}
        </InfoBox>
    )
}

