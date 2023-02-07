import {Alert} from "@mui/material";
import React from "react";
import {SxPropsMap} from "../../app/types";

const SX: SxPropsMap = {
    infoAlert: {justifyContent: "center"}
}

type Props = { text: string }

export function InfoAlert(props: Props) {
    return (
        <Alert sx={SX.infoAlert} severity={"info"} variant={"outlined"} icon={false}>
            {props.text}
        </Alert>
    )
}
