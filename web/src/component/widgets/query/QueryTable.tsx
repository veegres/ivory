import {useMemo} from "react"

import {QueryConnection, QueryFields} from "../../../api/query/type"
import {ErrorSmart} from "../../view/box/ErrorSmart"
import {VirtualizedTable} from "../../view/table/VirtualizedTable"
import {QueryTableActions} from "./QueryTableActions"

type Props = {
    connection: QueryConnection,
    refetch: () => void,
    data?: QueryFields,
    error?: Error | null,
    loading?: boolean,
    height?: number,
    width?: number,
    showIndexColumn?: boolean,
}

export function QueryTable(props: Props) {
    const {data, error, loading, connection, refetch, showIndexColumn, height, width} = props

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
        return <QueryTableActions connection={connection} refetch={refetch} pid={row[pidIndex]}/>
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
