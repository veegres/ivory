import {Box} from "@mui/material";
import {QueryFields} from "../../../type/query";
import {SxPropsMap} from "../../../type/common";
import {QueryBodyHistoryItem} from "./QueryBodyHistoryItem";
import {ClearAllIconButton} from "../../view/button/IconButtons";
import {useQuery} from "@tanstack/react-query";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    header: {display: "flex", justifyContent: "space-between", alignItems: "center", padding: "0 15px"},
    no: {
        display: "flex", alignItems: "center", justifyContent: "center", textTransform: "uppercase",
        padding: "8px 16px", border: "1px solid", borderColor: "divider", lineHeight: "1.54",
    },
}

type Props = {
    queryId: string,
}

const mock: QueryFields = {
    "fields": [
        {
            "name": "pid",
            "dataType": "int4",
            "dataTypeOID": 23
        },
        {
            "name": "state",
            "dataType": "text",
            "dataTypeOID": 25
        },
        {
            "name": "event",
            "dataType": "text",
            "dataTypeOID": 25
        },
        {
            "name": "transaction_duration",
            "dataType": "text",
            "dataTypeOID": 25
        },
        {
            "name": "query_duration",
            "dataType": "text",
            "dataTypeOID": 25
        },
        {
            "name": "query",
            "dataType": "text",
            "dataTypeOID": 25
        },
        {
            "name": "username",
            "dataType": "name",
            "dataTypeOID": 19
        },
        {
            "name": "application",
            "dataType": "text",
            "dataTypeOID": 25
        },
        {
            "name": "ip",
            "dataType": "inet",
            "dataTypeOID": 869
        },
        {
            "name": "blocked_by_process_id",
            "dataType": "_int4",
            "dataTypeOID": 1007
        }
    ],
    "rows": [
        [
            38,
            null,
            "AutoVacuumMain",
            "01:50:48.906087",
            null,
            "",
            null,
            "",
            null,
            []
        ],
        [
            40,
            null,
            "LogicalLauncherMain",
            "01:50:48.905375",
            null,
            "",
            "postgres",
            "",
            null,
            []
        ],
        [
            36,
            null,
            "BgWriterHibernate",
            "01:50:48.906901",
            null,
            "",
            null,
            "",
            null,
            []
        ],
        [
            35,
            null,
            "CheckpointerMain",
            "01:50:48.907261",
            null,
            "",
            null,
            "",
            null,
            []
        ],
        [
            37,
            null,
            "WalWriterMain",
            "01:50:48.906504",
            null,
            "",
            null,
            "",
            null,
            []
        ],
        [
            56,
            "active",
            "WalSenderMain",
            "01:50:33.649932",
            "01:50:33.638405",
            "START_REPLICATION SLOT \"patroni1\" 0/4000000 TIMELINE 1",
            "replicator",
            "patroni1",
            "172.19.0.3/32",
            []
        ],
        [
            57,
            "active",
            "WalSenderMain",
            "01:50:33.649378",
            "01:50:33.638405",
            "START_REPLICATION SLOT \"patroni3\" 0/3000000 TIMELINE 1",
            "replicator",
            "patroni3",
            "172.19.0.4/32",
            []
        ],
        [
            44,
            "idle",
            "ClientRead",
            "01:50:48.853267",
            "00:00:08.703696",
            "SELECT CASE WHEN pg_catalog.pg_is_in_recovery() THEN 0 ELSE ('x' || pg_catalog.substr(pg_catalog.pg_walfile_name(pg_catalog.pg_current_wal_lsn()), 1, 8))::bit(32)::int END, CASE WHEN pg_catalog.pg_is_in_recovery() THEN 0 ELSE pg_catalog.pg_wal_lsn_diff(pg_catalog.pg_current_wal_lsn(), '0/0')::bigint END, pg_catalog.pg_wal_lsn_diff(pg_catalog.pg_last_wal_replay_lsn(), '0/0')::bigint, pg_catalog.pg_wal_lsn_diff(COALESCE(pg_catalog.pg_last_wal_receive_lsn(), '0/0'), '0/0')::bigint, pg_catalog.pg_is_in_recovery() AND pg_catalog.pg_is_wal_replay_paused(), 0, CASE WHEN latest_end_lsn IS NULL THEN NULL ELSE received_tli END, slot_name, conninfo, NULL, 'on', '', NULL FROM pg_catalog.pg_stat_get_wal_receiver()",
            "postgres",
            "Patroni",
            "172.19.0.2/32",
            []
        ],
        [
            498,
            "active",
            null,
            "00:00:00.015404",
            "-00:00:00.001028",
            "SELECT\n    pid,\n    state,\n    wait_event AS event,\n    (now() - pg_stat_activity.backend_start)::text AS transaction_duration,\n    (now() - pg_stat_activity.query_start)::text AS query_duration,\n    query,\n    usename AS username,\n    application_name AS application,\n    client_addr AS ip,\n    pg_blocking_pids(pid) AS blocked_by_process_id\nFROM pg_stat_activity\nORDER BY now() - pg_stat_activity.query_start DESC;",
            "postgres",
            "",
            "172.19.0.5/32",
            []
        ]
    ],
    "startTime": 1703702116846,
    "endTime": 1703702116867,
    "url": "postgres://patroni2:5002/postgres"
}

export function QueryBodyHistory(props: Props) {
    const {queryId} = props

    const result = useQuery({
        queryKey: ["query", "history", queryId],
        queryFn: async () => {
            console.log("queryId", queryId)
            return Promise.resolve([mock, mock, mock, mock, mock]);
        },
        enabled: true, retry: false, refetchOnWindowFocus: false,
    })

    return (
        <Box sx={SX.box}>
            <Box sx={SX.header}>
                <Box>Previous Responses</Box>
                <ClearAllIconButton onClick={() => {}} color={"error"}/>
            </Box>
            <Box>
                {renderBody()}
            </Box>
        </Box>
    )


    function renderBody() {
        if (!result.data) {
            return <Box sx={SX.no}>Response is no history yet</Box>
        }

        return result.data.map((query, index) => (
            <QueryBodyHistoryItem key={index} query={query} />
        ))
    }
}
