import {Box, Skeleton} from "@mui/material";
import {instanceApi} from "../../../app/api";
import {useMutation, useQuery} from "@tanstack/react-query";
import {useAppearance} from "../../../provider/AppearanceProvider";
import {ErrorSmart} from "../../view/box/ErrorSmart";
import {useEffect, useState} from "react";
import ReactCodeMirror from "@uiw/react-codemirror";
import {json} from "@codemirror/lang-json";
import {InstanceRequest, InstanceWeb} from "../../../type/instance";
import {TabProps} from "./Overview";
import {ClusterNoInstanceError} from "./OverviewError";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {CodeThemes} from "../../../app/utils";
import {CancelIconButton, CopyIconButton, EditIconButton, SaveIconButton} from "../../view/button/IconButtons";
import {SxPropsMap} from "../../../type/common";
import {useSnackbar} from "notistack";

const SX: SxPropsMap = {
    box: {display: "flex", flexWrap: "nowrap", gap: 1, height: "100%"},
    input: {flexGrow: 1, borderWidth: "1px", borderStyle: "solid", overflowX: "auto", ">div": {height: "100%"}},
    buttons: {display: "flex", flexDirection: "column"},
}

export function OverviewConfig({info}: TabProps) {
    const appearance = useAppearance();
    const {enqueueSnackbar} = useSnackbar()
    const {defaultInstance, cluster} = info
    const [isEditable, setIsEditable] = useState(false)
    const [configState, setConfigState] = useState("")
    const {sidecar} = defaultInstance
    const requestBody: InstanceRequest = {...sidecar, credentialId: cluster.credentials.patroniId, certs: cluster.certs}

    const {data: config, isPending, isError, error} = useQuery({
        queryKey: ["instance/config", sidecar.host, sidecar.port],
        queryFn: () => instanceApi.config(requestBody),
        enabled: defaultInstance.inCluster,
    })
    const updateOptions = useMutationOptions([["instance/config", sidecar.host, sidecar.port]], () => setIsEditable(false))
    const updateConfig = useMutation({mutationFn: instanceApi.updateConfig, ...updateOptions})

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
        return navigator.clipboard.writeText(currentConfig)
    }

    function handleCancel() {
        setIsEditable(false)
        setConfigState(stringify(config))
    }

    function handleUpdate(instance: InstanceWeb, config: string) {
        if (configState) {
            try {
                const request: InstanceRequest = {
                    ...instance.sidecar,
                    credentialId: cluster.credentials.patroniId,
                    certs: cluster.certs,
                    body: JSON.parse(config)
                }
                updateConfig.mutate(request)
            } catch (e: any) {
                enqueueSnackbar(e?.message ?? "Unknown error", {variant: "error"})
            }
        }
    }

    function stringify(json?: any) {
        return JSON.stringify(json, null, 2)
    }
}
