import {ClusterNoInstanceError, ClusterNoLeaderError, ClusterNoPostgresPassword} from "./OverviewError";
import {Box, Button, TextField} from "@mui/material";
import {useState} from "react";
import {useMutationOptions} from "../../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {BloatApi, QueryApi} from "../../../../app/api";
import {SxPropsMap} from "../../../../type/common";
import {InstanceWeb} from "../../../../type/instance";
import {Cluster} from "../../../../type/cluster";
import {Bloat, BloatTarget} from "../../../../type/bloat";
import {AutocompleteFetch} from "../../../view/autocomplete/AutocompleteFetch";
import {getDomain} from "../../../../app/utils";
import {QueryPostgresRequest} from "../../../../type/query";

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

    const {queryClient, onError} = useMutationOptions()
    const start = useMutation({
        mutationFn: BloatApi.start,
        onSuccess: (job) => {
            queryClient.setQueryData<Bloat[]>(
                ["instance", "bloat", "list", cluster.name],
                (jobs) => [job, ...(jobs ?? [])]
            )
        },
        onError
    })

    if (!defaultInstance.inCluster) return <ClusterNoInstanceError/>
    if (!defaultInstance.leader) return <ClusterNoLeaderError/>
    if (!cluster.credentials.postgresId) return <ClusterNoPostgresPassword/>

    const postgresId = cluster.credentials.postgresId
    const req: QueryPostgresRequest = {credentialId: postgresId, db: {...defaultInstance.database, name: target?.dbName}}
    const keys = [getDomain(req.db), req.db.name ?? "postgres"]
    return (
        <Box sx={SX.form}>
            <AutocompleteFetch
                margin={"dense"} variant={"standard"}
                keys={["databases", ...keys]} label={"Database"}
                onFetch={(v) => QueryApi.databases({...req, name: v})}
                onUpdate={(v) => setTarget({...target, dbName: v || ""})}
            />
            <AutocompleteFetch
                margin={"dense"} variant={"standard"}
                keys={["schemas", ...keys]} label={"Schema"}
                disabled={!target?.dbName || !!target?.excludeSchema}
                onFetch={(v) => QueryApi.schemas({...req, name: v})}
                onUpdate={(v) => setTarget({...target, schema: v || ""})}
            />
            <AutocompleteFetch
                margin={"dense"} variant={"standard"}
                keys={["schemas", ...keys]} label={"Exclude Schema"}
                disabled={!target?.dbName || !!target?.schema}
                onFetch={(v) => QueryApi.schemas({...req, name: v})}
                onUpdate={(v) => setTarget({...target, excludeSchema: v || ""})}
            />
            <AutocompleteFetch
                margin={"dense"} variant={"standard"}
                keys={["tables", ...keys, target?.schema ?? ""]} label={"Table"}
                disabled={!target?.schema || !!target?.excludeTable}
                onFetch={(v) => QueryApi.tables({...req, schema: target?.schema ?? "", name: v})}
                onUpdate={(v) => setTarget({...target, table: v || ""})}
            />
            <AutocompleteFetch
                margin={"dense"} variant={"standard"}
                keys={["tables", ...keys, target?.excludeSchema ?? ""]} label={"Exclude Table"}
                disabled={!target?.schema || !!target?.table}
                onFetch={(v) => QueryApi.tables({...req, schema: target?.schema ?? "", name: v})}
                onUpdate={(v) => setTarget({...target, excludeTable: v || ""})}
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
        if (defaultInstance && postgresId) {
            const {database: {host, port}} = defaultInstance
            onClick()
            start.mutate({
                connection: {host, port, credId: postgresId},
                target,
                ratio,
                cluster: cluster.name
            })
        }
    }
}
