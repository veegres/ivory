import {Box, IconButton, Paper, TextField} from "@mui/material";
import {useState} from "react";
import {QueryType, SxPropsMap} from "../../../app/types";
import {QueryEditor} from "./QueryEditor";
import {QueryItemBody} from "./QueryItemBody";
import {useTheme} from "../../../provider/ThemeProvider";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {queryApi} from "../../../app/api";
import {LoadingButton} from "@mui/lab";
import {OpenIcon} from "../../view/OpenIcon";

const SX: SxPropsMap = {
    box: {fontSize: "15px"},
    head: {display: "flex", alignItems: "center", justifyContent: "space-between", padding: "5px 15px", cursor: "pointer", fontWeight: "bold"},
    body: {display: "flex", flexDirection: "column", gap: 1},
    buttons: {display: "flex", alignItems: "center", gap: 1},
    desc: {flexGrow: 1},
    editor: {borderRadius: "5px"},
}

type Props = {
    type: QueryType,
}

export function QueryAdd(props: Props) {
    const {info} = useTheme()
    const [name, setName] = useState("")
    const [description, setDesc] = useState("")
    const [query, setQuery] = useState("")
    const [open, setOpen] = useState(false)

    const createOptions = useMutationOptions([["query", "map"]], handleSuccess)
    const create = useMutation(queryApi.create, createOptions)

    return (
        <Paper sx={SX.box} variant={"outlined"}>
            <Box sx={SX.head} onClick={() => setOpen(!open)}>
                <Box>New</Box>
                <IconButton size={"small"}>
                    <OpenIcon open={open} size={15}/>
                </IconButton>
            </Box>
            <QueryItemBody show={open}>
                <Box sx={SX.body}>
                    <Box sx={SX.buttons}>
                        <TextField
                            size={"small"}
                            value={name}
                            placeholder={"Name"}
                            onChange={(e) => setName(e.target.value)}
                        />
                        <TextField
                            sx={SX.desc}
                            value={description}
                            size={"small"}
                            placeholder={"Description"}
                            onChange={(e) => setDesc(e.target.value)}
                        />
                        <LoadingButton loading={create.isLoading} variant={"outlined"} onClick={handleAdd}>
                            Add
                        </LoadingButton>
                    </Box>
                    <Box sx={{...SX.editor, border: `2px solid ${info?.palette.divider}`}}>
                        <QueryEditor value={query} editable={true} onUpdate={setQuery}/>
                    </Box>
                </Box>
            </QueryItemBody>
        </Paper>
    )

    function handleAdd() {
        create.mutate({name, type: props.type, description, query})
    }

    function handleSuccess() {
        setOpen(false)
        setName("")
        setDesc("")
        setQuery("")
    }
}
