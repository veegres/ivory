import {Box, Skeleton} from "@mui/material";
import {useAppearance} from "../../../../provider/AppearanceProvider";
import {ErrorSmart} from "../../../view/box/ErrorSmart";
import {useEffect, useState} from "react";
import ReactCodeMirror from "@uiw/react-codemirror";
import {json} from "@codemirror/lang-json";
import {InstanceRequest, InstanceWeb} from "../../../../type/instance";
import {ClusterNoInstanceError} from "./OverviewError";
import {CodeThemes} from "../../../../app/utils";
import {CancelIconButton, CopyIconButton, EditIconButton, SaveIconButton} from "../../../view/button/IconButtons";
import {SxPropsMap} from "../../../../type/general";
import {ActiveCluster} from "../../../../type/cluster";
import {useRouterInstanceConfig, useRouterInstanceConfigUpdate} from "../../../../router/instance";
import {useSnackbar} from "../../../../provider/SnackbarProvider";

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
    const appearance = useAppearance();
    const snackbar = useSnackbar()
    const {defaultInstance, cluster} = info
    const [isEditable, setIsEditable] = useState(false)
    const [configState, setConfigState] = useState("")
    const {sidecar} = defaultInstance
    const requestBody: InstanceRequest = {sidecar, credentialId: cluster.credentials.patroniId, certs: cluster.certs}

    const {data: config, isPending, isError, error} = useRouterInstanceConfig(requestBody, defaultInstance.inCluster)
    const updateConfig = useRouterInstanceConfigUpdate(sidecar, () => setIsEditable(false))

    useEffect(() => setConfigState(stringify(config)), [config])

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
                    theme={CodeThemes[appearance.theme]}
                    extensions={[json()]}
                    onChange={(value) => setConfigState(value)}
                />
            </Box>
            <Box sx={SX.buttons}>
                {renderUpdateButtons(defaultInstance, configState, updateConfig.isPending, isEditable)}
                <CopyIconButton placement={"left"} size={35} onClick={handleCopyAll}/>
            </Box>
        </Box>
    )

    function renderUpdateButtons(instance: InstanceWeb, configState: string, isLoading: boolean, isEditable: boolean) {
        if (!isEditable) return <EditIconButton placement={"left"} size={35} onClick={() => setIsEditable(true)}/>

        return (
            <>
                <CancelIconButton placement={"left"} size={35} disabled={isLoading} onClick={handleCancel}/>
                <SaveIconButton placement={"left"} size={35} loading={isLoading} onClick={() => handleUpdate(instance, configState)}/>
            </>
        )
    }

    function handleCopyAll() {
        const currentConfig = configState ? configState : stringify(config)
        navigator.clipboard.writeText(currentConfig).then(() => {
            snackbar("Config copied to clipboard!", "info")
        })
    }

    function handleCancel() {
        setIsEditable(false)
        setConfigState(stringify(config))
    }

    function handleUpdate(instance: InstanceWeb, config: string) {
        if (configState) {
            updateConfig.mutate({
                sidecar: instance.sidecar,
                credentialId: cluster.credentials.patroniId,
                certs: cluster.certs,
                body: JSON.parse(config)
            })
        }
    }

    function stringify(json?: any) {
        return JSON.stringify(json, null, 2)
    }
}
