import {Query, SxPropsMap} from "../../../app/types";
import {Box, Paper} from "@mui/material";
import {CancelIconButton, EditIconButton, PlayIconButton, RestoreIconButton} from "../../view/IconButtons";
import {useTheme} from "../../../provider/ThemeProvider";
import {useState} from "react";
import {QueryItemInfo} from "./QueryItemInfo";
import {QueryItemEdit} from "./QueryItemEdit";
import {QueryItemRun} from "./QueryItemRun";
import {QueryItemRestore} from "./QueryItemRestore";
import {QueryItemBody} from "./QueryItemBody";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", fontSize: "15px"},
    head: {display: "flex", padding: "5px 15px"},
    title: {flexGrow: 1, display: "flex", alignItems: "center", cursor: "pointer", gap: 1},
    name: {fontWeight: "bold"},
    creation: {fontSize: "12px", fontFamily: "monospace"},
    buttons: {display: "flex", alignItems: "center"},
    type: {padding: "0 8px", cursor: "pointer"},
}

enum BodyType {INFO, EDIT, RESTORE, RUN}

type Props = {
    id: string,
    query: Query,
}

export function QueryItem(props: Props) {
    const {id, query} = props
    const {info} = useTheme()
    const [body, setBody] = useState<BodyType>()
    const open = body !== undefined

    return (
        <Paper sx={SX.box} variant={"outlined"}>
            <Box sx={SX.head}>
                <Box sx={SX.title} onClick={open ? handleCancel : handleToggleBody(BodyType.INFO)}>
                    <Box sx={SX.name}>{query.name}</Box>
                    <Box sx={{...SX.creation, color: info?.palette.text.disabled}}>({query.creation})</Box>
                </Box>
                <Box sx={SX.buttons}>
                    {renderCancelButton()}
                    {!open && query.default !== query.custom && (
                        <RestoreIconButton onClick={handleToggleBody(BodyType.RESTORE)}/>
                    )}
                    {!open && <EditIconButton onClick={handleToggleBody(BodyType.EDIT)}/>}
                    <PlayIconButton color={"success"} onClick={handleToggleBody(BodyType.RUN)}/>
                </Box>
            </Box>
            <QueryItemBody show={body === BodyType.INFO}>
                <QueryItemInfo query={query.custom} description={query.description}/>
            </QueryItemBody>
            <QueryItemBody show={body === BodyType.EDIT}>
                <QueryItemEdit id={id} query={query.custom}/>
            </QueryItemBody>
            <QueryItemBody show={body === BodyType.RESTORE}>
                <QueryItemRestore id={id} def={query.default} custom={query.custom}/>
            </QueryItemBody>
            <QueryItemBody show={body === BodyType.RUN}>
                <QueryItemRun id={id}/>
            </QueryItemBody>
        </Paper>
    )

    function renderCancelButton() {
        if (!open) return
        return (
            <>
                <Box sx={{...SX.type, color: info?.palette.text.disabled}} onClick={handleCancel}>
                    {BodyType[body]}
                </Box>
                <CancelIconButton onClick={handleCancel}/>
            </>
        )
    }

    function handleCancel() {
        setBody(undefined)
    }

    function handleToggleBody(type: BodyType) {
        return () => setBody(type)
    }
}
