import {UploadButton} from "../../view/button/UploadButton";
import {useState} from "react";
import {getErrorMessage} from "../../../app/utils";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {certApi} from "../../../app/api";
import {CertType} from "../../../type/cert";
import {AxiosProgressEvent} from "axios";

type Props = {
    type: CertType,
}

export function CertsTabUpload(props: Props) {
    const { type } = props
    const [progress, setProgress] = useState<AxiosProgressEvent>()

    const uploadOptions = useMutationOptions([["certs"]])
    const upload = useMutation(certApi.upload, uploadOptions)
    const { loading, error: uploadError } = getUploadInfo()

    return (
        <UploadButton
            accept={".crt,.key,.chain"}
            maxSize={"1MB"}
            onUpload={handleUpload}
            loading={loading}
            error={uploadError}
        />
    )

    function handleUpload(file: File) {
        upload.mutate({file, setProgress, type})
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
