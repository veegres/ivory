import {CheckCircle} from "@mui/icons-material";
import {Alert, Box, Button, TextField} from "@mui/material";
import {useState} from "react";

import {useRouterCertAdd} from "../../../../api/cert/hook";
import {CertType} from "../../../../api/cert/type";
import {SxPropsMap} from "../../../../app/type";

const SX: SxPropsMap = {
    box: {
        display: "flex", flexDirection: "column", padding: "5px",
        justifyContent: "space-between", height: "100%", gap: "12px",
    },
    form: {display: "flex", alignItems: "center", gap: 1},
    textField: {flexGrow: 1},
    alert: {padding: "0 14px", borderRadius: "15px"},
    border: {borderRadius: "15px"},
}

type Props = {
    type: CertType,
}

export function CertsTabPath(props: Props) {
    const {type} = props
    const [path, setPath] = useState("")
    const add = useRouterCertAdd(type, () => setPath(""))

    return (
        <Box sx={SX.box}>
            <Alert sx={SX.alert} variant={"outlined"} severity={"info"} icon={false}>
                Ivory can work with files inside the container.
                Keep in mind that you need to mount the folder to the container for it to work.
            </Alert>
            <Box sx={SX.form}>
                <TextField
                    sx={SX.textField}
                    size={"small"}
                    label={"Path to The File"}
                    value={path}
                    onChange={(e) => setPath(e.target.value)}
                    slotProps={{input: {sx: SX.border}}}
                />
                <Button
                    sx={SX.border}
                    variant={"outlined"}
                    loading={add.isPending}
                    onClick={() => add.mutate({path, type})}
                    disabled={!path}
                    size={"large"}
                >
                    <CheckCircle/>
                </Button>
            </Box>
        </Box>
    )
}
