import {useMemo} from "react"

import {ErrorSmart} from "../../../../shared/component/box/ErrorSmart"
import {VirtualizedTable} from "../../../../shared/component/table/VirtualizedTable"
import {useRouterNodePlatformProcesses} from "../../api/hook"
import {PlatformVaultConnection, Process} from "../../api/type"

const COLUMNS = [
    {name: "PID", width: 70},
    {name: "PROGRAM", width: 140},
    {name: "COMMAND"},
    {name: "THREADS", width: 90},
    {name: "USER", width: 110},
    {name: "MEMB", width: 90},
    {name: "CPU%", width: 70},
]

type Props = {
    connection: PlatformVaultConnection,
}

export function PlatformProcesses(props: Props) {
    const {connection} = props
    const processes = useRouterNodePlatformProcesses(connection)

    const rows = useMemo(handleMemoRows, [processes.data])

    if (processes.error) return <ErrorSmart error={processes.error}/>

    return (
        <VirtualizedTable
            columns={COLUMNS}
            rows={rows}
            loading={processes.isLoading}
            showIndexColumn={false}
        />
    )

    function handleMemoRows() {
        return (processes.data ?? []).map(toRow)
    }
}

function toRow(process: Process) {
    return [
        process.pid,
        process.program,
        process.command,
        process.threads,
        process.user,
        formatBytes(process.memoryBytes),
        process.cpuPercent.toFixed(1),
    ]
}

function formatBytes(bytes: number) {
    const units = ["B", "KB", "MB", "GB", "TB"]
    let value = bytes
    let unitIndex = 0
    while (value >= 1024 && unitIndex < units.length - 1) {
        value /= 1024
        unitIndex++
    }
    return `${value.toFixed(1)} ${units[unitIndex]}`
}
