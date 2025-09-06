import {InputBase} from "@mui/material";
import {useEffect, useState} from "react";
import {SxPropsMap} from "../../../api/management/type";
import {QueryConnection, QueryRequest, QueryType} from "../../../api/query/type";
import {QueryBoxBody} from "./QueryBoxBody";
import {QueryInfoEdit} from "./QueryInfoEdit";
import {QueryTemplateWrapper} from "./QueryTemplateWrapper";
import {CancelIconButton} from "../../view/button/IconButtons";
import {QueryButtonCreate} from "./QueryButtonCreate";

const SX: SxPropsMap = {
    input: {fontSize: "inherit", padding: "0"},
}

type Props = {
    type: QueryType,
    connection: QueryConnection,
}

export function QueryTemplateNew(props: Props) {
    const {type, connection} = props
    const [body, setBody] = useState(false)
    const [queryCreate, setQueryCreate] = useState<QueryRequest>({name: "", query: "", type})

    useEffect(handleEffectClose, [queryCreate.name, setBody]);

    return (
        <QueryTemplateWrapper
            connection={connection}
            params={queryCreate.params}
            varieties={queryCreate.varieties}
            renderTitle={renderTitle()}
            renderButtons={renderButtons()}
            showButtons={queryCreate.name !== ""}
            showInfo={true}
            query={queryCreate.query}
        >
            <QueryBoxBody show={body}>
                <QueryInfoEdit query={{...queryCreate, type}} onChange={setQueryCreate}/>
            </QueryBoxBody>
        </QueryTemplateWrapper>
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
