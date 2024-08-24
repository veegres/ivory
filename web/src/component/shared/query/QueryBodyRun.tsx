import {Box, Tooltip} from "@mui/material";
import {SxPropsMap} from "../../../type/general";
import {ErrorSmart} from "../../view/box/ErrorSmart";
import {QueryRunRequest, QueryVariety} from "../../../type/query";
import {CancelIconButton, RefreshIconButton, TerminateIconButton} from "../../view/button/IconButtons";
import {QueryVarieties} from "./QueryVarieties";
import {VirtualizedTable} from "../../view/table/VirtualizedTable";
import {useCallback, useMemo} from "react";
import {useRouterQueryCancel, useRouterQueryRun, useRouterQueryTerminate} from "../../../router/query";
import {getPostgresUrl, SxPropsFormatter} from "../../../app/utils";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    info: {display: "flex", alignItems: "center", justifyContent: "space-between", gap: 2},
    label: {color: "text.secondary", cursor: "pointer", fontSize: "13.5px", whiteSpace: "nowrap"},
    word: {whiteSpace: "wrap", wordBreak: "break-all"},
    buttons: {display: "flex", alignItems: "center", gap: 1},
    pid: {display: "flex", justifyContent: "space-evenly", color: "text.secondary", padding: "0 3px"},
    cell: {padding: "0 1px"},
}

type Props = {
    varieties?: QueryVariety[],
    request: QueryRunRequest,
}

export function QueryBodyRun(props: Props) {
    const {varieties, request} = props
    const {connection} = request

    const {data, error, isFetching, refetch}  = useRouterQueryRun(request)
    const cancel = useRouterQueryCancel(request.queryUuid)
    const terminate = useRouterQueryTerminate(request.queryUuid)

    // NOTE: we need this only to fix linter when passing to callback function
    const cancelMutate = cancel.mutate
    const terminateMutate = terminate.mutate

    const pidIndex = useMemo(handleMemoPidIndex, [data])
    const renderRowButtons = useCallback(handleCallbackRenderRowButtons, [cancelMutate, terminateMutate, connection, pidIndex])
    const columns = useMemo(handleMemoColumns, [data])

    return (
        <Box sx={SX.box}>
            {renderInfo()}
            {renderTable()}
        </Box>
    )

    function renderInfo() {
        return (
            <Box sx={SX.info}>
                <Tooltip title={"SENT TO"} placement={"right"} arrow={true}>
                    <Box sx={SxPropsFormatter.merge(SX.label, SX.word)}>
                        [ {getPostgresUrl(connection)} ]
                    </Box>
                </Tooltip>
                <Box sx={SX.buttons}>
                    {!isFetching && data && (
                        <Tooltip title={"DURATION"} placement={"left"} arrow={true}>
                            <Box sx={SX.label}>[ {(data.endTime - data.startTime) / 1000}s ]</Box>
                        </Tooltip>
                    )}
                    {varieties && <QueryVarieties varieties={varieties}/>}
                    <RefreshIconButton
                        color={"success"}
                        loading={isFetching}
                        disabled={cancel.isPending || terminate.isPending}
                        onClick={() => refetch()}
                    />
                </Box>
            </Box>
        )
    }

    function renderTable() {
        if (error) return <ErrorSmart error={error}/>
        const rows = data?.rows ?? []

        return (
            <VirtualizedTable
                rows={rows}
                columns={columns}
                fetching={isFetching}
                renderRowCell={pidIndex === -1 ? undefined : renderRowButtons}
            />
        )
    }

    function handleMemoPidIndex() {
        if (!data) return -1
        return data.fields.findIndex(field => field.name === "pid")
    }

    function handleMemoColumns() {
        if (!data) return []
        return data.fields.map(field => ({name: field.name, description: field.dataType}))
    }

    function handleCallbackRenderRowButtons(row: any[]) {
        return (
            <Box sx={SX.pid}>
                <CancelIconButton
                    color={"inherit"}
                    size={25}
                    onClick={() => cancelMutate({connection, pid: row[pidIndex]})}
                />
                <TerminateIconButton
                    size={25}
                    onClick={() => terminateMutate({connection, pid: row[pidIndex]})}
                    color={"error"}
                />
            </Box>
        )
    }
}
