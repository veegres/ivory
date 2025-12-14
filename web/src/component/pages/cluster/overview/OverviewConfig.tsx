import {json} from "@codemirror/lang-json"
import {Box, Skeleton} from "@mui/material"
import ReactCodeMirror from "@uiw/react-codemirror"
import {useEffect, useState} from "react"

import {ActiveCluster} from "../../../../api/cluster/type"
import {useRouterInstanceConfig, useRouterInstanceConfigUpdate} from "../../../../api/instance/hook"
import {Permission} from "../../../../api/permission/type"
import {SxPropsMap} from "../../../../app/type"
import {CodeThemes, getSidecarConnection} from "../../../../app/utils"
import {useSettings} from "../../../../provider/AppProvider"
import {useSnackbar} from "../../../../provider/SnackbarProvider"
import {ErrorSmart} from "../../../view/box/ErrorSmart"
import {CancelIconButton, CopyIconButton, EditIconButton, SaveIconButton} from "../../../view/button/IconButtons"
import {Access} from "../../../widgets/access/Access"
import {ClusterNoInstanceError} from "./OverviewError"

const SX: SxPropsMap = {
    box: {display: "flex", flexWrap: "nowrap", gap: 1, height: "100%"},
    input: {flexGrow: 1, borderWidth: "1px", borderStyle: "solid", overflowX: "auto", ">div": {height: "100%"}},
    buttons: {display: "flex", flexDirection: "column"},
}

type Props = {
    info: ActiveCluster
}

export function OverviewConfig(props: Props) {
    const {info} = props
    const settings = useSettings()
    const snackbar = useSnackbar()
    const {defaultInstance, cluster} = info
    const [isEditable, setIsEditable] = useState(false)
    const [configState, setConfigState] = useState("")
    const {sidecar} = defaultInstance
    const connection = getSidecarConnection(cluster, sidecar)

    const config = useRouterInstanceConfig(connection, defaultInstance.inCluster)
    const updateConfig = useRouterInstanceConfigUpdate(sidecar, () => setIsEditable(false))

    const {data, isPending, isError, error} = config

    useEffect(() => setConfigState(stringify(data)), [data])

    if (!defaultInstance.inCluster) return <ClusterNoInstanceError/>
    if (isError) return <ErrorSmart error={error}/>
    if (isPending) return <Skeleton variant={"rectangular"} height={300}/>

    return (
        <Box sx={SX.box}>
            <Box sx={SX.input} borderColor={isEditable ? "divider" : "transparent"}>
                <ReactCodeMirror
                    height={"100%"}
                    width={"100%"}
                    value={configState}
                    editable={isEditable}
                    autoFocus={isEditable}
                    basicSetup={{highlightActiveLine: false, highlightActiveLineGutter: isEditable, highlightSelectionMatches: false}}
                    theme={CodeThemes[settings.theme]}
                    extensions={[json()]}
                    onChange={(value) => setConfigState(value)}
                />
            </Box>
            <Box sx={SX.buttons}>
                <Access permission={Permission.ManageInstanceConfigUpdate}>
                    {renderUpdateButtons()}
                </Access>
                <CopyIconButton placement={"left"} size={35} onClick={handleCopyAll}/>
            </Box>
        </Box>
    )

    function renderUpdateButtons() {
        if (!isEditable) return <EditIconButton placement={"left"} size={35} onClick={() => setIsEditable(true)}/>

        return (
            <>
                <CancelIconButton placement={"left"} size={35} disabled={updateConfig.isPending} onClick={handleCancel}/>
                <SaveIconButton placement={"left"} size={35} loading={updateConfig.isPending} onClick={handleUpdate}/>
            </>
        )
    }

    function handleCopyAll() {
        const currentConfig = configState ? configState : stringify(data)
        navigator.clipboard.writeText(currentConfig).then(() => {
            snackbar("Config copied to clipboard!", "info")
        })
    }

    function handleCancel() {
        setIsEditable(false)
        setConfigState(stringify(data))
    }

    function handleUpdate() {
        if (configState) {
            updateConfig.mutate({...connection, body: JSON.parse(configState)})
        }
    }

    function stringify(json?: any) {
        return JSON.stringify(json, null, 2)
    }
}
