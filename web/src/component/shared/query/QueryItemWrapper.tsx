import {CancelIconButton, PlayIconButton, QueryParamsIconButton,} from "../../view/button/IconButtons";
import {ReactNode, useState} from "react";
import {QueryBody} from "./QueryBody";
import {QueryBodyRun} from "./QueryBodyRun";
import {Database} from "../../../type/common";
import {QueryType, QueryVariety} from "../../../type/query";
import {QueryHead} from "./QueryHead";
import {QueryBoxPaper} from "./QueryBoxPaper";
import {FixedInputs} from "../../view/input/FixedInputs";

enum ViewCheckType {RUN, PARAMS}

type ViewCheckMap = {[key in ViewCheckType]: boolean}
type Props = {
    queryUuid: string,
    varieties?: QueryVariety[],
    params?: string[],
    credentialId: string,
    db: Database,
    type: QueryType,
    children: ReactNode,
    renderButtons: ReactNode,
    renderTitle: ReactNode,
}

export function QueryItemWrapper(props: Props) {
    const {queryUuid, params, varieties} = props
    const {credentialId, db} = props
    const {children, renderButtons, renderTitle} = props
    const [checkView, setCheckView] = useState<ViewCheckMap>({[ViewCheckType.RUN]: false, [ViewCheckType.PARAMS]: false})
    const [paramsValues, setParamsValues] = useState<string[]>()

    const isNoParams = !params || params.length == 0
    return (
        <QueryBoxPaper>
            <QueryHead renderTitle={renderTitle} renderButtons={renderTitleButtons()}/>
            {children}
            <QueryBody show={!isNoParams && checkView[ViewCheckType.PARAMS]}>
                <FixedInputs placeholders={params ?? []} onChange={setParamsValues}/>
            </QueryBody>
            <QueryBody show={checkView[ViewCheckType.RUN]}>
                <QueryBodyRun
                    request={{queryUuid, credentialId, db, queryParams: paramsValues}}
                    varieties={varieties}
                />
            </QueryBody>
        </QueryBoxPaper>
    )

    function renderTitleButtons() {
        return (
            <>
                {renderButtons}
                {renderQueryParamsButton()}
                {renderRunButton()}
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
