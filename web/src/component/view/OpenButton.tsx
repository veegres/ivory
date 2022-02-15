import {KeyboardArrowDown, KeyboardArrowUp} from "@mui/icons-material";

type Button = {
    show?: boolean,
    open: boolean,
    size?: number
}

export function OpenButton({open, show = true, size = 20}: Button) {
    if (!show) return null

    return open ? <KeyboardArrowUp sx={{fontSize: size}}/> : <KeyboardArrowDown sx={{fontSize: size}}/>
}
