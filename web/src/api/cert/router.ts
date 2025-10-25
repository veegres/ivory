import {AxiosProgressEvent, AxiosRequestConfig} from "axios";

import {api} from "../api";
import {R} from "../management/type";
import {Cert, CertAddRequest, CertMap, CertType, CertUploadRequest} from "./type";

export const CertApi = {
    list: {
        key: (type?: CertType) => ["cert", "list", type],
        fn: (type?: CertType) => api.get<R<CertMap>>("/cert", {params: {type}})
            .then((response => response.data.response)),
    },
    upload: {
        key: () => ["cert", "upload"],
        fn: async (request: CertUploadRequest) => {
            const {file, type, setProgress} = request
            const formData = new FormData()
            formData.append("type", type.toString())
            formData.append("file", file)
            const config: AxiosRequestConfig = {
                headers: {"Content-Type": "multipart/form-data"},
                onUploadProgress: (progressEvent: AxiosProgressEvent) => setProgress && setProgress(progressEvent)
            }
            const response = await api.post<R<Cert>>("/cert/upload", formData, config)
            return response.data.response
        },
    },
    add: {
        key: () => ["cert", "add"],
        fn: (request: CertAddRequest) => api.post<R<Cert>>("/cert/add", request)
            .then((response) => response.data.response),
    },
    delete: {
        key: () => ["cert", "delete"],
        fn: (uuid: string) => api.delete(`/cert/${uuid}`)
            .then((response) => response.data.response),
    },
}