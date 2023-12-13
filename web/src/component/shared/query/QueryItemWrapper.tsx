import {CancelIconButton, InfoIconButton, PlayIconButton, QueryParamsIconButton,} from "../../view/button/IconButtons";
import {ReactNode, useState} from "react";
import {QueryBody} from "./QueryBody";
import {QueryBodyRun} from "./QueryBodyRun";
import {Database} from "../../../type/common";
import {QueryRequest, QueryType, QueryVariety} from "../../../type/query";
import {QueryHead} from "./QueryHead";
import {QueryBoxPaper} from "./QueryBoxPaper";
import {FixedInputs} from "../../view/input/FixedInputs";
import {Alert} from "@mui/material";
import {QueryButtonCreate} from "./QueryButtonCreate";
import {QueryButtonUpdate} from "./QueryButtonUpdate";

enum ViewCheckType {RUN, PARAMS, EDIT_INFO}

type ViewCheckMap = {[key in ViewCheckType]: boolean}
const initialViewCheckMap: ViewCheckMap = {
    [ViewCheckType.RUN]: false,
    [ViewCheckType.PARAMS]: false,
    [ViewCheckType.EDIT_INFO]: false,
}

type Props = {
    queryUuid?: string,
    varieties?: QueryVariety[],
    params?: string[],
    credentialId: string,
    db: Database,
    type: QueryType,
    children: ReactNode,
    renderButtons?: ReactNode,
    renderTitle: ReactNode,

    edit?: "create" | "update",
    editRequest: QueryRequest,
    onEditSuccess?: () => void,
}

export function QueryItemWrapper(props: Props) {
    const {queryUuid, params, varieties} = props
    const {edit, editRequest, onEditSuccess} = props
    const {credentialId, db} = props
    const {children, renderButtons, renderTitle} = props
    const [checkView, setCheckView] = useState(initialViewCheckMap)
    const [paramsValues, setParamsValues] = useState<string[]>()

    const isNoParams = !params || params.length == 0
    return (
        <QueryBoxPaper>
            <QueryHead renderTitle={renderTitle} renderButtons={renderTitleButtons()}/>
            <QueryBody show={checkView[ViewCheckType.EDIT_INFO]}>
                <Alert severity={"info"}>
                    The fields <i>name</i> and <i>query</i> are required. To enable termination and query
                    cancellation buttons in the table you need to call postgres <i>process_id</i> as <i>pid</i> (buttons
                    will always appear if you have this field). Example: <i>SELECT pid FROM table;</i> It is
                    also possible to specify prepare statement fields to make query template more flexible in usage.
                    You need to provide params names in a block with params and use <i>$1</i>, <i>$2</i>, etc to
                    define param fields in the query.
                </Alert>
            </QueryBody>
            {children}
            <QueryBody show={!isNoParams && checkView[ViewCheckType.PARAMS]}>
                <FixedInputs placeholders={params ?? []} onChange={setParamsValues}/>
            </QueryBody>
            <QueryBody show={checkView[ViewCheckType.RUN]}>
                <QueryBodyRun
                    request={{queryUuid, query: editRequest.query, credentialId, db, queryParams: paramsValues}}
                    varieties={varieties}
                />
            </QueryBody>
        </QueryBoxPaper>
    )

    function renderTitleButtons() {
        return (
            <>
                {renderButtons}
                {renderEditButtons()}
                {renderQueryParamsButton()}
                {renderRunButton()}
            </>
        )
    }

    function renderEditButtons() {
        if (edit === undefined) return

        return (
            <>
                {edit === "create" && <QueryButtonCreate query={editRequest} onSuccess={onEditSuccess}/>}
                {edit === "update" && queryUuid && <QueryButtonUpdate id={queryUuid} query={editRequest} onSuccess={onEditSuccess}/>}
                {!checkView[ViewCheckType.EDIT_INFO] ? (
                    <InfoIconButton color={"primary"} onClick={handleCheckView(ViewCheckType.EDIT_INFO)}/>
                ) : (
                    <CancelIconButton color={"primary"} onClick={handleCheckView(ViewCheckType.EDIT_INFO)}/>
                )}
            </>
        )
    }

    function renderQueryParamsButton() {
        return !checkView[ViewCheckType.PARAMS] ? (
            <QueryParamsIconButton color={"secondary"} disabled={isNoParams} onClick={handleCheckView(ViewCheckType.PARAMS)}/>
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

    function handleCheckView(type: ViewCheckType) {
        return () => setCheckView({...checkView, [type]: !checkView[type]})
    }
}
