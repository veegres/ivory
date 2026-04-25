import {Box} from "@mui/material"
import {useState} from "react"

import {Query as QueryRequest} from "../../../api/database/type"
import {Feature} from "../../../api/feature"
import {Connection, Query, QueryCreation} from "../../../api/query/type"
import {SxPropsMap} from "../../../app/type"
import {
    CancelIconButton,
    EditIconButton,
    LogIconButton,
    QueryViewIconButton,
    RestoreIconButton,
} from "../../view/button/IconButtons"
import {Access} from "../access/Access"
import {QueryBoxBody} from "./QueryBoxBody"
import {QueryButtonDelete} from "./QueryButtonDelete"
import {QueryButtonUpdate} from "./QueryButtonUpdate"
import {QueryInfoEdit} from "./QueryInfoEdit"
import {QueryInfoView} from "./QueryInfoView"
import {QueryLog} from "./QueryLog"
import {QueryRestore} from "./QueryRestore"
import {QueryTemplateWrapper} from "./QueryTemplateWrapper"

const SX: SxPropsMap = {
    name: {fontWeight: "bold"},
    creation: {fontSize: "12px", fontFamily: "monospace",  color: "text.disabled"},
    type: {padding: "0 3px", cursor: "pointer", color: "text.disabled"},
}

enum ViewToggleType {VIEW, EDIT, RESTORE, LOG}

type Props = {
    query: Query,
    connection: Connection,
}

export function QueryTemplateView(props: Props) {
    const {query, connection} = props
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
                <QueryInfoView query={query as any}/>
            </QueryBoxBody>
            <QueryBoxBody show={toggleView === ViewToggleType.EDIT}>
                <QueryInfoEdit query={queryUpdate as any} onChange={setUpdateQuery as any}/>
            </QueryBoxBody>
            <QueryBoxBody show={toggleView === ViewToggleType.RESTORE}>
                <QueryRestore query={query as any} onSuccess={handleToggleBody(ViewToggleType.VIEW)}/>
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
                <Access feature={Feature.ViewQueryLogList}>
                    <LogIconButton onClick={handleToggleBody(ViewToggleType.LOG)}/>
                </Access>
                <Access feature={Feature.ManageQueryCrudUpdate}>
                    {query.default !== query.custom && (
                        <RestoreIconButton onClick={handleToggleBody(ViewToggleType.RESTORE)}/>
                    )}
                    <EditIconButton onClick={handleToggleBody(ViewToggleType.EDIT)}/>
                </Access>
                <Access feature={Feature.ManageQueryCrudDelete}>
                    {query.creation === QueryCreation.Manual && (
                        <QueryButtonDelete id={query.id} type={query.type}/>
                    )}
                </Access>
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
