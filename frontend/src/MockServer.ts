import {createServer} from "miragejs";

export function MockServer({environment = "test"}) {
    if (environment !== 'development') return;

    createServer({
        environment,
        routes() {
            this.namespace = "api"
            this.get("/cluster/info", () => ({
                members: [
                    {
                        name: "P4-FR-PPL-1",
                        timeline: 47,
                        lag: 59999,
                        state: "running",
                        host: "p4-fr-ppl-1",
                        role: "replica",
                        port: 5432,
                        api_url: "http://P4-FR-PPL-1:8008/patroni"
                    },
                    {
                        name: "P4-FR-PPL-2",
                        timeline: 47,
                        lag: 111,
                        state: "running",
                        host: "p4-fr-ppl-2",
                        role: "replica",
                        port: 5432,
                        api_url: "http://P4-FR-PPL-2:8008/patroni"
                    },
                    {
                        name: "P4-FR4-PPL-1",
                        timeline: 47,
                        state: "running",
                        host: "p4-fr4-ppl-1",
                        role: "leader",
                        port: 5432,
                        api_url: "http://P4-FR4-PPL-1:8008/patroni"
                    },
                    {
                        name: "P4-FR4-PPL-2",
                        timeline: 47,
                        lag: 0,
                        state: "running",
                        host: "p4-fr4-ppl-2",
                        role: "replica",
                        port: 5432,
                        api_url: "http://P4-FR4-PPL-2:8008/patroni"
                    }
                ]
            }))
            this.get("/cluster/list", () => ([
                {
                    name: "alpha",
                    nodes: ["p4-fr-ppl-1", "p4-fr-ppl-2", "p4-fr4-ppl-1", "p4-fr4-ppl-2"],
                },
                {
                    name: "betta",
                    nodes: ["p4-ava-ppl-1", "p4-ava-ppl-2", "p4-ava-ppl-3"],
                }
            ]))
            this.get("/node/:name/patroni", () => ({
                database_system_identifier: "6470821328108916785",
                postmaster_start_time: "2021-08-23 13:06:15.549 UTC",
                timeline: 47,
                cluster_unlocked: false,
                patroni: {scope: "alpha", version: "2.0.1"},
                state: "running",
                role: "master",
                xlog: {
                    received_location: 609432235694872,
                    replayed_timestamp: "2021-09-05 11:35:15.833 UTC",
                    paused: false,
                    replayed_location: 609432235694872
                },
                server_version: 90621
            }))
            this.get("/node/:name/config", () => ({
                retry_timeout: 60,
                postgresql: {
                    use_slots: false, remove_data_directory_on_diverged_timelines: true, parameters: {
                        autovacuum_vacuum_scale_factor: "0.01",
                        log_disconnections: false,
                        checkpoint_timeout: "30min",
                        autovacuum_vacuum_threshold: 50,
                        log_lock_waits: "true",
                        log_rotation_size: "0",
                        max_connections: "700",
                        hot_standby: true,
                        wal_log_hints: true,
                        log_line_prefix: "%m (user)%u (database)%d (PID)%p (SQLSTATE)%e started at %s ",
                        wal_keep_segments: "6250",
                        dynamic_shared_memory_type: "posix",
                        checkpoint_completion_target: "0.9",
                        max_wal_senders: "10",
                        work_mem: "64MB",
                        wal_buffers: "16MB",
                        max_wal_size: "50GB",
                        log_destination: "csvlog",
                        logging_collector: true,
                        wal_sync_method: "fdatasync",
                        log_directory: "/var/log/people-patroni-db",
                        shared_preload_libraries: "pg_stat_statements",
                        log_filename: "postgresql.log",
                        wal_level: "replica",
                        hot_standby_feedback: true,
                        constraint_exclusion: "partition",
                        autovacuum_vacuum_cost_limit: -1,
                        log_checkpoints: "true",
                        autovacuum_naptime: 60,
                        default_statistics_target: "100",
                        autovacuum_max_workers: 5,
                        log_statement: "ddl",
                        log_replication_commands: "true",
                        maintenance_work_mem: "512MB",
                        log_autovacuum_min_duration: "250",
                        shared_buffers: "8GB",
                        random_page_cost: "4",
                        effective_cache_size: "16GB",
                        autovacuum_vacuum_cost_delay: "2ms",
                        vacuum_cost_limit: "1000"
                    },
                    use_pg_rewind: true
                },
                maximum_lag_on_failover: 17592186044416,
                loop_wait: 10,
                master_start_timeout: 0,
                ttl: 90
            }))
        },
    })
}
