import {ClusterNoInstanceError, ClusterNoLeaderError, ClusterNoPostgresPassword} from "./OverviewError";
import {Box, Button, TextField} from "@mui/material";
import {useState} from "react";
import {SxPropsMap} from "../../../../type/general";
import {InstanceWeb} from "../../../../type/instance";
import {Cluster} from "../../../../type/cluster";
import {BloatTarget} from "../../../../type/bloat";
import {AutocompleteFetch} from "../../../view/autocomplete/AutocompleteFetch";
import {useRouterBloatStart} from "../../../../router/bloat";
import {getQueryConnection} from "../../../../app/utils";
import {useRouterQueryDatabase, useRouterQuerySchemas, useRouterQueryTables} from "../../../../router/query";

const SX: SxPropsMap = {
    form: {display: "grid", gridTemplateColumns: "repeat(3, 1fr)", columnGap: "30px"},
    buttons: {display: "flex", alignItems: "center", gap: 1}
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
    const [ratio, setRadio] = useState<number>()

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
                margin={"dense"} variant={"standard"}
                label={"Database"}
                connection={connection}
                useFetch={useRouterQueryDatabase}
                onUpdate={(v) => setTarget({dbName: v || undefined})}
            />
            <AutocompleteFetch
                value={target?.schema || null}
                margin={"dense"} variant={"standard"}
                label={"Schema"}
                connection={connection}
                useFetch={useRouterQuerySchemas}
                onUpdate={(v) => setTarget({dbName: target?.dbName, schema: v || undefined})}
                disabled={!target?.dbName || !!target?.excludeSchema}
            />
            <AutocompleteFetch
                value={target?.excludeSchema || null}
                margin={"dense"} variant={"standard"}
                label={"Exclude Schema"}
                connection={connection}
                useFetch={useRouterQuerySchemas}
                onUpdate={(v) => setTarget({dbName: target?.dbName, excludeSchema: v || undefined})}
                disabled={!target?.dbName || !!target?.schema}
            />
            <AutocompleteFetch
                value={target?.table || null}
                margin={"dense"} variant={"standard"}
                label={"Table"}
                connection={connection}
                params={{schema: target?.schema ?? ""}}
                useFetch={useRouterQueryTables}
                onUpdate={(v) => setTarget({...target, table: v || undefined})}
                disabled={!target?.schema || !!target?.excludeTable}
            />
            <AutocompleteFetch
                value={target?.excludeTable || null}
                margin={"dense"} variant={"standard"}
                label={"Exclude Table"}
                connection={connection}
                params={{schema: target?.schema ?? ""}}
                useFetch={useRouterQueryTables}
                onUpdate={(v) => setTarget({...target, excludeTable: v || undefined})}
                disabled={!target?.schema || !!target?.table}
            />
            <Box sx={SX.buttons}>
                <TextField
                    sx={{flexGrow: 1}} size={"small"} label={"Ratio"} type={"number"} variant={"standard"}
                    onChange={(e) => setRadio(parseInt(e.target.value))}
                />
                <Box sx={SX.buttons}>
                    <Button variant={"text"} disabled={start.isPending} onClick={handleRun}>
                        Clean
                    </Button>
                </Box>
            </Box>
        </Box>
    )

    function handleRun() {
        const credentialId = cluster.credentials.postgresId
        if (defaultInstance && credentialId) {
            const {database: {host, port}} = defaultInstance
            onClick()
            start.mutate({
                connection: {host, port, credentialId},
                target,
                ratio,
                cluster: cluster.name
            })
        }
    }
}
