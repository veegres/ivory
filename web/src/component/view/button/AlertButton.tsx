import {Box} from "@mui/material";
import {AlertDialog} from "../dialog/AlertDialog";
import {LoadingButton} from "@mui/lab";
import {ReactNode, useState} from "react";
import {SxPropsMap} from "../../../type/common";

const SX: SxPropsMap = {
    button: {minWidth: "10px", textWrap: "nowrap"},
}

type Props = {
    children: ReactNode,
    title: string,
    content: string,
    onClick: () => void,
    loading?: boolean,
    disabled?: boolean,
    size?: "small" | "medium" | "large",
    color?: "secondary" | "success" | "inherit" | "warning" | "error" | "primary" | "info",
}

export function AlertButton(props: Props) {
    const {children, title, content, onClick} = props
    const {loading, disabled, size, color} = props
    const [open, setOpen] = useState(false)

    return (
        <Box>
            <AlertDialog
                open={open}
                title={title}
                content={content}
                onAgree={onClick}
                onClose={() => setOpen(false)}
            />
            <LoadingButton
                sx={SX.button}
                size={size}
                color={color}
                disabled={disabled}
                loading={loading}
                onClick={() => setOpen(true)}>
                {children}
            </LoadingButton>
        </Box>
    )
}
