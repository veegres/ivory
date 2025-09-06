import {ClusterNoInstanceError, ClusterNoLeaderError, ClusterNoPostgresPassword} from "./OverviewError";
import {Box, Button, Checkbox, FormControlLabel, TextField} from "@mui/material";
import {useState} from "react";
import {SxPropsMap} from "../../../../api/management/type";
import {InstanceWeb} from "../../../../api/instance/type";
import {Cluster} from "../../../../api/cluster/type";
import {BloatOptions, BloatTarget} from "../../../../api/bloat/type";
import {AutocompleteFetch} from "../../../view/autocomplete/AutocompleteFetch";
import {useRouterBloatStart} from "../../../../api/bloat/hook";
import {getQueryConnection} from "../../../../app/utils";
import {useRouterQueryDatabase, useRouterQuerySchemas, useRouterQueryTables} from "../../../../api/query/hook";

const SX: SxPropsMap = {
    form: {display: "grid", gridTemplateColumns: "repeat(4, 1fr)", columnGap: "30px"},
    group: {display: "flex", alignItems: "center", justifyContent: "space-between", gap: 1},
}

type Props = {
    defaultInstance: InstanceWeb,
    cluster: Cluster,
    onClick: () => void,
    target?: BloatTarget,
    setTarget: (target: BloatTarget) => void
}

export function OverviewBloatJobForm(props: Props) {
    const {defaultInstance, cluster, target, onClick, setTarget} = props
    const [options, setOptions] = useState<BloatOptions>({force: false, noReindex: false, routineVacuum: false, initialReindex: false, noInitialVacuum: false})

    const start = useRouterBloatStart(cluster.name)

    if (!defaultInstance.inCluster) return <ClusterNoInstanceError/>
    if (!defaultInstance.leader) return <ClusterNoLeaderError/>
    if (!cluster.credentials.postgresId) return <ClusterNoPostgresPassword/>

    const db = {...defaultInstance.database, name: target?.dbName}
    const connection = getQueryConnection(cluster, db)
    return (
        <Box sx={SX.form}>
            <AutocompleteFetch
                value={target?.dbName || null}
                margin={"dense"} variant={"standard"} size={"small"}
                label={"Database"}
                connection={connection}
                useFetch={useRouterQueryDatabase}
                onUpdate={(v) => setTarget({dbName: v || undefined})}
            />
            <AutocompleteFetch
                value={target?.schema || null}
                margin={"dense"} variant={"standard"} size={"small"}
                label={"Schema"}
                connection={connection}
                useFetch={useRouterQuerySchemas}
                onUpdate={(v) => setTarget({dbName: target?.dbName, schema: v || undefined})}
                disabled={!target?.dbName || !!target?.excludeSchema}
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
                onUpdate={(v) => setTarget({dbName: target?.dbName, excludeSchema: v || undefined})}
                disabled={!target?.dbName || !!target?.schema}
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
        const credentialId = cluster.credentials.postgresId
        if (defaultInstance && credentialId) {
            const {database: {host, port}} = defaultInstance
            onClick()
            start.mutate({
                connection: {host, port, credentialId},
                target, options, cluster: cluster.name,
            })
        }
    }
}
