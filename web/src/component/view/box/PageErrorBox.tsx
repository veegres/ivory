import {PageMainBox} from "./PageMainBox";
import {ErrorSmart} from "./ErrorSmart";
import {Stack} from "@mui/material";
import {SxPropsMap} from "../../../type/general";

const SX: SxPropsMap = {
    stack: {width: "100%", height: "100%", gap: 4},
}

type Props = {
    error: any,
}

export function PageErrorBox(props: Props) {
    return (
        <Stack sx={SX.stack}>
            <PageMainBox><ErrorSmart error={props.error}/></PageMainBox>
        </Stack>
    )
}
