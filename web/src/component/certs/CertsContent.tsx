import React, {useState} from "react";
import {Box, Collapse} from "@mui/material";
import {useMutation, useQuery} from "@tanstack/react-query";
import {certApi} from "../../app/api";
import {LinearProgressStateful} from "../view/LinearProgressStateful";
import {InfoAlert} from "../view/InfoAlert";
import scroll from "../../style/scroll.module.css";
import {TransitionGroup} from "react-transition-group";
import {ErrorAlert} from "../view/ErrorAlert";
import {UploadButton} from "../view/UploadButton";
import {CertsItem} from "./CertsItem";
import {getErrorMessage} from "../../app/utils";
import {useMutationOptions} from "../../app/hooks";

const SX = {
    progress: { margin: "10px 0" },
    list: { maxHeight: "500px", overflowY: "auto" },
}

export function CertsContent() {
    const [progress, setProgress] = useState<ProgressEvent>()

    const { data: certs, isError, error, isFetching } = useQuery(["certs"], certApi.list)
    const uploadOptions = useMutationOptions(["certs"])
    const upload = useMutation(certApi.upload, uploadOptions)

    if (isError) return <ErrorAlert error={error}/>
    const { loading, error: uploadError } = getUploadInfo()

    return (
        <Box>
            <UploadButton accept={".crt"} maxSize={"1MB"} onUpload={handleUpload} loading={loading} error={uploadError} />
            <LinearProgressStateful sx={SX.progress} color={"inherit"} isFetching={isFetching} line />
            {renderList()}
        </Box>
    )

    function renderList() {
        const list = certs ?? []
        if (list.length === 0) return <InfoAlert text={"There is no certs yet"}/>

        return (
            <Box sx={SX.list} className={scroll.tiny}>
                <TransitionGroup>
                    {list.map((cert) => (
                        <Collapse key={cert.fileId}>
                            <CertsItem cert={cert}/>
                        </Collapse>
                    ))}
                </TransitionGroup>
            </Box>
        )
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

    function handleUpload(file: File) {
        upload.mutate({file, setProgress})
    }
}
