import {HelpOutlined} from "@mui/icons-material"
import {Alert, Box, Collapse, FormControl, FormLabel} from "@mui/material"
import {ReactNode, useState} from "react"

import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    box: {
        display: "flex", flexDirection: "column", padding: "5px 12px",
        border: 1, borderColor: "divider", borderRadius: 1, width: "100%"
    },
    form: {display: "flex", flexDirection: "row", justifyContent: "space-between", alignItems: "center"},
    formAction: {position: "absolute", right: "0px"},
    formLabel: {display: "flex", gap: "4px", alignItems: "center", cursor: "pointer"},
    child: {marginTop: 2},
    alert: {flexGrow: 1, display: "flex", gap: 1, flexDirection: "column", marginTop: 1},
    desc: {height: "100%", display: "flex", gap: 1, flexDirection: "column"},
}

type Props = {
    label: string,
    renderAction: ReactNode,
    renderBody?: ReactNode,
    showBody?: boolean,
    description?: ReactNode,
    recommendation?: ReactNode,
}

export function ConfigBox(props: Props) {
    const {description, label, showBody, renderBody, recommendation, renderAction} = props
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
                    {description && <HelpOutlined fontSize={"small"}/>}
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
                <Alert sx={SX.alert} icon={false} variant={"outlined"} severity={"info"}>
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
