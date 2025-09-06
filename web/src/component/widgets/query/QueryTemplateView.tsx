import {Query, QueryConnection, QueryCreation, QueryRequest} from "../../../api/query/type";
import {QueryTemplateWrapper} from "./QueryTemplateWrapper";
import {QueryBoxBody} from "./QueryBoxBody";
import {QueryInfoView} from "./QueryInfoView";
import {QueryInfoEdit} from "./QueryInfoEdit";
import {QueryRestore} from "./QueryRestore";
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
import {QueryLog} from "./QueryLog";
import {QueryButtonUpdate} from "./QueryButtonUpdate";
import {SxPropsMap} from "../../../app/type";

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

export function QueryTemplateView(props: Props) {
    const {query, connection, manual} = props
    const [toggleView, setToggleView] = useState<ViewToggleType>()
    const initQueryUpdate: QueryRequest = {...query, query: query.custom}
    const [queryUpdate, setUpdateQuery] = useState(initQueryUpdate)

    return (
        <QueryTemplateWrapper
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
            <QueryBoxBody show={toggleView === ViewToggleType.VIEW}>
                <QueryInfoView query={query}/>
            </QueryBoxBody>
            <QueryBoxBody show={toggleView === ViewToggleType.EDIT}>
                <QueryInfoEdit query={queryUpdate} onChange={setUpdateQuery}/>
            </QueryBoxBody>
            <QueryBoxBody show={toggleView === ViewToggleType.RESTORE}>
                <QueryRestore query={query} onSuccess={handleToggleBody(ViewToggleType.VIEW)}/>
            </QueryBoxBody>
            <QueryBoxBody show={toggleView === ViewToggleType.LOG}>
                <QueryLog queryId={query.id}/>
            </QueryBoxBody>
        </QueryTemplateWrapper>
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
