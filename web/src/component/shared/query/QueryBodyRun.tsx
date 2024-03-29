import {Box, TableCell, Tooltip} from "@mui/material";
import {SxPropsMap} from "../../../type/general";
import {ErrorSmart} from "../../view/box/ErrorSmart";
import {QueryRunRequest, QueryVariety} from "../../../type/query";
import {CancelIconButton, RefreshIconButton, TerminateIconButton} from "../../view/button/IconButtons";
import {QueryVarieties} from "./QueryVarieties";
import {SimpleStickyTable} from "../../view/table/SimpleStickyTable";
import {useCallback, useMemo} from "react";
import {NoBox} from "../../view/box/NoBox";
import {useRouterQueryCancel, useRouterQueryRun, useRouterQueryTerminate} from "../../../router/query";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    info: {display: "flex", alignItems: "center", justifyContent: "space-between", gap: 2},
    label: {color: "text.secondary", cursor: "pointer", fontSize: "13.5px"},
    buttons: {display: "flex", alignItems: "center", gap: 1},
    pid: {display: "flex"},
    cell: {padding: "0 1px"},
}

type Props = {
    varieties?: QueryVariety[],
    request: QueryRunRequest,
}

export function QueryBodyRun(props: Props) {
    const {varieties, request} = props
    const {credentialId, db} = request

    const result = useRouterQueryRun(request)
    const {data, error, isFetching} = result

    const cancel = useRouterQueryCancel(request.queryUuid)
    const terminate = useRouterQueryTerminate(request.queryUuid)

    // NOTE: we need this only to fix linter when passing to callback function
    const cancelMutate = cancel.mutate
    const terminateMutate = terminate.mutate

    const pidIndex = useMemo(handleMemoPidIndex, [data])
    const renderRowButtons = useCallback(handleCallbackRenderRowButtons, [cancelMutate, terminateMutate, credentialId, db, pidIndex])
    const renderHeaderCell = useCallback(handleCallbackRenderHeaderCell, [pidIndex])
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
                    <Box sx={SX.label}>[ {isFetching ? "request is loading" : data?.url ?? "unknown"} ]</Box>
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
                        disabled={isFetching || cancel.isPending || terminate.isPending}
                        onClick={result.refetch}
                    />
                </Box>
            </Box>
        )
    }

    function renderTable() {
        if (error) return <ErrorSmart error={error}/>
        if (!isFetching && (!data || !data.fields.length || !data.rows.length)) {
            return <NoBox text={"Response is empty"}/>
        }

        return (
            <SimpleStickyTable
                rows={data?.rows ?? []}
                columns={columns}
                fetching={isFetching}
                renderHeaderCell={renderHeaderCell}
                renderRowCell={renderRowButtons}
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

    function handleCallbackRenderHeaderCell() {
        return pidIndex !== -1 && <TableCell/>
    }

    function handleCallbackRenderRowButtons(row: any[]) {
        return pidIndex !== -1 && (
            <TableCell sx={SX.cell}>
                <Box sx={SX.pid}>
                    <CancelIconButton
                        size={25}
                        onClick={() => cancelMutate({pid: row[pidIndex], credentialId, db})}
                    />
                    <TerminateIconButton
                        size={25}
                        onClick={() => terminateMutate({pid: row[pidIndex], credentialId, db})}
                        color={"error"}
                    />
                </Box>
            </TableCell>
        )
    }
}
