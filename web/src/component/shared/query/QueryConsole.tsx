import {QueryBody} from "./QueryBody";
import {QueryBodyRun} from "./QueryBodyRun";
import {QueryBoxCodeEditor} from "./QueryBoxCodeEditor";
import {useState} from "react";
import {SxPropsMap} from "../../../type/general";
import {CancelIconButton, PlayIconButton} from "../../view/button/IconButtons";
import {Box, Paper} from "@mui/material";
import {QueryBoxWrapper} from "./QueryBoxWrapper";
import {useStore, useStoreAction} from "../../../provider/StoreProvider";
import {useDebounceFunction} from "../../../hook/Debounce";
import {QueryConnection} from "../../../type/query";

const SX: SxPropsMap = {
    query: {position: "relative", display: "flex", flexDirection: "column", padding: "2px 2px 15px 2px", fontSize: "13px"},
    run: {position: "absolute", right: "10px", top: "4.5px", zIndex: 1, minWidth: 0},

}

type Props = {
    connection: QueryConnection,
}

export function QueryConsole(props: Props) {
    const {connection} = props
    const {instance: {queryConsole}} = useStore()
    const {setConsoleQuery} = useStoreAction()
    const [run, setRun] = useState(false)
    const [query, setQuery] = useState(queryConsole)

    useDebounceFunction(query, setConsoleQuery)

    return (
        <Box>
            <Box sx={SX.query}>
                <QueryBoxWrapper editable={true}>
                    <QueryBoxCodeEditor
                        value={query}
                        editable={true}
                        autoFocus={true}
                        onUpdate={setQuery}
                    />
                </QueryBoxWrapper>
                <Box sx={SX.run}>
                    {!run ? (
                        <PlayIconButton color={"success"} onClick={() => setRun(true)}/>
                    ) : (
                        <CancelIconButton color={"success"} tooltip={"Close"} onClick={() => setRun(false)}/>
                    )}
                </Box>
            </Box>
            <Paper>
                <QueryBody show={run}>
                    <QueryBodyRun request={{connection, query}}/>
                </QueryBody>
            </Paper>
        </Box>
    )
}
