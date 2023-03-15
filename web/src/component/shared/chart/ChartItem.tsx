import {Box} from "@mui/material";
import {SxPropsMap} from "../../../type/common";
import {
    blue, blueGrey, brown, cyan, deepOrange, deepPurple, green, indigo, orange, pink, purple, red
} from "@mui/material/colors";

const SX: SxPropsMap = {
    box: {
        display: "flex", flexDirection: "column", justifyContent: "center",
        alignItems: "center", color: "common.white", borderRadius: "5px",
        whiteSpace: "nowrap", padding: "20px 40px",
    },
    label: {fontSize: "13px", textTransform: "uppercase"},
    value: {fontSize: "30px", fontWeight: "bold"},
}

export enum Color {
    RED, PINK, PURPLE, DEEP_PURPLE, INDIGO, BLUE, CYAN,
    GREEN, ORANGE, DEEP_ORANGE, BROWN, BLUE_GREY
}
const colorContrast = 500
const colors: {[key in Color]: string} = {
    [Color.RED]: red[colorContrast],
    [Color.PINK]: pink[colorContrast],
    [Color.PURPLE]: purple[colorContrast],
    [Color.DEEP_PURPLE]: deepPurple[colorContrast],
    [Color.INDIGO]: indigo[colorContrast],
    [Color.BLUE]: blue[colorContrast],
    [Color.CYAN]: cyan[colorContrast],
    [Color.GREEN]: green[colorContrast],
    [Color.ORANGE]: orange[colorContrast],
    [Color.DEEP_ORANGE]: deepOrange[colorContrast],
    [Color.BROWN]: brown[colorContrast],
    [Color.BLUE_GREY]: blueGrey[colorContrast],
}

type Props = {
    label: string,
    value?: string | number,
    color?: Color,
    width?: string,
}

export function ChartItem(props: Props) {
    const {value, label, color, width} = props
    const w = width ? width : "200px", h = "100px"

    const bg = color !== undefined ? colors[color] : colors[Math.floor(Math.random() * Object.keys(colors).length) as Color]

    return (
        <Box sx={SX.box} minWidth={w} height={h} bgcolor={bg}>
            <Box sx={SX.label}>{label}</Box>
            <Box sx={SX.value}>{value ?? "unknown"}</Box>
        </Box>
    )
}
