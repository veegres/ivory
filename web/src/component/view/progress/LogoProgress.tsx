import {Box} from "@mui/material";
import {StylePropsMap, SxPropsMap} from "../../../app/type";

const SX: SxPropsMap = {
    box: {display: "flex", alignItems: "center", justifyContent: "center", position: "relative"},
    loader: {
        position: "absolute", borderRadius: "10em",
        boxShadow: "0px 20px 40px 25px #04c5fe, 0px -20px 40px 25px #8949ff, 20px 0px 40px 25px #ff07fb, -20px 0px 40px 25px #ff07fb",
        animation: "pulse 1s ease-out infinite alternate", "@keyframes pulse": {
            from: {transform: "scale(1)"},
            to: {transform: "scale(2)"},
        }
    },
}

const style: StylePropsMap = {
    img: {zIndex: 1, opacity: 0.7, filter: "grayscale(100%)"},
}

export function LogoProgress() {
    return (
        <Box sx={SX.box}>
            <img style={style.img} src={"/ivory.png"} width={200} height={200} alt={"Ivory"}/>
            <Box sx={SX.loader}/>
        </Box>
    )
}