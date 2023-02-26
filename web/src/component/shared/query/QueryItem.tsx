import {Query, SxPropsMap} from "../../../app/types";
import {Box, Collapse, Divider, Paper} from "@mui/material";
import {EditIconButton, PlayIconButton, UndoIconButton} from "../../view/IconButtons";
import {useTheme} from "../../../provider/ThemeProvider";
import {useState} from "react";
import {QueryItemInfo} from "./QueryItemInfo";
import {QueryItemEdit} from "./QueryItemEdit";
import {QueryItemPlay} from "./QueryItemPlay";
import {InfoAlert} from "../../view/InfoAlert";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", fontSize: "15px"},
    head: {display: "flex", padding: "5px 15px"},
    body: {padding: "8px 15px", fontSize: "13px"},
    title: {flexGrow: 1, display: "flex", alignItems: "center", cursor: "pointer", gap: 1},
    name: {fontWeight: "bold"},
    creation: {fontSize: "12px", fontFamily: "monospace"},
    buttons: {display: "flex", alignItems: "center", justifyContent: "space-between", gap: 1}
}

enum BodyType {INFO, EDIT, PLAY}

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
                <Box sx={SX.title} onClick={handleToggleBody(BodyType.INFO)}>
                    <Box sx={SX.name}>{query.name}</Box>
                    <Box sx={{...SX.creation, color: info?.palette.text.disabled}}>({query.creation})</Box>
                </Box>
                <Box sx={SX.buttons}>
                    <UndoIconButton size={28} onClick={() => {
                    }}/>
                    <EditIconButton size={28} onClick={handleToggleBody(BodyType.EDIT)}/>
                    <PlayIconButton size={28} color={"success"} onClick={handleToggleBody(BodyType.PLAY)}/>
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
                return <QueryItemInfo query={query}/>
            case BodyType.EDIT:
                return <QueryItemEdit query={query}/>
            case BodyType.PLAY:
                return <QueryItemPlay id={id}/>
            default:
                return <InfoAlert text={"Not implemented yet"}/>
        }
    }

    function handleToggleBody(type: BodyType) {
        return () => {
            if (body === type && open) setOpen(false)
            else setOpen(true)
            setBody(type)
        }
    }
}
