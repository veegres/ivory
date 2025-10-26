import {KeyboardDoubleArrowRight} from "@mui/icons-material"
import {Box, Button, Tooltip} from "@mui/material"

import {useRouterQueryUpdate} from "../../../api/query/hook"
import {Query} from "../../../api/query/type"
import {SxPropsMap} from "../../../app/type"
import {QueryBoxCodeEditor} from "./QueryBoxCodeEditor"
import {QueryBoxWrapper} from "./QueryBoxWrapper"

const SX: SxPropsMap = {
    box: {display: "grid", gridTemplateColumns: "minmax(0, 1fr) auto minmax(0, 1fr)", rowGap: 1, columnGap: 2},
    button: {display: "flex", justifyContent: "center", alignItems: "center"},
    bold: {fontWeight: "bold", padding: "0 10px"},
}

type Props = {
    query: Query,
    onSuccess: () => void,
}

export function QueryRestore(props: Props) {
    const {query, onSuccess} = props

    const update = useRouterQueryUpdate(query.type, onSuccess)

    return (
        <Box sx={SX.box}>
            <Box sx={SX.bold}>Default:</Box>
            <Box></Box>
            <Box sx={SX.bold}>Current:</Box>
            <QueryBoxWrapper><QueryBoxCodeEditor value={query.default} editable={false}/></QueryBoxWrapper>
            <Box sx={SX.button}>
                <Tooltip title={"Restore"} placement={"top"} disableInteractive>
                    <span>
                        <Button
                            disabled={query.default === query.custom}
                            loading={update.isPending}
                            size={"small"}
                            variant={"outlined"}
                            onClick={handleUpdate}
                        >
                            <KeyboardDoubleArrowRight/>
                        </Button>
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
