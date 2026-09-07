import {Alert} from "@mui/material"
import {ReactNode} from "react"

import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    infoAlert: {justifyContent: "center", padding: "0px 8px"},
    neutralAlert: {color: "text.secondary", border: "1px solid", borderColor: "divider", backgroundColor: "transparent"},
    center: {"& .MuiAlert-message": {textAlign: "center"}},
    left: {"& .MuiAlert-message": {textAlign: "left"}},
    right: {"& .MuiAlert-message": {textAlign: "right"}},
    justify: {"& .MuiAlert-message": {textAlign: "justify"}},
}

type Props = {
    text: ReactNode,
    severity?: "neutral" | "success" | "info" | "warning" | "error",
    align?: "left" | "center" | "right" | "justify",
}

export function AlertAlign(props: Props) {
    const {severity = "neutral", align = "center"} = props
    const al = SX[align]
    return (
        <Alert
            sx={[SX.infoAlert, al, severity === "neutral" && SX.neutralAlert]}
            severity={severity === "neutral" ? "info" : severity}
            variant={"outlined"}
            icon={false}
        >
            {props.text}
        </Alert>
    )
}
