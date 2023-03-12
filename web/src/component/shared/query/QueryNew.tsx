import {Alert, Box, Paper, TextField, ToggleButton} from "@mui/material";
import {useState} from "react";
import {SxPropsMap} from "../../../type/common";
import {QueryEditor} from "./QueryEditor";
import {useTheme} from "../../../provider/ThemeProvider";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {queryApi} from "../../../app/api";
import {LoadingButton} from "@mui/lab";
import {QueryType} from "../../../type/query";
import {InfoOutlined} from "@mui/icons-material";
import {InfoBox} from "../../view/InfoBox";

const SX: SxPropsMap = {
    body: {display: "flex", flexDirection: "column", gap: 1, marginBottom: "8px", padding: "10px 15px"},
    buttons: {display: "flex", alignItems: "center", gap: 1},
    title: {display: "flex", alignItems: "center", justifyContent: "space-between"},
    info: {display: "flex", gap: 1},
    toggle: {padding: "3px"},
    desc: {flexGrow: 1},
    type: {color: "secondary.main"},
    editor: {borderRadius: "5px"},
    name: {color: "text.disabled"},
}

type Props = {
    type: QueryType,
}

export function QueryNew(props: Props) {
    const {type} = props
    const {info} = useTheme()
    const [alert, setAlert] = useState(false)
    const [name, setName] = useState("")
    const [description, setDesc] = useState("")
    const [query, setQuery] = useState("")

    const createOptions = useMutationOptions([["query", "map", type]], handleSuccess)
    const create = useMutation(queryApi.create, createOptions)

    return (
        <Paper sx={SX.body} variant={"outlined"}>
            <Box sx={SX.title}>
                <Box><b>NEW QUERY</b></Box>
                <Box sx={SX.info}>
                    <InfoBox tooltip={""}><Box sx={SX.type}>{QueryType[type]}</Box></InfoBox>
                    <ToggleButton
                        sx={SX.toggle}
                        value={"info"}
                        size={"small"}
                        selected={alert}
                        onClick={() => setAlert(!alert)}
                    >
                        <InfoOutlined/>
                    </ToggleButton>
                </Box>
            </Box>
            {alert && (
                <Alert variant={"outlined"} severity={"info"} icon={false}>
                    Fields <i>name</i> and <i>query</i> are required for a new query. If you want to have termination
                    and query cancellation buttons in the table you need to call postgres <i>process_id</i> as
                    a <i>pid</i>. Example: <i>SELECT pid FROM table;</i>
                </Alert>
            )}
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
                <LoadingButton
                    loading={create.isLoading}
                    variant={"outlined"}
                    onClick={handleAdd}
                    disabled={!name || !query}
                >
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
