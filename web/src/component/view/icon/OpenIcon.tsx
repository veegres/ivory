import {KeyboardArrowDown, KeyboardArrowUp} from "@mui/icons-material"

type Icon = {
    open: boolean,
    show?: boolean,
    size?: number
}

export function OpenIcon({open, show = true, size = 20}: Icon) {
    if (!show) return null

    return open ? (
        <KeyboardArrowUp sx={{fontSize: size}}/>
    ) : (
        <KeyboardArrowDown sx={{fontSize: size}}/>
    )
}
