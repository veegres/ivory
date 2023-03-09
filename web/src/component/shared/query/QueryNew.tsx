import {Box, Paper, TextField} from "@mui/material";
import {useState} from "react";
import {QueryType, SxPropsMap} from "../../../app/types";
import {QueryEditor} from "./QueryEditor";
import {useTheme} from "../../../provider/ThemeProvider";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {queryApi} from "../../../app/api";
import {LoadingButton} from "@mui/lab";

const SX: SxPropsMap = {
    body: {display: "flex", flexDirection: "column", gap: 1, marginBottom: "8px", padding: "10px 15px"},
    buttons: {display: "flex", alignItems: "center", gap: 1},
    desc: {flexGrow: 1},
    editor: {borderRadius: "5px"},
    name: {color: "text.disabled"},
}

type Props = {
    type: QueryType,
}

export function QueryNew(props: Props) {
    const {type} = props
    const {info} = useTheme()
    const [name, setName] = useState("")
    const [description, setDesc] = useState("")
    const [query, setQuery] = useState("")

    const createOptions = useMutationOptions([["query", "map", type]], handleSuccess)
    const create = useMutation(queryApi.create, createOptions)

    return (
        <Paper sx={SX.body} variant={"outlined"}>
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
        </Paper>
    )

    function handleAdd() {
        create.mutate({name, type, description, query})
    }

    function handleSuccess() {
        setName("")
        setDesc("")
        setQuery("")
    }
}
