import {Alert, Box, TextField} from "@mui/material";
import {useState} from "react";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {CertApi} from "../../../app/api";
import {CheckCircle} from "@mui/icons-material";
import {LoadingButton} from "@mui/lab";
import {SxPropsMap} from "../../../type/common";
import {CertType} from "../../../type/cert";

const SX: SxPropsMap = {
    box: { display: "flex", flexDirection: "column", padding: "5px", justifyContent: "space-between", height: "100%", gap: "12px" },
    form: { display: "flex", alignItems: "center", gap: 1 },
    textField: { flexGrow: 1 },
    alert: { padding: "0 14px", borderRadius: "15px" },
    border: { borderRadius: "15px" },
}

type Props = {
    type: CertType,
}

export function CertsTabPath(props: Props) {
    const { type } = props
    const [path, setPath] = useState("")

    const addOptions = useMutationOptions([["certs"]], () => setPath(""))
    const add = useMutation({mutationFn: CertApi.add, ...addOptions})

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
                    InputProps={{sx: SX.border}}
                    value={path}
                    onChange={(e) => setPath(e.target.value)}
                />
                <LoadingButton
                    sx={SX.border}
                    variant={"outlined"}
                    loading={add.isPending}
                    onClick={() => add.mutate({ path, type })}
                    disabled={!path}
                    size={"large"}
                >
                    <CheckCircle/>
                </LoadingButton>
            </Box>
        </Box>
    )
}
