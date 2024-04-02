import {CancelIconButton, InfoIconButton, PlayIconButton, QueryParamsIconButton,} from "../../view/button/IconButtons";
import {ReactNode, useMemo, useState} from "react";
import {QueryBody} from "./QueryBody";
import {QueryBodyRun} from "./QueryBodyRun";
import {QueryConnection, QueryVariety} from "../../../type/query";
import {QueryHead} from "./QueryHead";
import {QueryBoxPaper} from "./QueryBoxPaper";
import {FixedInputs} from "../../view/input/FixedInputs";
import {Alert} from "@mui/material";

enum ViewCheckType {RUN, PARAMS, EDIT_INFO}

type ViewCheckMap = {[key in ViewCheckType]: boolean}
const initialViewCheckMap: ViewCheckMap = {
    [ViewCheckType.RUN]: false,
    [ViewCheckType.PARAMS]: false,
    [ViewCheckType.EDIT_INFO]: false,
}

type Props = {
    connection: QueryConnection,
    queryUuid?: string,
    varieties?: QueryVariety[],
    params?: string[],
    children: ReactNode,
    renderButtons: ReactNode,
    showButtons: boolean,
    showInfo: boolean,
    renderTitle: ReactNode,
    query: string,
}

export function QueryItemWrapper(props: Props) {
    const {queryUuid, params, varieties, query} = props
    const {connection, showInfo, showButtons} = props
    const {children, renderButtons, renderTitle} = props

    const [checkView, setCheckView] = useState(initialViewCheckMap)
    const [paramsValues, setParamsValues] = useState<string[]>(params?.map(() => "") ?? [])
    const queryRunRequest = useMemo( () => ({queryUuid, query, connection, queryParams: paramsValues}), [connection, paramsValues, query, queryUuid])

    return (
        <QueryBoxPaper>
            <QueryHead
                renderTitle={renderTitle}
                showAllButtons={showButtons}
                renderHiddenButtons={renderHiddeTitleButtons()}
                renderButtons={renderTitleButtons()}
            />
            <QueryBody show={showInfo && checkView[ViewCheckType.EDIT_INFO]}>
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
            <QueryBody show={checkView[ViewCheckType.PARAMS]} unmountOnExit={false}>
                <FixedInputs placeholders={params ?? []} values={paramsValues} onChange={setParamsValues}/>
            </QueryBody>
            <QueryBody show={checkView[ViewCheckType.RUN]}>
                <QueryBodyRun request={queryRunRequest} varieties={varieties}/>
            </QueryBody>
        </QueryBoxPaper>
    )

    function renderTitleButtons() {
        return (
            <>
                {renderQueryParamsButton()}
                {renderRunButton()}
            </>
        )
    }

    function renderHiddeTitleButtons() {
        return (
            <>
                {renderButtons}
                {renderEditButtons()}
            </>
        )
    }

    function renderEditButtons() {
        if (!showInfo) return

        return !checkView[ViewCheckType.EDIT_INFO] ? (
            <InfoIconButton color={"primary"} onClick={handleCheckView(ViewCheckType.EDIT_INFO)}/>
        ) : (
            <CancelIconButton color={"primary"} onClick={handleCheckView(ViewCheckType.EDIT_INFO)}/>
        )
    }

    function renderQueryParamsButton() {
        return !checkView[ViewCheckType.PARAMS] ? (
            <QueryParamsIconButton color={"secondary"} disabled={!params || params.length == 0} onClick={handleCheckView(ViewCheckType.PARAMS)}/>
        ) : (
            <CancelIconButton color={"secondary"} tooltip={"Close Query Params"} onClick={handleCheckView(ViewCheckType.PARAMS)}/>
        )
    }

    function renderRunButton() {
        return !checkView[ViewCheckType.RUN] ? (
            <PlayIconButton color={"success"} disabled={!query && !queryUuid} onClick={handleCheckView(ViewCheckType.RUN)}/>
        ) : (
            <CancelIconButton color={"success"} tooltip={"Close"} onClick={handleCheckView(ViewCheckType.RUN)}/>
        )
    }

    function handleCheckView(type: ViewCheckType) {
        return () => setCheckView({...checkView, [type]: !checkView[type]})
    }
}
