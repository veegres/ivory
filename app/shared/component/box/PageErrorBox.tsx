import {Box, Button} from "@mui/material"

import {StylePropsMap, SxPropsMap} from "../../helper/HelperType"
import select from "../../style/select.module.css"
import {ErrorSmart} from "./ErrorSmart"

const SX: SxPropsMap = {
    page: {display: "flex", justifyContent: "center", width: "100%", padding: "50px 16px"},
    box: {display: "flex", flexDirection: "column", alignItems: "center", gap: 3, width: "100%", maxWidth: "720px"},
    header: {fontSize: "25px", fontWeight: 900, fontFamily: "monospace", textAlign: "center"},
    description: {textAlign: "center", color: "text.secondary"},
    error: {width: "100%"},
}

const style: StylePropsMap = {
    img: {opacity: 0.7, filter: "grayscale(80%)"},
}

type Props = {
    error: any,
    description?: string,
}

export function PageErrorBox(props: Props) {
    const {error, description} = props

    return (
        <Box sx={SX.page}>
            <Box sx={SX.box}>
                <img style={style.img} src={"/ivory.png"} width={120} height={120} alt={"Ivory"}/>
                <Box sx={SX.header} className={select.none}>Something went wrong</Box>
                {description && <Box sx={SX.description}>{description}</Box>}
                <Box sx={SX.error}><ErrorSmart error={error}/></Box>
                <Button variant={"contained"} onClick={handleReload}>Reload</Button>
            </Box>
        </Box>
    )

    function handleReload() {
        window.location.reload()
    }
}
