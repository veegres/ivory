import {SvgIconProps, Tooltip} from "@mui/material"
import {cloneElement, ReactElement} from "react"

import {IconButton} from "./IconButtons"
import {SimpleButton} from "./SimpleButton"

type Props = {
    title: string,
    label?: string,
    icon: ReactElement<SvgIconProps>,
    size?: number,
    disabled?: boolean,
    variant?: "button" | "icon" | "button_label",
    onClick: () => void,
}

export function TriggerButton(props: Props) {
    const {title, label, icon, size, disabled = false, variant = "icon", onClick} = props

    if (variant === "button") {
        // NOTE: mirrors the Refresher button dimensions (28px high, +4px wide,
        // 16px glyph) so the outlined buttons in a toolbar all match
        const side = size ?? 28
        return (
            <Tooltip title={title} arrow={true} placement={"top"}>
                <SimpleButton
                    sx={{height: `${side}px`, width: `${side + 4}px`}}
                    disabled={disabled}
                    onClick={onClick}
                >
                    {cloneElement(icon, {sx: {fontSize: "16px"}})}
                </SimpleButton>
            </Tooltip>
        )
    }
    if (variant === "button_label") {
        return (
            <SimpleButton
                variant={"text"}
                startIcon={icon}
                disabled={disabled}
                sx={{height: `${size}px`, padding: "5px 8px 5px 13px", lineHeight: 1}}
                onClick={onClick}
            >
                {label}
            </SimpleButton>
        )
    }
    return <IconButton tooltip={title} icon={icon} size={size} disabled={disabled} onClick={onClick}/>
}
