import {Box, Skeleton} from "@mui/material";
import {instanceApi} from "../../../app/api";
import {useMutation, useQuery} from "@tanstack/react-query";
import {useAppearance} from "../../../provider/AppearanceProvider";
import {ErrorSmart} from "../../view/box/ErrorSmart";
import {useEffect, useState} from "react";
import ReactCodeMirror from "@uiw/react-codemirror";
import {json} from "@codemirror/lang-json";
import {Instance, InstanceRequest} from "../../../type/Instance";
import {TabProps} from "./Overview";
import {ClusterNoInstanceError} from "./OverviewError";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {CodeThemes} from "../../../app/utils";
import {CancelIconButton, CopyIconButton, EditIconButton, SaveIconButton} from "../../view/button/IconButtons";
import {SxPropsMap} from "../../../type/common";

const SX: SxPropsMap = {
    box: {display: "flex", flexWrap: "nowrap", gap: 1},
    input: {flexGrow: 1, borderWidth: "1px", borderStyle: "solid", overflowX: "auto"},
    buttons: {display: "flex", flexDirection: "column"},
}

export function OverviewConfig({info}: TabProps) {
    const theme = useAppearance();
    const {defaultInstance, cluster} = info
    const [isEditable, setIsEditable] = useState(false)
    const [configState, setConfigState] = useState("")
    const {sidecar} = defaultInstance
    const requestBody: InstanceRequest = {...sidecar, credentialId: cluster.credentials.patroniId, certs: cluster.certs}

    const {data: config, isLoading, isError, error} = useQuery(
        ["instance/config", sidecar.host, sidecar.port],
        () => instanceApi.config(requestBody),
        {enabled: defaultInstance.inCluster}
    )
    const updateOptions = useMutationOptions([["instance/config", sidecar.host, sidecar.port]], () => setIsEditable(false))
    const updateConfig = useMutation(instanceApi.updateConfig, updateOptions)

    useEffect(() => setConfigState(stringify(config)), [config])

    if (!defaultInstance.inCluster) return <ClusterNoInstanceError/>
    if (isError) return <ErrorSmart error={error}/>
    if (isLoading) return <Skeleton variant={"rectangular"} height={300}/>

    return (
        <Box sx={SX.box}>
            <Box sx={SX.input} borderColor={isEditable ? "divider" : "transparent"}>
                <ReactCodeMirror
                    value={configState}
                    editable={isEditable}
                    autoFocus={isEditable}
                    basicSetup={{highlightActiveLine: false, highlightActiveLineGutter: isEditable}}
                    theme={CodeThemes[theme.state.mode]}
                    extensions={[json()]}
                    onChange={(value) => setConfigState(value)}
                />
            </Box>
            <Box sx={SX.buttons}>
                {renderUpdateButtons(defaultInstance, configState, updateConfig.isLoading, isEditable)}
                <CopyIconButton placement={"left"} size={35} onClick={handleCopyAll}/>
            </Box>
        </Box>
    )

    function renderUpdateButtons(instance: Instance, configState: string, isLoading: boolean, isEditable: boolean) {
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
        return navigator.clipboard.writeText(currentConfig)
    }

    function handleCancel() {
        setIsEditable(false)
        setConfigState(stringify(config))
    }

    function handleUpdate(instance: Instance, config: string) {
        if (configState) {
            const request: InstanceRequest = {
                ...instance.sidecar,
                credentialId: cluster.credentials.patroniId,
                certs: cluster.certs,
                body: JSON.parse(config)
            }
            updateConfig.mutate(request)
        }
    }

    function stringify(json?: any) {
        return JSON.stringify(json, null, 2)
    }
}
