import {HelpOutline} from "@mui/icons-material";
import {Alert, AlertColor, Box, Collapse, FormControl, FormLabel} from "@mui/material";
import {ReactNode, useState} from "react";

import {SxPropsMap} from "../../../app/type";

const SX: SxPropsMap = {
    box: {border: 1, borderColor: "divider", borderRadius: 1, width: "100%"},
    form: {display: "flex", flexDirection: "row", justifyContent: "space-between", alignItems: "center", height: "23px", margin: "10px 15px"},
    formAction: {position: "absolute", right: "0px"},
    formLabel: {display: "flex", gap: "4px", alignItems: "center", cursor: "pointer"},
    child: {margin: "5px 15px 10px"},
    alert: {flexGrow: 1, display: "flex", gap: 1, flexDirection: "column", margin: "0px 15px 8px"},
    desc: {height: "100%", display: "flex", gap: 1, flexDirection: "column"},
}

type Props = {
    label: string,
    renderAction: ReactNode,
    renderBody?: ReactNode,
    showBody?: boolean,
    description?: ReactNode,
    recommendation?: ReactNode,
    severity?: AlertColor,
}

export function ConfigBox(props: Props) {
    const {description, label, showBody, renderBody, severity, recommendation, renderAction} = props
    const [open, setOpen] = useState(false)
    return (
        <Box sx={SX.box}>
            {renderForm()}
            {renderInfo()}
            <Collapse in={showBody}>
                <Box sx={SX.child}>
                    {renderBody}
                </Box>
            </Collapse>
        </Box>
    )

    function renderForm() {
        return (
            <FormControl sx={SX.form}>
                <FormLabel sx={SX.formLabel} onClick={() => setOpen(!open)}>
                    <Box>{label}</Box>
                    {description && <HelpOutline fontSize={"small"}/>}
                </FormLabel>
                <Box sx={SX.formAction}>
                    {renderAction}
                </Box>
            </FormControl>
        )
    }

    function renderInfo() {
        return (
            <Collapse in={open && (!!description || !!recommendation)}>
                <Alert sx={SX.alert} icon={false} variant={"outlined"} severity={severity}>
                    <Box sx={SX.desc}>
                        {description && <Box>{description}</Box>}
                        {recommendation && (
                            <Box>
                                <Box><b>Recommendation</b></Box>
                                <Box>{recommendation}</Box>
                            </Box>
                        )}
                    </Box>
                </Alert>
            </Collapse>
        )
    }
}
