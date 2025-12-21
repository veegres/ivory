import {
    Add,
    ArrowBack,
    AutoFixHigh,
    AutoMode,
    Cached,
    Cancel,
    CheckCircle,
    ClearAll,
    Close,
    CopyAll,
    Delete,
    Edit,
    Info,
    ManageSearch, MoreVert,
    PendingActions,
    PlayArrow,
    Receipt,
    Restore
} from "@mui/icons-material"
import {Box, IconButton as MuiIconButton, Tooltip} from "@mui/material"
import {SvgIconProps} from "@mui/material/SvgIcon/SvgIcon"
import {cloneElement, ReactElement, ReactNode, SyntheticEvent} from "react"

type Color = "inherit" | "default" | "primary" | "secondary" | "error" | "info" | "success" | "warning"
type Placement = "bottom-end" | "bottom-start" | "bottom" | "left-end" | "left-start" | "left" | "right-end" | "right-start" | "right" | "top-end" | "top-start" | "top"

type ButtonProps = Props & {
    icon: ReactElement<SvgIconProps>,
    tooltip: ReactNode,
    color?: Color,
    placement?: Placement,
    arrow?: boolean,
}

export function IconButton(props: ButtonProps) {
    const {icon, color, placement, arrow, onClick, tooltip, disabled = false, size = 32, loading} = props
    // 0.56 is the ratio for size = 32 and fontSize = 18
    const fontSize = Math.floor(size * 0.56)

    return (
        <Tooltip title={tooltip} placement={placement ?? "top"} arrow={arrow} disableInteractive>
            <Box component={"span"}>
                <MuiIconButton
                    sx={{height: `${size}px`, width: `${size}px`}}
                    color={color}
                    loading={loading}
                    disabled={loading || disabled}
                    onClick={onClick}
                >
                    {cloneElement(icon, {sx: {fontSize}})}
                </MuiIconButton>
            </Box>
        </Tooltip>
    )
}

type Props = {
    loading?: boolean,
    onClick: (event: Event | SyntheticEvent, reason?: string) => void,
    size?: number,
    disabled?: boolean,
    color?: Color,
    placement?: Placement,
    tooltip?: ReactNode,
    arrow?: boolean,
}

export function EditIconButton(props: Props) {
    return <IconButton icon={<Edit/>} tooltip={"Edit"} {...props}/>
}

export function DeleteIconButton(props: Props) {
    return <IconButton icon={<Delete/>} tooltip={"Delete"} {...props}/>
}

export function CancelIconButton(props: Props) {
    return <IconButton icon={<Cancel/>} tooltip={"Cancel"} {...props}/>
}

export function SaveIconButton(props: Props) {
    return <IconButton icon={<CheckCircle/>} tooltip={"Save"} {...props}/>
}

export function RefreshIconButton(props: Props) {
    return <IconButton icon={<Cached/>} tooltip={"Refresh"} {...props}/>
}

export function AutoRefreshIconButton(props: Props) {
    return <IconButton icon={<AutoMode/>} tooltip={"Auto Detection"} {...props}/>
}

export function BackIconButton(props: Props) {
    return <IconButton icon={<ArrowBack/>} tooltip={"Back"} {...props}/>
}

export function CloseIconButton(props: Props) {
    return <IconButton icon={<Close/>} tooltip={"Close"} {...props}/>
}

export function AddIconButton(props: Props) {
    return <IconButton icon={<Add/>} tooltip={"Add"} {...props}/>
}

export function AutoIconButton(props: Props) {
    return <IconButton icon={<AutoFixHigh/>} tooltip={"Auto"} {...props}/>
}

export function PlayIconButton(props: Props) {
    return <IconButton icon={<PlayArrow/>} tooltip={"RUN"} {...props}/>
}

export function QueryParamsIconButton(props: Props) {
    return <IconButton icon={<ManageSearch/>} tooltip={"Query Params"} {...props}/>
}

export function RestoreIconButton(props: Props) {
    return <IconButton icon={<Restore/>} tooltip={"Restore"} {...props}/>
}

export function CopyIconButton(props: Props) {
    return <IconButton icon={<CopyAll/>} tooltip={"Copy"} {...props}/>
}

export function QueryViewIconButton(props: Props) {
    return <IconButton icon={<Receipt/>} tooltip={"Query View"} {...props}/>
}

export function LogIconButton(props: Props) {
    return <IconButton icon={<PendingActions/>} tooltip={"Log"} {...props}/>
}

export function InfoIconButton(props: Props) {
    return <IconButton icon={<Info/>} tooltip={"Info"} {...props}/>
}

export function ClearAllIconButton(props: Props) {
    return <IconButton icon={<ClearAll/>} tooltip={"Clear All"} {...props}/>
}

export function MoreIconButton(props: Props) {
    return <IconButton icon={<MoreVert/>} tooltip={""} {...props}/>
}
