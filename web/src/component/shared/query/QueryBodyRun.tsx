import {Box, TableCell, Tooltip} from "@mui/material";
import {SxPropsMap} from "../../../type/common";
import {ErrorSmart} from "../../view/box/ErrorSmart";
import {QueryRunRequest, QueryVariety} from "../../../type/query";
import {CancelIconButton, RefreshIconButton, TerminateIconButton} from "../../view/button/IconButtons";
import {useMutation, useQuery} from "@tanstack/react-query";
import {queryApi} from "../../../app/api";
import {QueryVarieties} from "./QueryVarieties";
import {SimpleStickyTable} from "../../view/table/SimpleStickyTable";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    no: {
        display: "flex", alignItems: "center", justifyContent: "center", textTransform: "uppercase",
        padding: "8px 16px", border: "1px solid", borderColor: "divider", lineHeight: "1.54",
    },
    info: {display: "flex", alignItems: "center", justifyContent: "space-between", gap: 2},
    label: {color: "text.secondary", padding: "0 5px", cursor: "pointer", fontSize: "13.5px"},
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

    const result = useQuery({
        queryKey: ["query", "run", request.queryUuid ?? "standalone"],
        queryFn: () => queryApi.run(request.queryUuid ? {...request, query: undefined} : request),
        enabled: true, retry: false, refetchOnWindowFocus: false,
    })
    const {data, error, isFetching} = result

    const onSuccess = () => result.refetch()
    const cancel = useMutation({mutationFn: queryApi.cancel, onSuccess})
    const terminate = useMutation({mutationFn: queryApi.terminate, onSuccess})

    const pidIndex = result.data?.fields.findIndex(field => field.name === "pid") ?? -1

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
            return <Box sx={SX.no}>Response is empty</Box>
        }

        const columns = data?.fields.map(field => (
            {name: field.name, description: field.dataType}
        ))

        return (
            <SimpleStickyTable
                rows={data?.rows ?? []}
                columns={columns ?? []}
                fetching={isFetching}
                renderHeaderCell={() => pidIndex !== -1 && <TableCell/>}
                renderRowCell={renderRowButtons}
            />
        )
    }

    function renderRowButtons(row: any[]) {
        return pidIndex !== -1 && (
            <TableCell sx={SX.cell}>
                <Box sx={SX.pid}>
                    <CancelIconButton
                        size={25}
                        onClick={() => cancel.mutate({pid: row[pidIndex], credentialId, db})}
                    />
                    <TerminateIconButton
                        size={25}
                        onClick={() => terminate.mutate({pid: row[pidIndex], credentialId, db})}
                        color={"error"}
                    />
                </Box>
            </TableCell>
        )
    }
}
