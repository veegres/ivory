import {UploadButton} from "../../../view/button/UploadButton";
import {useState} from "react";
import {getErrorMessage} from "../../../../app/utils";
import {CertType} from "../../../../api/cert/type";
import {AxiosProgressEvent} from "axios";
import {useRouterCertUpload} from "../../../../api/cert/hook";

type Props = {
    type: CertType,
}

export function CertsTabUpload(props: Props) {
    const {type} = props
    const [progress, setProgress] = useState<AxiosProgressEvent>()
    const upload = useRouterCertUpload(type)
    const {loading, error: uploadError} = getUploadInfo()

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
            isLoading: upload.isPending,
            loaded: progress?.loaded,
            total: progress?.total,
        }
        return {loading, error}
    }
}
