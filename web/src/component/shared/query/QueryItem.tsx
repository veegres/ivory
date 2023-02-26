import {Query, SxPropsMap} from "../../../app/types";
import {Box, Collapse, Divider, Paper} from "@mui/material";
import {CancelIconButton, EditIconButton, PlayIconButton, RestoreIconButton} from "../../view/IconButtons";
import {useTheme} from "../../../provider/ThemeProvider";
import {useState} from "react";
import {QueryItemInfo} from "./QueryItemInfo";
import {QueryItemEdit} from "./QueryItemEdit";
import {QueryItemPlay} from "./QueryItemPlay";
import {InfoAlert} from "../../view/InfoAlert";
import {QueryItemRestore} from "./QueryItemRestore";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", fontSize: "15px"},
    head: {display: "flex", padding: "5px 15px"},
    body: {padding: "8px 15px", fontSize: "13px"},
    title: {flexGrow: 1, display: "flex", alignItems: "center", cursor: "pointer", gap: 1},
    name: {fontWeight: "bold"},
    creation: {fontSize: "12px", fontFamily: "monospace"},
    buttons: {display: "flex", alignItems: "center"},
    type: {padding: "0 8px", cursor: "pointer"},
}

enum BodyType {INFO, EDIT, RESTORE, PLAY}

type Props = {
    id: string,
    query: Query,
}

export function QueryItem(props: Props) {
    const {id, query} = props
    const {info} = useTheme()
    const [open, setOpen] = useState(false)
    const [body, setBody] = useState<BodyType>(BodyType.INFO)

    return (
        <Paper sx={SX.box} variant={"outlined"}>
            <Box sx={SX.head}>
                <Box sx={SX.title} onClick={open ? handleCancel : handleToggleBody(BodyType.INFO)}>
                    <Box sx={SX.name}>{query.name}</Box>
                    <Box sx={{...SX.creation, color: info?.palette.text.disabled}}>({query.creation})</Box>
                </Box>
                <Box sx={SX.buttons}>
                    {open && (
                        <>
                            <Box sx={{...SX.type, color: info?.palette.text.disabled}} onClick={handleCancel}>
                                {BodyType[body]}
                            </Box>
                            <CancelIconButton onClick={handleCancel}/>
                        </>
                    )}
                    {!open && query.default !== query.custom && (
                        <RestoreIconButton onClick={handleToggleBody(BodyType.RESTORE)}/>
                    )}
                    {!open && <EditIconButton onClick={handleToggleBody(BodyType.EDIT)}/>}
                    <PlayIconButton color={"success"} onClick={handleToggleBody(BodyType.PLAY)}/>
                </Box>
            </Box>
            <Collapse in={open}>
                <Divider/>
                <Box sx={SX.body}>{renderBody()}</Box>
            </Collapse>
        </Paper>
    )

    function renderBody() {
        switch (body) {
            case BodyType.INFO:
                return <QueryItemInfo query={query.custom} description={query.description}/>
            case BodyType.EDIT:
                return <QueryItemEdit id={id} query={query.custom}/>
            case BodyType.RESTORE:
                return <QueryItemRestore id={id} def={query.default} custom={query.custom}/>
            case BodyType.PLAY:
                return <QueryItemPlay id={id}/>
            default:
                return <InfoAlert text={"Not implemented yet"}/>
        }
    }

    function handleCancel() {
        setBody(BodyType.INFO)
        setOpen(false)
    }

    function handleToggleBody(type: BodyType) {
        return () => {
            setBody(type)
            setOpen(true)
        }
    }
}
