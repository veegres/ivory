import {QueryBody} from "./QueryBody";
import {QueryBodyRun} from "./QueryBodyRun";
import {QueryBoxCodeEditor} from "./QueryBoxCodeEditor";
import {useState} from "react";
import {Database, SxPropsMap} from "../../../type/common";
import {CancelIconButton, PlayIconButton} from "../../view/button/IconButtons";
import {Box} from "@mui/material";
import {QueryBoxWrapper} from "./QueryBoxWrapper";
import {useStore, useStoreAction} from "../../../provider/StoreProvider";
import {useDebounceFunction} from "../../../hook/Debounce";

const SX: SxPropsMap = {
    query: {position: "relative", display: "flex", flexDirection: "column", padding: "2px 2px 15px 2px", fontSize: "13px"},
    run: {position: "absolute", right: "10px", top: "4px", zIndex: 1, minWidth: 0},

}

type Props = {
    credentialId: string,
    db: Database,
}

export function QueryConsole(props: Props) {
    const {credentialId, db} = props
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
            <QueryBody show={run}>
                <QueryBodyRun request={{query, credentialId, db}}/>
            </QueryBody>
        </Box>
    )
}
