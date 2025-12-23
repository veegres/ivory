import {FileUploadOutlined} from "@mui/icons-material"
import {Box, Button, CircularProgress} from "@mui/material"
import {ChangeEvent, DragEvent, useState} from "react"

import {SxPropsMap} from "../../../app/type"

const SX: SxPropsMap = {
    box: {padding: "5px", height: "100%"},
    upload: {
        minHeight: "120px", minWidth: "350px", width: "100%", height: "100%", display: "flex",
        alignItems: "center", flexDirection: "column", justifyContent: "space-evenly",
        padding: "10px 20px", border: "2px dashed", borderRadius: "15px", gap: 1,
    },
    error: {color: "error.main", whiteSpace: "pre-wrap"},
    label: {textTransform: "uppercase", fontSize: "12px"},
    active: (theme) => ({
        pointerEvents: "none",
        boxShadow: `inset 0px 0px 50px 50px ${color[theme.palette.mode]}`,
    }),
    button: (theme) => ({boxShadow: `inset 0px 0px 50px 5px ${color[theme.palette.mode]}`}),
}

const color = {
    light: "rgba(0,0,0,0.1)",
    dark: "rgba(255,255,255,0.1)",
}

type Props = {
    accept: string,
    maxSize?: string,
    loading: {
        loading: boolean,
        loaded?: number,
        total?: number,
    },
    error?: string,
    onUpload: (file: File) => void,
}

export function UploadButton(props: Props) {
    const {accept, maxSize, loading: {loading, loaded, total}, error} = props
    const [active, setActive] = useState(false)

    return (
        <Box
            sx={SX.box}
            onDragEnter={handleDragEnter}
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
            onDrop={handleDrop}
        >
            {renderBody()}
        </Box>
    )

    function renderBody() {
        if (loading) return renderLoading()
        if (active) return renderDrop()
        return renderButton()
    }

    function renderButton() {
        return (
            <Button sx={[SX.upload, SX.button]} color={"inherit"} component={"label"}>
                <input hidden accept={accept} multiple type={"file"} onChange={handleChange}/>
                <FileUploadOutlined fontSize={"large"}/>
                <Box sx={SX.label}>Drag and drop a file or click to browse</Box>
                {maxSize && !error && <Box textTransform={"none"} fontSize={"small"}>Maximum file size is {maxSize}</Box>}
                {error && <Box sx={SX.error} textTransform={"none"} fontSize={"small"}>{error}</Box>}
            </Button>
        )
    }

    function renderDrop() {
        return (
            <Box sx={[SX.upload, SX.active]}>
                <FileUploadOutlined fontSize={"large"}/>
                <Box sx={SX.label}>Drop the file here</Box>
            </Box>
        )
    }

    function renderLoading() {
        const value = loaded && total ? Math.floor(loaded * 100 / total) : undefined
        const variant = value ? "determinate" : "indeterminate"
        return (
            <Box sx={[SX.upload, SX.active]}>
                <CircularProgress variant={variant} value={value}/>
                <Box fontSize={"small"}>Loading {value === 100 ? "is finished" : `in process ${value}%`}</Box>
            </Box>
        )
    }

    function handleChange(event: ChangeEvent<HTMLInputElement>) {
        handleUpload(event.target.files?.item(0))
        // clean value to be able to upload its one more time
        event.target.value = ""
    }

    function handleDragEnter() {
        if (!active) setActive(true)

    }

    function handleDragOver(event: DragEvent<HTMLInputElement>) {
        // allow drop function to be fired
        event.preventDefault()
    }

    function handleDragLeave() {
        if (active) setActive(false)
    }

    function handleDrop(event: DragEvent<HTMLInputElement>) {
        event.preventDefault()
        handleUpload(event.dataTransfer.files.item(0))
        setActive(false)
    }

    function handleUpload(file: File | null | undefined) {
        if (file) props.onUpload(file)
    }
}
