import {Box, CircularProgress, Grid, IconButton, Skeleton, Tooltip} from "@mui/material";
import {instanceApi} from "../../../app/api";
import {useMutation, useQuery} from "@tanstack/react-query";
import {useTheme} from "../../../provider/ThemeProvider";
import {ErrorAlert} from "../../view/ErrorAlert";
import {ReactElement, useEffect, useState} from "react";
import ReactCodeMirror from "@uiw/react-codemirror";
import {json} from "@codemirror/lang-json";
import {Cancel, CopyAll, Edit, SaveAlt} from "@mui/icons-material";
import {InstanceLocal} from "../../../app/types";
import {TabProps} from "./Overview";
import {ClusterNoInstanceError} from "./OverviewError";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {darculaInit} from "@uiw/codemirror-theme-darcula";
import {materialLightInit} from "@uiw/codemirror-theme-material";

const options = {settings: {background: "transparent", gutterBackground: "transparent"}}
const themes = {
    dark: darculaInit(options),
    light: materialLightInit(options),
}

export function OverviewConfig({info}: TabProps) {
    const theme = useTheme();
    const {instance, cluster} = info
    const [isEditable, setIsEditable] = useState(false)
    const [configState, setConfigState] = useState("")

    const {data: config, isLoading, isError, error} = useQuery(
        ["instance/config", instance.sidecar.host],
        () => instanceApi.config({...instance.sidecar, cluster: cluster.name}),
        {enabled: instance.inCluster}
    )
    const updateOptions = useMutationOptions(["instance/config", instance.sidecar.host], () => setIsEditable(false))
    const updateConfig = useMutation(instanceApi.updateConfig, updateOptions)

    useEffect(() => setConfigState(stringify(config)), [config])

    if (!instance.inCluster) return <ClusterNoInstanceError/>
    if (isError) return <ErrorAlert error={error}/>
    if (isLoading) return <Skeleton variant={"rectangular"} height={300}/>

    const border = `1px solid ${isEditable && theme.info ? theme.info.palette.divider : "transparent"}`
    return (
        <Grid container flexWrap={"nowrap"}>
            <Grid item flexGrow={1} sx={{border}}>
                <ReactCodeMirror
                    value={configState}
                    editable={isEditable}
                    theme={themes[theme.mode]}
                    extensions={[json()]}
                    onChange={(value) => setConfigState(value)}
                />
            </Grid>
            <Grid item xs={"auto"} display={"flex"} flexDirection={"column"}>
                {renderUpdateButtons(instance, configState, updateConfig.isLoading, isEditable)}
                {renderButton(<CopyAll/>, "Copy", handleCopyAll)}
            </Grid>
        </Grid>
    )

    function renderUpdateButtons(instance: InstanceLocal, configState: string, isLoading: boolean, isEditable: boolean) {
        if (isLoading) return renderButton(<CircularProgress size={25}/>, "Saving", undefined, true)
        if (!isEditable) return renderButton(<Edit/>, "Edit", () => setIsEditable(true))

        return (
            <>
                {renderButton(<Cancel/>, "Cancel", handleCancel)}
                {renderButton(<SaveAlt/>, "Save", () => handleUpdate(instance, configState))}
            </>
        )
    }

    function renderButton(icon: ReactElement, tooltip: string, onClick = () => {}, disabled = false) {
        return (
            <Tooltip title={tooltip} placement={"left"} arrow>
                <Box component={"span"}>
                    <IconButton onClick={onClick} disabled={disabled}>{icon}</IconButton>
                </Box>
            </Tooltip>
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

    function handleUpdate(instance: InstanceLocal, config: string) {
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
