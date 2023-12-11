import {Box} from "@mui/material";
import {
    CancelIconButton,
    DeleteIconButton,
    EditIconButton,
    InfoIconButton,
    PlayIconButton,
    QueryParamsIconButton,
    RestoreIconButton
} from "../../view/button/IconButtons";
import {useAppearance} from "../../../provider/AppearanceProvider";
import {useState} from "react";
import {QueryItemInfo} from "./QueryItemInfo";
import {QueryItemEdit} from "./QueryItemEdit";
import {QueryItemRestore} from "./QueryItemRestore";
import {QueryItemBody} from "./QueryItemBody";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {useMutation, useQuery} from "@tanstack/react-query";
import {queryApi} from "../../../app/api";
import {QueryItemRun} from "./QueryItemRun";
import {Database, SxPropsMap} from "../../../type/common";
import {Query, QueryCreation, QueryType} from "../../../type/query";
import {FixedInputs} from "../../view/input/FixedInputs";
import {QueryItemHead} from "./QueryItemHead";
import {QueryItemPaper} from "./QueryItemPaper";

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
    const {query, credentialId, db, type, editable} = props
    const {info} = useAppearance()
    const [toggleView, setToggleView] = useState<ViewToggleType>()
    const [checkView, setCheckView] = useState<{[key in ViewCheckType]: boolean}>({[ViewCheckType.RUN]: false, [ViewCheckType.PARAMS]: false})
    const [params, setParams] = useState(query.params?.map(() => ''))
    const [currentDb, setCurrentDb] = useState(db)
    const open = toggleView !== undefined

    const result = useQuery({
        queryKey: ["query", "run", query.id],
        queryFn: () => {
            setCurrentDb(db)
            return queryApi.run({queryUuid: query.id, credentialId, db, queryParams: params})
        },
        enabled: false, retry: false,
    })

    const removeOptions = useMutationOptions([["query", "map", type]])
    const remove = useMutation({mutationFn: queryApi.delete, ...removeOptions})
    const cancel = useMutation({mutationFn: queryApi.cancel, onSuccess: () => result.refetch()})
    const terminate = useMutation({mutationFn: queryApi.terminate, onSuccess: () => result.refetch()})

    return (
        <QueryItemPaper>
            <QueryItemHead renderTitle={renderTitle()} renderButtons={renderTitleButtons()} onClick={handleCloseAll}/>
            <QueryItemBody show={checkView[ViewCheckType.PARAMS]}>
                <FixedInputs inputs={params ?? []} placeholders={query.params ?? []} onChange={v => setParams(v)}/>
            </QueryItemBody>
            <QueryItemBody show={toggleView === ViewToggleType.INFO}>
                <QueryItemInfo query={query.custom} description={query.description} varieties={query.varieties}/>
            </QueryItemBody>
            <QueryItemBody show={toggleView === ViewToggleType.EDIT}>
                <QueryItemEdit id={query.id} query={query.custom}/>
            </QueryItemBody>
            <QueryItemBody show={toggleView === ViewToggleType.RESTORE}>
                <QueryItemRestore id={query.id} def={query.default} custom={query.custom}/>
            </QueryItemBody>
            <QueryItemBody show={checkView[ViewCheckType.RUN]}>
                <QueryItemRun
                    data={result.data}
                    error={result.error}
                    varieties={query.varieties}
                    db={currentDb}
                    loading={result.isFetching || cancel.isPending || terminate.isPending}
                    onRefresh={() => result.refetch().then()}
                    onCancel={(pid) => cancel.mutate({pid, credentialId, db})}
                    onTerminate={(pid) => terminate.mutate({pid, credentialId, db})}
                />
            </QueryItemBody>
        </QueryItemPaper>
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
                <InfoIconButton onClick={handleToggleBody(ViewToggleType.INFO)}/>
                {query.default !== query.custom && (
                    <RestoreIconButton onClick={handleToggleBody(ViewToggleType.RESTORE)}/>
                )}
                {editable && (
                    <EditIconButton onClick={handleToggleBody(ViewToggleType.EDIT)}/>
                )}
                {query.creation === QueryCreation.Manual && (
                    <DeleteIconButton loading={remove.isPending} onClick={handleDelete}/>
                )}
            </>
        )
    }

    function renderQueryParamsButton() {
        const disabled = !query.params || query.params.length == 0
        return !checkView[ViewCheckType.PARAMS] ? (
            <QueryParamsIconButton color={"secondary"} disabled={disabled} onClick={() => handleCheckView(ViewCheckType.PARAMS)}/>
        ) : (
            <CancelIconButton color={"secondary"} tooltip={"Close Query Params"} onClick={() => handleCheckView(ViewCheckType.PARAMS)}/>
        )
    }

    function renderRunButton() {
        return !checkView[ViewCheckType.RUN] ? (
            <PlayIconButton color={"success"} loading={result.isFetching} onClick={handleRun}/>
        ) : (
            <CancelIconButton color={"success"} tooltip={"Close"} onClick={() => handleCheckView(ViewCheckType.RUN)}/>
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

    function handleRun() {
        if (!checkView[ViewCheckType.RUN]) handleCheckView(ViewCheckType.RUN)
        result.refetch().then()
    }

    function handleCheckView(type: ViewCheckType) {
        setCheckView({...checkView, [type]: !checkView[type]})
    }

    function handleDelete() {
        remove.mutate(query.id)
    }
}
