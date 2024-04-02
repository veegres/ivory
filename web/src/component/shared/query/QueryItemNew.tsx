import {InputBase} from "@mui/material";
import {useEffect, useState} from "react";
import {SxPropsMap} from "../../../type/general";
import {QueryConnection, QueryRequest, QueryType} from "../../../type/query";
import {QueryBody} from "./QueryBody";
import {QueryBodyInfoEdit} from "./QueryBodyInfoEdit";
import {QueryItemWrapper} from "./QueryItemWrapper";
import {CancelIconButton} from "../../view/button/IconButtons";
import {QueryButtonCreate} from "./QueryButtonCreate";

const SX: SxPropsMap = {
    input: {fontSize: "inherit", padding: "0"},
}

type Props = {
    type: QueryType,
    connection: QueryConnection,
}

export function QueryItemNew(props: Props) {
    const {type, connection} = props
    const [body, setBody] = useState(false)
    const [queryCreate, setQueryCreate] = useState<QueryRequest>({name: "", query: "", type})

    useEffect(handleEffectClose, [queryCreate.name, setBody]);

    return (
        <QueryItemWrapper
            connection={connection}
            params={queryCreate.params}
            varieties={queryCreate.varieties}
            renderTitle={renderTitle()}
            renderButtons={renderButtons()}
            showButtons={queryCreate.name !== ""}
            showInfo={true}
            query={queryCreate.query}
        >
            <QueryBody show={body}>
                <QueryBodyInfoEdit query={{...queryCreate, type}} onChange={setQueryCreate}/>
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

    function renderButtons() {
        return (
            <>
                {queryCreate.name !== "" && <CancelIconButton onClick={handleSuccess}/>}
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
