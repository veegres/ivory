import {Box, CircularProgress, Grid, IconButton, Skeleton, Tooltip} from "@mui/material";
import {nodeApi} from "../../app/api";
import {useMutation, useQuery} from "@tanstack/react-query";
import {useTheme} from "../../provider/ThemeProvider";
import {ErrorAlert} from "../view/ErrorAlert";
import React, {ReactElement, useEffect, useState} from "react";
import ReactCodeMirror from "@uiw/react-codemirror";
import {json} from "@codemirror/lang-json";
import {Cancel, CopyAll, Edit, SaveAlt} from "@mui/icons-material";
import {oneDarkHighlightStyle} from "@codemirror/theme-one-dark";
import {syntaxHighlighting, defaultHighlightStyle} from "@codemirror/language";
import {Instance} from "../../app/types";
import {TabProps} from "./Cluster";
import {ClusterNoInstanceError} from "./ClusterError";
import {EditorView} from "@codemirror/view";
import {useMutationOptions} from "../../app/hooks";

const highlightExtension = {
    dark: syntaxHighlighting(oneDarkHighlightStyle),
    light: syntaxHighlighting(defaultHighlightStyle)
}

export function ClusterConfig({info}: TabProps) {
    const theme = useTheme();
    const { instance } = info
    const [isEditable, setIsEditable] = useState(false)
    const [configState, setConfigState] = useState('')

    const {data: config, isLoading, isError, error} = useQuery(
        ["node/config", instance.api_domain],
        () => nodeApi.config(instance.api_domain),
        { enabled: instance.inCluster }
    )
    const updateOptions = useMutationOptions(["node/config", instance.api_domain])
    const updateConfig = useMutation(nodeApi.updateConfig, updateOptions)

    useEffect(() => setConfigState(stringify(config)), [config])

    if (isError) return <ErrorAlert error={error}/>
    if (isLoading) return <Skeleton variant="rectangular" height={300}/>

    const isDark = theme.mode === "dark"
    const codeMirrorTheme = EditorView.theme({}, {dark: isDark})

    if (!instance.inCluster) return <ClusterNoInstanceError />

    return (
        <Grid container flexWrap={"nowrap"}>
            <Grid item flexGrow={1}>
                <ReactCodeMirror
                    value={configState}
                    theme={codeMirrorTheme}
                    editable={isEditable}
                    extensions={[json(), isDark ? highlightExtension.dark : highlightExtension.light]}
                    onChange={(value) => setConfigState(value)}
               />
            </Grid>
            <Grid item xs={"auto"} display={"flex"} flexDirection={"column"}>
                {renderUpdateButtons(instance, configState, updateConfig.isLoading, isEditable)}
                {renderButton(<CopyAll/>, "Copy", handleCopyAll)}
            </Grid>
        </Grid>
    )

    function renderUpdateButtons(instance: Instance, configState: string, isLoading: boolean, isEditable: boolean) {
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
            <Tooltip title={tooltip} placement="left" arrow>
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

    function handleUpdate(node: Instance, config: string) {
        setIsEditable(false)
        if (configState) updateConfig.mutate({node: node.api_domain, config})
    }

    function stringify(json?: any) {
        return JSON.stringify(json, null, 2)
    }
}
