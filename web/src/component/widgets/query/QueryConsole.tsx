import {Box, Paper} from "@mui/material"
import {useState} from "react"

import {ConnectionRequest} from "../../../api/postgres"
import {SxPropsMap} from "../../../app/type"
import {useDebounceFunction} from "../../../hook/Debounce"
import {useStore, useStoreAction} from "../../../provider/StoreProvider"
import {CancelIconButton, PlayIconButton} from "../../view/button/IconButtons"
import {QueryBoxBody} from "./QueryBoxBody"
import {QueryBoxCodeEditor} from "./QueryBoxCodeEditor"
import {QueryBoxWrapper} from "./QueryBoxWrapper"
import {QueryRun} from "./QueryRun"

const SX: SxPropsMap = {
    query: {position: "relative", display: "flex", flexDirection: "column", padding: "2px 2px 15px 2px", fontSize: "13px"},
    run: {position: "absolute", right: "10px", top: "4.5px", zIndex: 1, minWidth: 0},

}

type Props = {
    connection: ConnectionRequest,
}

export function QueryConsole(props: Props) {
    const {connection} = props
    const {queryConsole} = useStore(s => s.instance)
    const {setConsoleQuery} = useStoreAction
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
                <QueryBoxBody show={run}>
                    <QueryRun connection={connection} query={query}/>
                </QueryBoxBody>
            </Paper>
        </Box>
    )
}
