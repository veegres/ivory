import {ClusterNoInstanceError, ClusterNoLeaderError, ClusterNoPostgresPassword} from "./OverviewError";
import {Box, Button, TextField} from "@mui/material";
import {useState} from "react";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {bloatApi, queryApi} from "../../../app/api";
import {SxPropsMap} from "../../../type/common";
import {Instance} from "../../../type/Instance";
import {Cluster} from "../../../type/cluster";
import {Bloat, BloatTarget} from "../../../type/bloat";
import {AutocompleteFetch} from "../../view/././autocomplete/AutocompleteFetch";
import {getDomain} from "../../../app/utils";

const SX: SxPropsMap = {
    form: {display: "grid", gridTemplateColumns: "repeat(3, 1fr)", columnGap: "30px"},
    buttons: {display: "flex", alignItems: "center", gap: 1}
}


type Props = {
    defaultInstance: Instance,
    cluster: Cluster,
    onClick: () => void,
    onSuccess: (job: Bloat) => void,
    target?: BloatTarget,
    setTarget: (target: BloatTarget) => void
}

export function OverviewBloatJobForm(props: Props) {
    const {defaultInstance, cluster, target, onSuccess, onClick, setTarget} = props
    const [ratio, setRadio] = useState<number>()

    const {onError} = useMutationOptions()
    const start = useMutation(bloatApi.start, {onSuccess, onError})

    if (!defaultInstance.inCluster) return <ClusterNoInstanceError/>
    if (!defaultInstance.leader) return <ClusterNoLeaderError/>
    if (!cluster.credentials.postgresId) return <ClusterNoPostgresPassword/>

    const req = {clusterName: cluster.name, db: {...defaultInstance.database, database: target?.dbName}}
    return (
        <Box sx={SX.form}>
            {renderInput(
                ["databases"],
                "Database",
                (v) => setTarget({...target, dbName: v}),
                (v) => queryApi.databases({...req, name: v})
            )}
            {renderInput(
                ["schemas"],
                "Schema",
                (v) => setTarget({...target, schema: v}),
                (v) => queryApi.schemas({...req, name: v})
            )}
            {renderInput(
                ["tables", target?.schema ?? ""],
                "Table",
                (v) => setTarget({...target, table: v}),
                (v) => queryApi.tables({...req, schema: target?.schema, name: v})
            )}
            {renderInput(
                ["schemas"],
                "Exclude Schema",
                (v) => setTarget({...target, excludeSchema: v}),
                (v) => queryApi.schemas({...req, name: v})
            )}
            {renderInput(
                ["tables", target?.excludeSchema ?? ""],
                "Exclude Table",
                (v) => setTarget({...target, excludeTable: v}),
                (v) => queryApi.tables({...req, schema: target?.excludeSchema, name: v})
            )}
            <Box sx={SX.buttons}>
                <TextField
                    sx={{flexGrow: 1}} size={"small"} label={"Ratio"} type={"number"} variant={"standard"}
                    onChange={(e) => setRadio(parseInt(e.target.value))}
                />
                <Box sx={SX.buttons}>
                    <Button variant={"text"} disabled={start.isLoading} onClick={handleRun}>
                        Start
                    </Button>
                </Box>
            </Box>
        </Box>
    )

    function renderInput(key: string[], label: string, onChange: (v: string) => void, onFetch: (value: string) => Promise<string[]>) {
        const [type, schema] = key
        const keys = ["query", type, req.clusterName, getDomain(req.db), req.db.database ?? ""]
        return (
            <AutocompleteFetch
                keys={schema ? [...keys, schema] : keys}
                onFetch={onFetch} label={label} margin={"dense"} variant={"standard"}
                noOptionsText={`Select previous value`}
                onUpdate={(v) => onChange(v || "")}
            />
        )
    }

    function handleRun() {
        if (defaultInstance && cluster.credentials?.postgresId) {
            const {database: {host, port}} = defaultInstance
            onClick()
            start.mutate({
                connection: {host, port, credId: cluster.credentials.postgresId},
                target,
                ratio,
                cluster: cluster.name
            })
        }
    }
}
