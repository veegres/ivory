import {Query, QueryConnection, QueryCreation, QueryRequest} from "../../../type/query";
import {SxPropsMap} from "../../../type/general";
import {QueryItemWrapper} from "./QueryItemWrapper";
import {QueryBody} from "./QueryBody";
import {QueryBodyInfoView} from "./QueryBodyInfoView";
import {QueryBodyInfoEdit} from "./QueryBodyInfoEdit";
import {QueryBodyRestore} from "./QueryBodyRestore";
import {useState} from "react";
import {
    CancelIconButton,
    EditIconButton,
    LogIconButton,
    QueryViewIconButton,
    RestoreIconButton
} from "../../view/button/IconButtons";
import {QueryButtonDelete} from "./QueryButtonDelete";
import {Box} from "@mui/material";
import {QueryBodyLog} from "./QueryBodyLog";
import {QueryButtonUpdate} from "./QueryButtonUpdate";

const SX: SxPropsMap = {
    name: {fontWeight: "bold"},
    creation: {fontSize: "12px", fontFamily: "monospace",  color: "text.disabled"},
    type: {padding: "0 3px", cursor: "pointer", color: "text.disabled"},
}

enum ViewToggleType {VIEW, EDIT, RESTORE, LOG}

type Props = {
    query: Query,
    connection: QueryConnection,
    manual: boolean,
}

export function QueryItemView(props: Props) {
    const {query, connection, manual} = props
    const [toggleView, setToggleView] = useState<ViewToggleType>()
    const initQueryUpdate: QueryRequest = {...query, query: query.custom}
    const [queryUpdate, setUpdateQuery] = useState(initQueryUpdate)

    return (
        <QueryItemWrapper
            connection={connection}
            queryUuid={query.id}
            varieties={query.varieties}
            params={query.params}
            renderButtons={renderToggleButtons()}
            showButtons={toggleView !== undefined}
            showInfo={toggleView === ViewToggleType.EDIT}
            renderTitle={renderTitle()}
            query={query.custom}
        >
            <QueryBody show={toggleView === ViewToggleType.VIEW}>
                <QueryBodyInfoView query={query}/>
            </QueryBody>
            <QueryBody show={toggleView === ViewToggleType.EDIT}>
                <QueryBodyInfoEdit query={queryUpdate} onChange={setUpdateQuery}/>
            </QueryBody>
            <QueryBody show={toggleView === ViewToggleType.RESTORE}>
                <QueryBodyRestore query={query} onSuccess={handleToggleBody(ViewToggleType.VIEW)}/>
            </QueryBody>
            <QueryBody show={toggleView === ViewToggleType.LOG}>
                <QueryBodyLog queryId={query.id}/>
            </QueryBody>
        </QueryItemWrapper>
    )

    function renderTitle() {
        return (
            <>
                <Box sx={SX.name}>{query.name}</Box>
                <Box sx={SX.creation}>({query.creation})</Box>
            </>
        )
    }

    function renderToggleButtons() {
        if (toggleView !== undefined) return renderCancelButton(ViewToggleType[toggleView])
        return (
            <>
                <QueryViewIconButton onClick={handleToggleBody(ViewToggleType.VIEW)}/>
                <LogIconButton onClick={handleToggleBody(ViewToggleType.LOG)}/>
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
                <Box sx={SX.type} onClick={handleClose}>{type}</Box>
                <CancelIconButton tooltip={`Close ${type}`} onClick={handleClose}/>
                {toggleView === ViewToggleType.EDIT && (
                    <QueryButtonUpdate id={query.id} query={queryUpdate} onSuccess={handleToggleBody(ViewToggleType.VIEW)}/>
                )}
            </>
        )
    }

    function handleToggleBody(type?: ViewToggleType) {
        return () => setToggleView(type)
    }

    function handleClose() {
        setUpdateQuery(initQueryUpdate)
        handleToggleBody(undefined)()
    }
}
