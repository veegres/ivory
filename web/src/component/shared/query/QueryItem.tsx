import {Box} from "@mui/material";
import {
    CancelIconButton,
    EditIconButton,
    PlayIconButton,
    QueryParamsIconButton,
    QueryViewIconButton,
    RestoreIconButton
} from "../../view/button/IconButtons";
import {useAppearance} from "../../../provider/AppearanceProvider";
import {useState} from "react";
import {QueryBodyInfoView} from "./QueryBodyInfoView";
import {QueryBodyInfoEdit} from "./QueryBodyInfoEdit";
import {QueryBodyRestore} from "./QueryBodyRestore";
import {QueryBody} from "./QueryBody";
import {QueryBodyRun} from "./QueryBodyRun";
import {Database, SxPropsMap} from "../../../type/common";
import {Query, QueryCreation, QueryRequest, QueryType} from "../../../type/query";
import {FixedInputs} from "../../view/input/FixedInputs";
import {QueryHead} from "./QueryHead";
import {QueryBoxPaper} from "./QueryBoxPaper";
import {QueryButtonDelete} from "./QueryButtonDelete";

const SX: SxPropsMap = {
    name: {fontWeight: "bold"},
    creation: {fontSize: "12px", fontFamily: "monospace"},
    type: {padding: "0 3px", cursor: "pointer", color: "text.disabled"},
}

enum ViewToggleType {INFO, EDIT, RESTORE, }
enum ViewCheckType {RUN, PARAMS}

type Props = {
    query: Query,
    credentialId?: string,
    db: Database,
    type: QueryType,
    editable: boolean,
}

export function QueryItem(props: Props) {
    const {query, credentialId, db, editable} = props
    const {info} = useAppearance()
    const [toggleView, setToggleView] = useState<ViewToggleType>()
    const [checkView, setCheckView] = useState<{[key in ViewCheckType]: boolean}>({[ViewCheckType.RUN]: false, [ViewCheckType.PARAMS]: false})
    const [params, setParams] = useState(query.params?.map(() => ''))
    const [queryUpdate, setUpdateQuery] = useState<QueryRequest>({...query, query: query.custom})

    const open = toggleView !== undefined
    return (
        <QueryBoxPaper>
            <QueryHead renderTitle={renderTitle()} renderButtons={renderTitleButtons()} onClick={handleCloseAll}/>
            <QueryBody show={toggleView === ViewToggleType.INFO}>
                <QueryBodyInfoView query={query}/>
            </QueryBody>
            <QueryBody show={toggleView === ViewToggleType.EDIT}>
                <QueryBodyInfoEdit query={queryUpdate} onChange={setUpdateQuery}/>
            </QueryBody>
            <QueryBody show={toggleView === ViewToggleType.RESTORE}>
                <QueryBodyRestore query={query}/>
            </QueryBody>
            <QueryBody show={checkView[ViewCheckType.PARAMS]}>
                <FixedInputs inputs={params ?? []} placeholders={query.params ?? []} onChange={v => setParams(v)}/>
            </QueryBody>
            <QueryBody show={checkView[ViewCheckType.RUN]}>
                <QueryBodyRun
                    request={{queryUuid: query.id, credentialId, db, queryParams: params}}
                    varieties={query.varieties}
                />
            </QueryBody>
        </QueryBoxPaper>
    )

    function renderTitle() {
        return (
            <>
                <Box sx={SX.name}>{query.name}</Box>
                <Box sx={{...SX.creation, color: info?.palette.text.disabled}}>({query.creation})</Box>
            </>
        )
    }

    function renderTitleButtons() {
        return (
            <>
                {renderToggleButtons()}
                {renderQueryParamsButton()}
                {renderRunButton()}
            </>
        )
    }

    function renderToggleButtons() {
        if (open) return renderCancelButton(ViewToggleType[toggleView])
        return (
            <>
                <QueryViewIconButton onClick={handleToggleBody(ViewToggleType.INFO)}/>
                {query.default !== query.custom && (
                    <RestoreIconButton onClick={handleToggleBody(ViewToggleType.RESTORE)}/>
                )}
                {editable && (
                    <EditIconButton onClick={handleToggleBody(ViewToggleType.EDIT)}/>
                )}
                {query.creation === QueryCreation.Manual && (
                    <QueryButtonDelete id={query.id} type={query.type}/>
                )}
            </>
        )
    }

    function renderQueryParamsButton() {
        const disabled = !query.params || query.params.length == 0
        return !checkView[ViewCheckType.PARAMS] ? (
            <QueryParamsIconButton color={"secondary"} disabled={disabled} onClick={handleCheckView(ViewCheckType.PARAMS)}/>
        ) : (
            <CancelIconButton color={"secondary"} tooltip={"Close Query Params"} onClick={handleCheckView(ViewCheckType.PARAMS)}/>
        )
    }

    function renderRunButton() {
        return !checkView[ViewCheckType.RUN] ? (
            <PlayIconButton color={"success"} onClick={handleCheckView(ViewCheckType.RUN)}/>
        ) : (
            <CancelIconButton color={"success"} tooltip={"Close"} onClick={handleCheckView(ViewCheckType.RUN)}/>
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

    function handleCloseAll() {
        setToggleView(undefined)
        setCheckView({
            [ViewCheckType.PARAMS]: false,
            [ViewCheckType.RUN]: false,
        })
    }

    function handleToggleBody(type?: ViewToggleType) {
        return () => setToggleView(type)
    }

    function handleCheckView(type: ViewCheckType) {
        return () => setCheckView({...checkView, [type]: !checkView[type]})
    }
}
