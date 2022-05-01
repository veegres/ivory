import {Alert} from "@mui/material";
import React from "react";

const SX = {
    infoAlert: {justifyContent: 'center'}
}

type Props = { text: string }

export function Info(props: Props) {
    return (
        <Alert sx={SX.infoAlert} severity={"info"} variant={"outlined"} icon={false}>
            {props.text}
        </Alert>
    )
}
