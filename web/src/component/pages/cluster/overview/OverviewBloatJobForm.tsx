import {Box, Button, Checkbox, FormControlLabel, TextField} from "@mui/material"
import {useState} from "react"

import {useRouterBloatStart} from "../../../../api/bloat/hook"
import {BloatOptions, BloatTarget} from "../../../../api/bloat/type"
import {Cluster, Node} from "../../../../api/cluster/type"
import {DatabaseType} from "../../../../api/database/type"
import {useRouterQueryDatabase, useRouterQuerySchemas, useRouterQueryTables} from "../../../../api/query/hook"
import {SxPropsMap} from "../../../../app/type"
import {getQueryConnection} from "../../../../app/utils"
import {AutocompleteFetch} from "../../../view/autocomplete/AutocompleteFetch"
import {ClusterNoLeaderError, ClusterNoPostgresVault, NoDatabaseError} from "./OverviewError"

const SX: SxPropsMap = {
    form: {display: "grid", gridTemplateColumns: "repeat(4, 1fr)", columnGap: "30px"},
    group: {display: "flex", alignItems: "center", justifyContent: "space-between", gap: 1},
}

type Props = {
    node: Node,
    cluster: Cluster,
    onClick: () => void,
    target?: BloatTarget,
    setTarget: (target: BloatTarget) => void
}

export function OverviewBloatJobForm(props: Props) {
    const {node, cluster, target, onClick, setTarget} = props
    const [options, setOptions] = useState<BloatOptions>({force: false, noReindex: false, routineVacuum: false, initialReindex: false, noInitialVacuum: false})

    const start = useRouterBloatStart(cluster.name)

    if (node.keeper.role !== "leader") return <ClusterNoLeaderError/>
    const vaultId = cluster.vaults.databaseId
    if (!vaultId) return <ClusterNoPostgresVault/>
    if (!node.connection.host || !node.connection.dbPort) return <NoDatabaseError/>

    const db = {type: DatabaseType.POSTGRES, host: node.connection.host, port: node.connection.dbPort, name: target?.database}
    const connection = getQueryConnection(cluster, db)
    return (
        <Box sx={SX.form}>
            <AutocompleteFetch
                value={target?.database || null}
                margin={"dense"} variant={"standard"} size={"small"}
                label={"Database"}
                connection={connection}
                useFetch={useRouterQueryDatabase}
                onUpdate={(v) => setTarget({database: v || undefined})}
            />
            <AutocompleteFetch
                value={target?.schema || null}
                margin={"dense"} variant={"standard"} size={"small"}
                label={"Schema"}
                connection={connection}
                useFetch={useRouterQuerySchemas}
                onUpdate={(v) => setTarget({database: target?.database, schema: v || undefined})}
                disabled={!target?.database || !!target?.excludeSchema}
            />
            <AutocompleteFetch
                value={target?.table || null}
                margin={"dense"} variant={"standard"} size={"small"}
                label={"Table"}
                connection={connection}
                params={{schema: target?.schema ?? ""}}
                useFetch={useRouterQueryTables}
                onUpdate={(v) => setTarget({...target, table: v || undefined})}
                disabled={!target?.schema || !!target?.excludeTable}
            />
            <Box sx={SX.group}>
                <TextField
                    label={"Min table size (MB)"} type={"number"} variant={"standard"} size={"small"}
                    onChange={(e) => setOptions({...options, minTableSize: parseInt(e.target.value)})}
                />
                <TextField
                    label={"Max table size (MB)"} type={"number"} variant={"standard"} size={"small"}
                    onChange={(e) => setOptions({...options, maxTableSize: parseInt(e.target.value)})}
                />
            </Box>
            <AutocompleteFetch
                value={target?.excludeSchema || null}
                margin={"dense"} variant={"standard"} size={"small"}
                label={"Exclude Schema"}
                connection={connection}
                useFetch={useRouterQuerySchemas}
                onUpdate={(v) => setTarget({database: target?.database, excludeSchema: v || undefined})}
                disabled={!target?.database || !!target?.schema}
            />
            <AutocompleteFetch
                value={target?.excludeTable || null}
                margin={"dense"} variant={"standard"} size={"small"}
                label={"Exclude Table"}
                connection={connection}
                params={{schema: target?.schema ?? ""}}
                useFetch={useRouterQueryTables}
                onUpdate={(v) => setTarget({...target, excludeTable: v || undefined})}
                disabled={!target?.schema || !!target?.table}
            />
            <Box sx={SX.group}>
                <TextField
                    label={"Delay ratio"} type={"number"} variant={"standard"} size={"small"}
                    onChange={(e) => setOptions({...options, delayRatio: parseInt(e.target.value)})}
                />
            </Box>
            <Box sx={SX.group}>
                <FormControlLabel label={"Force"} control={<Checkbox checked={options.force} onChange={(e) => setOptions({...options, force: e.target.checked})}/>}/>
                <Button variant={"text"} disabled={start.isPending} onClick={handleRun}>Clean</Button>
            </Box>
            <FormControlLabel label={"Routing vacuum"} control={<Checkbox checked={options.routineVacuum} onChange={(e) => setOptions({...options, routineVacuum: e.target.checked})}/>}/>
            <FormControlLabel label={"No initial vacuum"} control={<Checkbox checked={options.noInitialVacuum} onChange={(e) => setOptions({...options, noInitialVacuum: e.target.checked})}/>}/>
            <FormControlLabel label={"Initial reindex"} disabled={options.noReindex} control={<Checkbox checked={options.initialReindex} onChange={(e) => setOptions({...options, initialReindex: e.target.checked})}/>}/>
            <FormControlLabel label={"No reindex"} disabled={options.initialReindex} control={<Checkbox checked={options.noReindex} onChange={(e) => setOptions({...options, noReindex: e.target.checked})}/>}/>
        </Box>
    )

    function handleRun() {
        if (node && vaultId) {
            const dbRun = {...db, schema: target?.schema}
            onClick()
            start.mutate({
                connection: getQueryConnection(cluster, dbRun),
                target, options, cluster: cluster.name,
            })
        }
    }
}
