import {Box, CircularProgress, Tooltip} from "@mui/material";
import {SxPropsMap} from "../../../type/general";
import {
    blue, blueGrey, brown, cyan, deepOrange, deepPurple, green, indigo, orange, pink, purple, red
} from "@mui/material/colors";
import {useMemo} from "react";
import {IconButton, PlayIconButton, RefreshIconButton} from "../../view/button/IconButtons";
import {InfoOutlined} from "@mui/icons-material";
import {AxiosError} from "axios";

const SX: SxPropsMap = {
    box: {
        display: "flex", flexDirection: "column", justifyContent: "center",
        alignItems: "center", color: "common.white", borderRadius: "5px",
        whiteSpace: "nowrap", flexGrow: 1, minWidth: "150px", height: "100px",
        cursor: "pointer", gap: 1, padding: "0px 15px"
    },
    label: {
        display: "flex", fontSize: "13px", textTransform: "uppercase", gap: 1,
        justifyContent: "space-around", width: "100%", alignItems: "center",
    },
    value: {height: "50px"},
    text: {fontSize: "30px", fontWeight: "bold"},
}

export enum Color {
    RED, PINK, PURPLE, DEEP_PURPLE, INDIGO, BLUE, CYAN,
    GREEN, ORANGE, DEEP_ORANGE, BROWN, BLUE_GREY
}
const colorContrast = 400
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
    loading: boolean,
    label?: string,
    value?: string | number,
    color?: Color,
    error?: AxiosError<any>,
    onClick: () => void,
}

export function ChartItem(props: Props) {
    const {value, label, color, loading, error, onClick} = props

    const bg = useMemo(handleMemoBackground, [color])

    return (
        <Box sx={SX.box} bgcolor={!error ? bg : colors[Color.RED]}>
            <Box sx={SX.label}>
                <Box>{label}</Box>
                <Box>{renderIcon()}</Box>
            </Box>
            <Box sx={SX.value}>
                {!loading ? renderValue() : (
                    <CircularProgress/>
                )}
            </Box>
        </Box>
    )

    function renderIcon() {
        if (value || error) return <RefreshIconButton onClick={onClick} size={23} />

        return (
            <IconButton
                size={23}
                icon={<InfoOutlined />}
                tooltip={"These charts can create some load on the database!"}
                placement={"top"}
                onClick={() => void 0}
            />
        )
    }

    function renderValue() {
        if (!value && !error) return <PlayIconButton onClick={onClick} size={50}/>

        return (
            <Tooltip title={value ?? error?.response?.data.error ?? error?.message} placement={"top"} arrow={true}>
                <Box sx={SX.text}>{value ?? "ERROR"}</Box>
            </Tooltip>
        )
    }

    function handleMemoBackground() {
        return color !== undefined ? colors[color] : colors[Math.floor(Math.random() * Object.keys(colors).length) as Color]
    }
}
