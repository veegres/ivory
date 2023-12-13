import {Query, QueryCreation, QueryRequest, QueryType} from "../../../type/query";
import {Database, SxPropsMap} from "../../../type/common";
import {QueryItemWrapper} from "./QueryItemWrapper";
import {QueryBody} from "./QueryBody";
import {QueryBodyInfoView} from "./QueryBodyInfoView";
import {QueryBodyInfoEdit} from "./QueryBodyInfoEdit";
import {QueryBodyRestore} from "./QueryBodyRestore";
import {useMemo, useState} from "react";
import {CancelIconButton, EditIconButton, QueryViewIconButton, RestoreIconButton} from "../../view/button/IconButtons";
import {QueryButtonDelete} from "./QueryButtonDelete";
import {Box} from "@mui/material";

const SX: SxPropsMap = {
    name: {fontWeight: "bold"},
    creation: {fontSize: "12px", fontFamily: "monospace"},
    type: {padding: "0 3px", cursor: "pointer", color: "text.disabled"},
}

enum ViewToggleType {INFO, EDIT, RESTORE}

type Props = {
    query: Query,
    credentialId: string,
    db: Database,
    type: QueryType,
    manual: boolean,
}

export function QueryItemView(props: Props) {
    const {query, credentialId, db, type, manual} = props
    const [toggleView, setToggleView] = useState<ViewToggleType>()
    const initQueryUpdate: QueryRequest = useMemo(() => ({...query, query: query.custom}), [query])
    const [queryUpdate, setUpdateQuery] = useState(initQueryUpdate)

    return (
        <QueryItemWrapper
            queryUuid={query.id}
            varieties={query.varieties}
            params={query.params}
            db={db} type={type}
            credentialId={credentialId}
            renderButtons={renderToggleButtons()}
            renderTitle={renderTitle()}
        >
            <QueryBody show={toggleView === ViewToggleType.INFO}>
                <QueryBodyInfoView query={query}/>
            </QueryBody>
            <QueryBody show={toggleView === ViewToggleType.EDIT}>
                <QueryBodyInfoEdit query={queryUpdate} onChange={setUpdateQuery}/>
            </QueryBody>
            <QueryBody show={toggleView === ViewToggleType.RESTORE}>
                <QueryBodyRestore query={query}/>
            </QueryBody>
        </QueryItemWrapper>
    )

    function renderTitle() {
        return (
            <>
                <Box sx={SX.name}>{query.name}</Box>
                <Box sx={{...SX.creation, color: "text.disabled"}}>({query.creation})</Box>
            </>
        )
    }

    function renderToggleButtons() {
        if (toggleView !== undefined) return renderCancelButton(ViewToggleType[toggleView])
        return (
            <>
                <QueryViewIconButton onClick={handleToggleBody(ViewToggleType.INFO)}/>
                {query.default !== query.custom && (
                    <RestoreIconButton onClick={handleToggleBody(ViewToggleType.RESTORE)}/>
                )}
                {manual && (
                    <EditIconButton onClick={handleToggleBody(ViewToggleType.EDIT)}/>
                )}
                {query.creation === QueryCreation.Manual && (
                    <QueryButtonDelete id={query.id} type={query.type}/>
                )}
            </>
        )
    }

    function renderCancelButton(type: string) {
        return (
            <>
                <Box sx={SX.type} onClick={handleToggleBody(undefined)}>{type}</Box>
                <CancelIconButton tooltip={"Close Query Body"} onClick={handleToggleBody(undefined)}/>
            </>
        )
    }

    function handleToggleBody(type?: ViewToggleType) {
        return () => setToggleView(type)
    }
}
