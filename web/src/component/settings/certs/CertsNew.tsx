import {
    Alert,
    Box,
    TextField,
    ToggleButton,
    ToggleButtonGroup,
    Tooltip
} from "@mui/material";
import {UploadButton} from "../../view/UploadButton";
import React, {useState} from "react";
import {useMutation} from "@tanstack/react-query";
import {certApi} from "../../../app/api";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {getErrorMessage} from "../../../app/utils";
import {FilePresent, UploadFile} from "@mui/icons-material";
import {SaveIconButton} from "../../view/IconButtons";

const SX = {
    box: { display: "flex", height: "120px" },
    pathTab: { display: "flex", flexDirection: "column", padding: "5px", justifyContent: "space-between", height: "100%" },
    pathTabForm: { display: "flex", alignItems: "center", gap: 1 },
    pathTabField: { flexGrow: 1 },
    pathTabAlert: { padding: "0 14px", borderRadius: "15px" },
    group: { minHeight: "25px", padding: "5px" },
    button: { flexGrow: 1, borderRadius: "15px" },
    border: { borderRadius: "15px" },
}

enum Tabs {
    UPLOAD,
    PATH,
}

export function CertsNew() {
    const [tab, setTab] = useState(Tabs.PATH)
    const [progress, setProgress] = useState<ProgressEvent>()

    const uploadOptions = useMutationOptions(["certs"])
    const upload = useMutation(certApi.upload, uploadOptions)
    const { loading, error: uploadError } = getUploadInfo()

    return (
        <Box sx={SX.box}>
            <ToggleButtonGroup sx={SX.group} orientation={"vertical"} size={"small"} exclusive value={tab}>
                <ToggleButton sx={SX.button} onClick={() => setTab(Tabs.UPLOAD)} value={Tabs.UPLOAD}>
                    <Tooltip placement={"right"} title={"Cert Upload"}><UploadFile /></Tooltip>
                </ToggleButton>
                <ToggleButton sx={SX.button} onClick={() => setTab(Tabs.PATH)} value={Tabs.PATH}>
                    <Tooltip placement={"right"} title={"Cert Path"}><FilePresent /></Tooltip>
                </ToggleButton>
            </ToggleButtonGroup>
            <Box sx={{ flexGrow: 1 }}>{renderTab()}</Box>
        </Box>
    )

    function renderTab() {
        switch (tab) {
            case Tabs.UPLOAD: return renderUploadPath()
            case Tabs.PATH: return renderPathTab()
        }
    }

    function renderUploadPath() {
        return (
            <UploadButton
                accept={".crt,.key"}
                maxSize={"1MB"}
                onUpload={handleUpload}
                loading={loading}
                error={uploadError}
            />
        )
    }

    function renderPathTab() {
        return (
            <Box sx={SX.pathTab}>
                <Alert sx={SX.pathTabAlert} variant={"outlined"} severity={"info"} icon={false}>
                    Ivory can work with files inside the container.
                    Keep in mind that you need to mount the folder to the container for it to work.
                </Alert>
                <Box sx={SX.pathTabForm}>
                    <TextField sx={SX.pathTabField} size={"small"} label={"Path to The File"} InputProps={{sx: SX.border}} />
                    <SaveIconButton loading={false} onClick={() => {}} />
                </Box>
            </Box>
        )
    }

    function handleUpload(file: File) {
        upload.mutate({file, setProgress})
    }

    function getUploadInfo() {
        const error = upload.isError ? getErrorMessage(upload.error) : undefined
        const loading = {
            isLoading: upload.isLoading,
            loaded: progress?.loaded,
            total: progress?.total,
        }
        return { loading, error }
    }
}
