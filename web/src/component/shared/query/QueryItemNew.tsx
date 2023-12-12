import {Alert, InputBase} from "@mui/material";
import {useEffect, useState} from "react";
import {SxPropsMap} from "../../../type/common";
import {QueryRequest, QueryType} from "../../../type/query";
import {QueryBoxPaper} from "./QueryBoxPaper";
import {QueryHead} from "./QueryHead";
import {QueryBody} from "./QueryBody";
import {CancelIconButton, InfoIconButton} from "../../view/button/IconButtons";
import {QueryButtonCreate} from "./QueryButtonCreate";
import {QueryBodyInfoEdit} from "./QueryBodyInfoEdit";

const SX: SxPropsMap = {
    input: {fontSize: "inherit", padding: "0"},
}

type Props = {
    type: QueryType,
}

export function QueryItemNew(props: Props) {
    const {type} = props
    const [body, setBody] = useState(false)
    const [info, setInfo] = useState(false)
    const [query, setQuery] = useState<QueryRequest>({name: "", query: "", type})

    useEffect(handleEffectClose, [query.name, setBody]);

    return (
        <QueryBoxPaper>
            <QueryHead renderTitle={renderTitle()} renderButtons={renderTitleButtons()}/>
            <QueryBody show={info}>
                <Alert severity={"info"}>
                    Fields <i>name</i> and <i>query</i> are required for a new query. If you want to have termination
                    and query cancellation buttons in the table you need to call postgres <i>process_id</i> as
                    a <i>pid</i>. Example: <i>SELECT pid FROM table;</i>
                </Alert>
            </QueryBody>
            <QueryBody show={body}>
                <QueryBodyInfoEdit query={query} onChange={setQuery}/>
            </QueryBody>
        </QueryBoxPaper>
    )



    function renderTitle() {
        return (
            <InputBase
                sx={SX.input}
                fullWidth
                required
                placeholder={"Type query name"}
                value={query.name}
                onChange={(e) => setQuery({...query, name: e.target.value})}
            />
        )
    }

    function renderTitleButtons() {
        return (
            <>
                {!info ? (
                    <InfoIconButton color={"primary"} onClick={() => setInfo(true)}/>
                ) : (
                    <CancelIconButton color={"primary"} onClick={() => setInfo(false)}/>
                )}
                <QueryButtonCreate query={query} onSuccess={handleSuccess}/>
            </>
        )
    }

    function handleEffectClose() {
        if (query.name === "") setBody(false)
        else setBody(true)
    }

    function handleSuccess() {
        setQuery({name: "", query: "", type})
    }
}
