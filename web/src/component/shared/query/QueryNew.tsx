import {Alert, Box, InputBase, ToggleButton, ToggleButtonGroup} from "@mui/material";
import {useEffect, useState} from "react";
import {SxPropsMap} from "../../../type/common";
import {QueryEditor} from "./QueryEditor";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {queryApi} from "../../../app/api";
import {LoadingButton} from "@mui/lab";
import {QueryType, QueryVariety} from "../../../type/query";
import {InfoOutlined} from "@mui/icons-material";
import {InfoBox} from "../../view/box/InfoBox";
import {QueryDescription} from "./QueryDescription";
import {QueryItemPaper} from "./QueryItemPaper";
import {QueryItemHead} from "./QueryItemHead";
import {QueryItemBody} from "./QueryItemBody";
import {QueryVarietyOptions} from "../../../app/utils";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    toggle: {padding: "3px"},
    type: {color: "secondary.main"},
    input: {fontSize: "inherit", padding: "0"},
}

type Props = {
    type: QueryType,
}

export function QueryNew(props: Props) {
    const {type} = props
    const [body, setBody] = useState(false)
    const [alert, setAlert] = useState(false)
    const [name, setName] = useState("")
    const [description, setDesc] = useState("")
    const [query, setQuery] = useState("")
    const [varieties, setVarietis] = useState<QueryVariety[]>([])

    const createOptions = useMutationOptions([["query", "map", type]], handleSuccess)
    const create = useMutation({mutationFn: queryApi.create, ...createOptions})

    useEffect(handleEffectClose, [name, setBody]);

    return (
        <QueryItemPaper>
            <QueryItemHead renderTitle={renderTitle()} renderButtons={renderTitleButtons()}/>
            <QueryItemBody show={alert}>
                <Alert severity={"info"} onClose={() => setAlert(false)}>
                    Fields <i>name</i> and <i>query</i> are required for a new query. If you want to have termination
                    and query cancellation buttons in the table you need to call postgres <i>process_id</i> as
                    a <i>pid</i>. Example: <i>SELECT pid FROM table;</i>
                </Alert>
            </QueryItemBody>
            <QueryItemBody show={body}>
                <Box sx={SX.box}>
                    <QueryDescription>
                        <ToggleButtonGroup size={"small"} fullWidth value={varieties} onChange={(_e, v) => setVarietis(v)}>
                            {Object.entries(QueryVarietyOptions).map(([v, o]) => (
                                <ToggleButton key={v} sx={{padding: "0"}} value={v}>{o.label}</ToggleButton>
                            ))}
                        </ToggleButtonGroup>
                    </QueryDescription>
                    <QueryDescription editable>
                        <InputBase
                            sx={SX.input}
                            fullWidth
                            multiline
                            placeholder={"Description"}
                            value={description}
                            onChange={(e) => setDesc(e.target.value)}
                        />
                    </QueryDescription>
                    <QueryDescription editable>
                        <QueryEditor value={query} editable={true} autoFocus={false} onUpdate={setQuery}/>
                    </QueryDescription>
                </Box>
            </QueryItemBody>
        </QueryItemPaper>
    )

    function renderTitle() {
        return (
            <InputBase
                sx={SX.input}
                fullWidth
                required
                placeholder={"Type query name"}
                value={name}
                onChange={(e) => setName(e.target.value)}
            />
        )
    }

    function renderTitleButtons() {
        return (
            <>
                <InfoBox tooltip={"Type"} withPadding>
                    <Box sx={SX.type}>{QueryType[type]}</Box>
                </InfoBox>
                <LoadingButton
                    size={"small"}
                    loading={create.isPending}
                    variant={"outlined"}
                    onClick={handleAdd}
                    disabled={!name || !query}
                >
                    Add
                </LoadingButton>
                <ToggleButton
                    sx={SX.toggle} value={"info"} size={"small"} selected={alert}
                    onClick={() => setAlert(!alert)}>
                    <InfoOutlined/>
                </ToggleButton>
            </>
        )
    }

    function handleEffectClose() {
        if (name === "") setBody(false)
        else setBody(true)
    }

    function handleAdd() {
        create.mutate({name, type, description, query})
    }

    function handleSuccess() {
        setName("")
        setDesc("")
        setQuery("")
    }
}
