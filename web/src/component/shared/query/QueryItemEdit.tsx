import {QueryEditor} from "./QueryEditor";
import {SxPropsMap} from "../../../type/common";
import {Box} from "@mui/material";
import {CopyIconButton, SaveIconButton} from "../../view/button/IconButtons";
import {useMutation} from "@tanstack/react-query";
import {queryApi} from "../../../app/api";
import {useState} from "react";
import {useMutationOptions} from "../../../hook/QueryCustom";

const SX: SxPropsMap = {
    box: {display: "flex", gap: 1},
    editor: {flexGrow: 1, overflowX: "auto", border: "1px solid", borderColor: "divider"},
    buttons: {display: "flex", flexDirection: "column", gap: 1}
}

type Props = {
    id: string,
    query: string,
}

export function QueryItemEdit(props: Props) {
    const {id, query} = props
    const [newQuery, setNewQuery] = useState(query)

    const updateOptions = useMutationOptions([["query", "map"]])
    const update = useMutation(queryApi.update, updateOptions)

    return (
        <Box sx={SX.box}>
            <Box sx={SX.editor}>
                <QueryEditor value={query} editable={true} onUpdate={setNewQuery}/>
            </Box>
            <Box sx={SX.buttons}>
                <SaveIconButton loading={update.isLoading} placement={"left"} onClick={handleUpdate}/>
                <CopyIconButton placement={"left"} onClick={handleCopy}/>
            </Box>
        </Box>
    )

    function handleUpdate() {
        update.mutate({id, query: {query: newQuery}})
    }

    function handleCopy() {
        return navigator.clipboard.writeText(query)
    }
}
