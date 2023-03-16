import {Grid, Skeleton} from "@mui/material";
import {instanceApi} from "../../../app/api";
import {useMutation, useQuery} from "@tanstack/react-query";
import {useTheme} from "../../../provider/ThemeProvider";
import {ErrorAlert} from "../../view/ErrorAlert";
import {useEffect, useState} from "react";
import ReactCodeMirror from "@uiw/react-codemirror";
import {json} from "@codemirror/lang-json";
import {Instance} from "../../../type/Instance";
import {TabProps} from "./Overview";
import {ClusterNoInstanceError} from "./OverviewError";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {CodeThemes} from "../../../app/utils";
import {CancelIconButton, CopyIconButton, EditIconButton, SaveIconButton} from "../../view/IconButtons";

export function OverviewConfig({info}: TabProps) {
    const theme = useTheme();
    const {defaultInstance, cluster} = info
    const [isEditable, setIsEditable] = useState(false)
    const [configState, setConfigState] = useState("")

    const {data: config, isLoading, isError, error} = useQuery(
        ["instance/config", defaultInstance.sidecar.host],
        () => instanceApi.config({...defaultInstance.sidecar, cluster: cluster.name}),
        {enabled: defaultInstance.inCluster}
    )
    const updateOptions = useMutationOptions([["instance/config", defaultInstance.sidecar.host]], () => setIsEditable(false))
    const updateConfig = useMutation(instanceApi.updateConfig, updateOptions)

    useEffect(() => setConfigState(stringify(config)), [config])

    if (!defaultInstance.inCluster) return <ClusterNoInstanceError/>
    if (isError) return <ErrorAlert error={error}/>
    if (isLoading) return <Skeleton variant={"rectangular"} height={300}/>

    const border = `1px solid ${isEditable && theme.info ? theme.info.palette.divider : "transparent"}`
    return (
        <Grid container flexWrap={"nowrap"} gap={1}>
            <Grid item flexGrow={1} sx={{border}}>
                <ReactCodeMirror
                    value={configState}
                    editable={isEditable}
                    autoFocus={isEditable}
                    basicSetup={{highlightActiveLine: false, highlightActiveLineGutter: isEditable}}
                    theme={CodeThemes[theme.mode]}
                    extensions={[json()]}
                    onChange={(value) => setConfigState(value)}
                />
            </Grid>
            <Grid item xs={"auto"} display={"flex"} flexDirection={"column"}>
                {renderUpdateButtons(defaultInstance, configState, updateConfig.isLoading, isEditable)}
                <CopyIconButton placement={"left"} size={35} onClick={handleCopyAll}/>
            </Grid>
        </Grid>
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
        if (configState) updateConfig.mutate({
            ...instance.sidecar,
            cluster: cluster.name,
            body: JSON.parse(config)
        })
    }

    function stringify(json?: any) {
        return JSON.stringify(json, null, 2)
    }
}
