import {Alert} from "@mui/material"
import {ReactNode} from "react"

import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    infoAlert: {justifyContent: "center", "& .MuiAlert-message": {textAlign: "center"}},
    neutralAlert: {color: "text.secondary", border: "1px solid", borderColor: "divider", backgroundColor: "transparent"},
}

type Props = {
    text: ReactNode,
    severity?: "neutral" | "success" | "info" | "warning" | "error",
}

export function AlertCentered(props: Props) {
    const {severity = "neutral"} = props
    return (
        <Alert
            sx={severity === "neutral" ? [SX.infoAlert, SX.neutralAlert] : SX.infoAlert}
            severity={severity === "neutral" ? "info" : severity}
            variant={"outlined"}
            icon={false}
        >
            {props.text}
        </Alert>
    )
}
