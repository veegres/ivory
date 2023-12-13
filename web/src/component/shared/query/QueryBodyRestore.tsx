import {Box, Tooltip} from "@mui/material";
import {QueryBoxWrapper} from "./QueryBoxWrapper";
import {SxPropsMap} from "../../../type/common";
import {QueryBoxCodeEditor} from "./QueryBoxCodeEditor";
import {KeyboardDoubleArrowRight} from "@mui/icons-material";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {queryApi} from "../../../app/api";
import {LoadingButton} from "@mui/lab";
import {Query} from "../../../type/query";

const SX: SxPropsMap = {
    box: {display: "grid", gridTemplateColumns: "minmax(0, 1fr) auto minmax(0, 1fr)", rowGap: 1, columnGap: 2},
    button: {display: "flex", justifyContent: "center", alignItems: "center"},
    bold: {fontWeight: "bold", padding: "0 10px"},
}

type Props = {
    query: Query,
    onSuccess: () => void,
}

export function QueryBodyRestore(props: Props) {
    const {query, onSuccess} = props

    const updateOptions = useMutationOptions([["query", "map", query.type]], onSuccess)
    const update = useMutation({mutationFn: queryApi.update, ...updateOptions})

    return (
        <Box sx={SX.box}>
            <Box sx={SX.bold}>Default:</Box>
            <Box></Box>
            <Box sx={SX.bold}>Current:</Box>
            <QueryBoxWrapper><QueryBoxCodeEditor value={query.default} editable={false}/></QueryBoxWrapper>
            <Box sx={SX.button}>
                <Tooltip title={"Restore"} placement={"top"} disableInteractive>
                    <span>
                        <LoadingButton
                            disabled={query.default === query.custom}
                            loading={update.isPending}
                            size={"small"}
                            variant={"outlined"}
                            onClick={handleUpdate}
                        >
                            <KeyboardDoubleArrowRight/>
                        </LoadingButton>
                    </span>
                </Tooltip>
            </Box>
            <QueryBoxWrapper><QueryBoxCodeEditor value={query.custom} editable={false}/></QueryBoxWrapper>
        </Box>
    )

    function handleUpdate() {
        update.mutate({id: query.id, query: {...query, query: query.default}})
    }
}
