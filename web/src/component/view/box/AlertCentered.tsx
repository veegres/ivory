import {Alert} from "@mui/material"

import {SxPropsMap} from "../../../app/type"

const SX: SxPropsMap = {
    infoAlert: {justifyContent: "center", "& .MuiAlert-message": {textAlign: "center"}},
}

type Props = {
    text: string,
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
