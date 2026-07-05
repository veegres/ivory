import {Alert} from "@mui/material"
import {ReactNode} from "react"

import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    infoAlert: {justifyContent: "center", "& .MuiAlert-message": {textAlign: "center"}},
}

type Props = {
    text: ReactNode,
    severity?: "success" | "info" | "warning" | "error",
}

export function AlertCentered(props: Props) {
    const {severity = "info"} = props
    return (
        <Alert sx={SX.infoAlert} severity={severity} variant={"outlined"} icon={false}>
            {props.text}
        </Alert>
    )
}
