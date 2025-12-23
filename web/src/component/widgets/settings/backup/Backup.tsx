import {AxiosProgressEvent} from "axios"
import {useState} from "react"

import {useRouterExport, useRouterImport} from "../../../../api/management/hook"
import {Permission} from "../../../../api/permission/type"
import {SxPropsMap} from "../../../../app/type"
import {getErrorMessage} from "../../../../app/utils"
import {AlertInformative} from "../../../view/box/AlertInformative"
import {List} from "../../../view/box/List"
import {DownloadButton} from "../../../view/button/DownloadButton"
import {UploadButton} from "../../../view/button/UploadButton"
import {Access} from "../../access/Access"
import {MenuWrapperScroll} from "../menu/MenuWrapperScroll"

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

export function Backup() {
    const exp = useRouterExport()
    const imp = useRouterImport()

    const [progress, setProgress] = useState<AxiosProgressEvent>()
    const error = imp.isError ? getErrorMessage(imp.error) : undefined

    return (
        <MenuWrapperScroll sx={SX.box}>
            <AlertInformative
                title={"Backup and Restore Your Data"}
                subtitle={`
                    Ivory backs up all clusters, permissions, and manual queries. We ensure full
                     backward compatibility for backup files across all versions.
                `}
                description={<>
                    Please note that importing data may overwrite existing clusters or user
                    permissions, which could lead to data loss in the event of a conflict. For
                    security reasons, passwords and credential states are not included in
                    backups and must be re-entered manually. While direct database compatibility
                    between versions is not guaranteed, Ivory backups are designed for universal
                    compatibility. Use this tool to safely migrate data between versions. Detailed
                    version compatibility can be found on the Security page.
                </>}
            />
            <Access permission={Permission.ManageManagementExport}>
                <List name={"Export backup"}>
                   <DownloadButton onClick={exp.mutate}/>
                </List>
            </Access>
            <Access permission={Permission.ManageManagementImport}>
                <List name={"Import backup"}>
                    <UploadButton
                        accept={".bak"}
                        maxSize={"1GB"}
                        loading={{loading: imp.isPending, loaded: progress?.loaded, total: progress?.total}}
                        error={error}
                        onUpload={handleOnUpload}
                    />
                </List>
            </Access>
        </MenuWrapperScroll>
    )

    function handleOnUpload(file: File) {
        imp.mutate({file, setProgress})
    }
}
