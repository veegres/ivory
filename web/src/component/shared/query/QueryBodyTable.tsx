import {ErrorSmart} from "../../view/box/ErrorSmart";
import {VirtualizedTable} from "../../view/table/VirtualizedTable";
import {QueryConnection, QueryFields} from "../../../type/query";
import {useMemo} from "react";
import {QueryBodyTableKillButton} from "./QueryBodyTableKillButton";

type Props = {
    connection: QueryConnection,
    queryUuid?: string,
    data?: QueryFields,
    error?: Error | null,
    loading?: boolean,
    height?: number,
    width?: number,
    showIndexColumn?: boolean,
}

export function QueryBodyTable(props: Props) {
    const {data, error, loading, connection, queryUuid, showIndexColumn, height, width} = props

    const pidIndex = useMemo(handleMemoPidIndex, [data])
    const columns = useMemo(handleMemoColumns, [data])

    if (error) return <ErrorSmart error={error}/>
    const rows = data?.rows ?? []

    return (
        <VirtualizedTable
            rows={rows}
            columns={columns}
            loading={loading}
            height={height}
            width={width}
            showIndexColumn={showIndexColumn}
            renderRowActions={pidIndex === -1 ? undefined : renderActions}
        />
    )

    function renderActions(row: any[]) {
        return <QueryBodyTableKillButton connection={connection} queryUuid={queryUuid} pid={row[pidIndex]}/>
    }

    function handleMemoPidIndex() {
        if (!data) return -1
        return data.fields.findIndex(field => field.name === "pid")
    }

    function handleMemoColumns() {
        if (!data) return []
        return data.fields.map(field => ({name: field.name, description: field.dataType}))
    }
}
