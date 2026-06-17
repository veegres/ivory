import {Box, Button, Tooltip} from "@mui/material"
import {ReactNode, useState} from "react"

import {SxPropsMap} from "../../helper/type"
import {AlertDialog} from "../dialog/AlertDialog"

const SX: SxPropsMap = {
    button: {padding: "3px", minWidth: 0, textWrap: "nowrap", lineHeight: 1},
    box: {display: "flex", alignItems: "center"},
}

type Props = {
    children?: ReactNode,
    label: ReactNode,
    title: string,
    description: ReactNode | string,
    onClick?: () => void,
    tooltip?: ReactNode,
    loading?: boolean,
    disabled?: boolean,
    size?: "small" | "medium" | "large",
    color?: "secondary" | "success" | "inherit" | "warning" | "error" | "primary" | "info",
    variant?: "text" | "outlined" | "contained",
}

export function AlertButton(props: Props) {
    const {children, title, description, label, variant, onClick} = props
    const {loading, disabled, size, color, tooltip} = props
    const [open, setOpen] = useState(false)
    const borderColor = color === "inherit" ? "divider" : undefined

    return (
        <>
            <AlertDialog
                open={open}
                title={title}
                description={description}
                onAgree={onClick}
                onClose={() => setOpen(false)}
            >
                {children}
            </AlertDialog>
            <Tooltip title={tooltip} placement={"top"} arrow={true}>
                <Box sx={SX.box}>
                    <Button
                        sx={[SX.button, {borderColor}]}
                        size={size}
                        color={color}
                        variant={variant}
                        disabled={disabled}
                        loading={loading}
                        onClick={() => setOpen(true)}
                    >
                        {label}
                    </Button>
                </Box>
            </Tooltip>
        </>
    )
}
