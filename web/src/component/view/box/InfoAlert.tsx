import {Alert} from "@mui/material";
import {SxPropsMap} from "../../../type/general";

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
