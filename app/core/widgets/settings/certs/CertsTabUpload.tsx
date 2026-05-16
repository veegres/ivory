import {AxiosProgressEvent} from "axios"
import {useState} from "react"

import {useRouterCertUpload} from "../../../../features/cert/hook"
import {CertType} from "../../../../features/cert/type"
import {UploadButton} from "../../../../shared/component/button/UploadButton"
import {getErrorMessage} from "../../../../shared/helper/utils"

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
            loading: upload.isPending,
            loaded: progress?.loaded,
            total: progress?.total,
        }
        return {loading, error}
    }
}
