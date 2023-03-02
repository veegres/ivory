import {ClusterNoInstanceError, ClusterNoLeaderError, ClusterNoPostgresPassword} from "./OverviewError";
import {Box, Button, TextField} from "@mui/material";
import {Cluster, CompactTable, DefaultInstance, SxPropsMap, Target} from "../../../app/types";
import {ChangeEvent, useState} from "react";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {bloatApi} from "../../../app/api";

const SX: SxPropsMap = {
    form: {display: "grid", flexGrow: 1, gridTemplateColumns: "repeat(3, 1fr)", columnGap: "30px"},
    buttons: {display: "flex", alignItems: "center", gap: 1}
}


type Props = {
    defaultInstance: DefaultInstance,
    cluster: Cluster,
    disabled: boolean,
    onSuccess: (job: CompactTable) => void,
}

export function OverviewBloatJobForm(props: Props) {
    const {defaultInstance, cluster, onSuccess, disabled} = props
    const [target, setTarget] = useState<Target>()
    const [ratio, setRadio] = useState<number>()


    const {onError} = useMutationOptions()
    const start = useMutation(bloatApi.start, {onSuccess, onError})

    if (!defaultInstance.inCluster) return <ClusterNoInstanceError/>
    if (!defaultInstance.leader) return <ClusterNoLeaderError/>
    if (!cluster.credentials.postgresId) return <ClusterNoPostgresPassword/>

    return (
        <Box sx={SX.form}>
            {renderInput("Database Name", (e) => setTarget({...target, dbName: e.target.value}))}
            {renderInput("Schema", (e) => setTarget({...target, schema: e.target.value}))}
            {renderInput("Table", (e) => setTarget({...target, table: e.target.value}))}
            {renderInput("Exclude Schema", (e) => setTarget({...target, excludeSchema: e.target.value}))}
            {renderInput("Exclude Table", (e) => setTarget({...target, excludeTable: e.target.value}))}
            <Box sx={SX.buttons}>
                <TextField
                    sx={{flexGrow: 1}} size={"small"} label={"Ratio"} type={"number"} variant={"standard"}
                    disabled={disabled}
                    onChange={(e) => setRadio(parseInt(e.target.value))}
                />
                <Box sx={SX.buttons}>
                    <Button variant={"text"} disabled={start.isLoading || disabled} onClick={handleRun}>
                        RUN
                    </Button>
                </Box>
            </Box>
        </Box>
    )

    function renderInput(label: string, onChange: (e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => void) {
        return (
            <TextField
                size={"small"}
                label={label}
                variant={"standard"}
                margin={"dense"}
                disabled={disabled}
                onChange={onChange}
            />
        )
    }

    function handleRun() {
        if (defaultInstance && cluster.credentials?.postgresId) {
            const {database: {host, port}} = defaultInstance
            start.mutate({
                connection: {host, port, credId: cluster.credentials.postgresId},
                target,
                ratio,
                cluster: cluster.name
            })
        }
    }
}
