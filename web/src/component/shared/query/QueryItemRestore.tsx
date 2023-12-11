import {Box, Tooltip} from "@mui/material";
import {QueryDescription} from "./QueryDescription";
import {SxPropsMap} from "../../../type/common";
import {QueryEditor} from "./QueryEditor";
import {KeyboardDoubleArrowRight} from "@mui/icons-material";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {queryApi} from "../../../app/api";
import {LoadingButton} from "@mui/lab";

const SX: SxPropsMap = {
    box: {display: "grid", gridTemplateColumns: "minmax(0, 1fr) auto minmax(0, 1fr)", rowGap: 1, columnGap: 2},
    button: {display: "flex", justifyContent: "center", alignItems: "center"},
    bold: {fontWeight: "bold", padding: "0 10px"},
}

type Props = {
    id: string,
    def: string,
    custom: string,
}

export function QueryItemRestore(props: Props) {
    const {id, def, custom} = props

    const updateOptions = useMutationOptions([["query", "map"]])
    const update = useMutation({mutationFn: queryApi.update, ...updateOptions})

    return (
        <Box sx={SX.box}>
            <Box sx={SX.bold}>Default:</Box>
            <Box></Box>
            <Box sx={SX.bold}>Current:</Box>
            <QueryDescription><QueryEditor value={def} editable={false}/></QueryDescription>
            <Box sx={SX.button}>
                <Tooltip title={"Restore"} placement={"top"} disableInteractive>
                    <span>
                        <LoadingButton
                            disabled={def === custom}
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
            <QueryDescription><QueryEditor value={custom} editable={false}/></QueryDescription>
        </Box>
    )

    function handleUpdate() {
        update.mutate({id, query: {query: def}})
    }
}
