import {Box} from "@mui/material"
import {AxiosProgressEvent} from "axios"
import {useState} from "react"

import {AlertInformative} from "../../../shared/component/box/AlertInformative"
import {List} from "../../../shared/component/box/List"
import {DownloadButton} from "../../../shared/component/button/DownloadButton"
import {UploadButton} from "../../../shared/component/button/UploadButton"
import {SxPropsMap} from "../../../shared/helper/type"
import {getErrorMessage} from "../../../shared/helper/utils"
import {useRouterExport, useRouterImport} from "../api/hook"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 2},
    export: {
        display: "flex", justifyContent: "center", padding: "10px 20px", alignItems: "center",
        width: "100%", height: "100%", flexDirection: "column", gap: 1, border: "2px dashed",
        borderRadius: "15px", minHeight: "120px", minWidth: "350px",
    },
    exportText: {fontSize: "12px"},
    wrap: {padding: "5px", height: "100%"},
}

export function ManageBackup() {
    const exp = useRouterExport()
    const imp = useRouterImport()

    const [progress, setProgress] = useState<AxiosProgressEvent>()
    const error = imp.isError ? getErrorMessage(imp.error) : undefined

    return (
        <Box sx={SX.box}>
            <AlertInformative
                title={"Backup and Restore Your Data"}
                subtitle={`
                    Ivory backs up all clusters, permissions, and manual queries. We ensure full
                    backward compatibility for backup files across all versions.
                `}
                description={<>
                    Please note that importing data may overwrite existing clusters or user
                    permissions, which could lead to data loss in the event of a conflict. For
                    security reasons, credential are not included in
                    backups and must be re-entered manually. While direct database compatibility
                    between versions is not guaranteed, Ivory backups are designed for universal
                    compatibility. Use this tool to safely migrate data between versions. Detailed
                    version compatibility can be found on the Security page.
                </>}
            />
            <List name={"Export backup"}>
               <DownloadButton onClick={exp.mutate}/>
            </List>
            <List name={"Import backup"}>
                <UploadButton
                    accept={".bak"}
                    maxSize={"1GB"}
                    loading={{loading: imp.isPending, loaded: progress?.loaded, total: progress?.total}}
                    error={error}
                    onUpload={handleOnUpload}
                />
            </List>
        </Box>
    )

    function handleOnUpload(file: File) {
        imp.mutate({file, setProgress})
    }
}
