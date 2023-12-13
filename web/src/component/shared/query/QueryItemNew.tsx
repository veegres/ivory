import {Alert, InputBase} from "@mui/material";
import {useEffect, useState} from "react";
import {Database, SxPropsMap} from "../../../type/common";
import {QueryRequest, QueryType} from "../../../type/query";
import {QueryBody} from "./QueryBody";
import {CancelIconButton, InfoIconButton} from "../../view/button/IconButtons";
import {QueryButtonCreate} from "./QueryButtonCreate";
import {QueryBodyInfoEdit} from "./QueryBodyInfoEdit";
import {QueryItemWrapper} from "./QueryItemWrapper";

const SX: SxPropsMap = {
    input: {fontSize: "inherit", padding: "0"},
}

type Props = {
    type: QueryType,
    credentialId: string,
    db: Database,
}

export function QueryItemNew(props: Props) {
    const {type, credentialId, db} = props
    const [body, setBody] = useState(false)
    const [info, setInfo] = useState(false)
    const [queryCreate, setQueryCreate] = useState<QueryRequest>({name: "", query: "", type})

    useEffect(handleEffectClose, [queryCreate.name, setBody]);

    return (
        <QueryItemWrapper
            queryUuid={"SHOULD BE ADDED IN BACKEND"}
            params={queryCreate.params}
            varieties={queryCreate.varieties}
            credentialId={credentialId}
            db={db}
            type={type}
            renderTitle={renderTitle()}
            renderButtons={renderTitleButtons()}
        >
            <QueryBody show={info}>
                <Alert severity={"info"}>
                    Fields <i>name</i> and <i>query</i> are required for a new query. If you want to have termination
                    and query cancellation buttons in the table you need to call postgres <i>process_id</i> as
                    a <i>pid</i>. Example: <i>SELECT pid FROM table;</i>
                </Alert>
            </QueryBody>
            <QueryBody show={body}>
                <QueryBodyInfoEdit query={queryCreate} onChange={setQueryCreate}/>
            </QueryBody>
        </QueryItemWrapper>
    )



    function renderTitle() {
        return (
            <InputBase
                sx={SX.input}
                fullWidth
                required
                placeholder={"Type query name"}
                value={queryCreate.name}
                onChange={(e) => setQueryCreate({...queryCreate, name: e.target.value})}
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
                <QueryButtonCreate query={queryCreate} onSuccess={handleSuccess}/>
            </>
        )
    }

    function handleEffectClose() {
        if (queryCreate.name === "") setBody(false)
        else setBody(true)
    }

    function handleSuccess() {
        setQueryCreate({name: "", query: "", type})
    }
}
