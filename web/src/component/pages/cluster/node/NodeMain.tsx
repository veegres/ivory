import {Box, Link} from "@mui/material"

import {NodeTab, NodeTabType} from "../../../../api/node/type"
import {useRouterQueryDatabase, useRouterQuerySchemas} from "../../../../api/query/hook"
import {Database} from "../../../../api/query/postgres"
import {Connection} from "../../../../api/query/type"
import {SxPropsMap} from "../../../../app/type"
import {getConnection} from "../../../../app/utils"
import {useStore, useStoreAction} from "../../../../provider/StoreProvider"
import {AutocompleteFetch} from "../../../view/autocomplete/AutocompleteFetch"
import {Chart} from "../../../widgets/chart/Chart"
import {ClusterNoPostgresVault, NoDatabaseError} from "../overview/OverviewError"
import {NodeMainQueries} from "./NodeMainQueries"
import {NodeMainTitle} from "./NodeMainTitle"

const SX: SxPropsMap = {
    main: {flexGrow: 1, overflow: "auto", display: "flex", flexDirection: "column", gap: 1},
    inputs: {display: "flex", alignItems: "center", gap: 1, width: "300px"},
}

const Tabs: {[key in NodeTabType]: NodeTab} = {
    [NodeTabType.CHART]: {
        label: "Charts",
        body: (connection: Connection) => <Chart connection={connection}/>,
        info: <>
            Here you can check some basic charts about your overall database and each database separately
            by specifying database name in the input near by.
            If you have some proposal what can be added here, please, suggest
            it <Link href={"https://github.com/veegres/ivory/issues"} target={"_blank"}>here</Link>
        </>
    },
    [NodeTabType.QUERY]: {
        label: "Queries",
        body: (connection: Connection) => <NodeMainQueries connection={connection}/>,
        info: <>
            Here you can run some queries to troubleshoot your postgres (<b>always use LIMIT in queries
            to reduce number of rows, it will help to render and execute query faster</b>). There are some default queries
            which are provided by the <i>system</i>. If manual queries are enabled, you can do such
            things as:
            <ul style={{margin: "0"}}>
                <li>create your own <i>custom</i> queries</li>
                <li>edit <i>system</i> or <i>custom</i> queries</li>
                <li>rollback these changes at anytime to default state (the first query that was saved)</li>
            </ul>
        </>
    }
}

type Props = {
    tab: NodeTabType,
    database: Database,
}

export function NodeMain(props: Props) {
    const {tab, database} = props
    const activeCluster = useStore(s => s.activeCluster)
    const {dbName, dbSchema} = useStore(s => s.node)
    const {setDbName, setDbSchema} = useStoreAction

    if (!activeCluster) return null
    const {cluster} = activeCluster
    const {label, info, body} = Tabs[tab]

    const vaultId = cluster.vaults.databaseId
    const db = {...database, name: dbName, schema: dbSchema} as Database
    const connection = getConnection(cluster, db)

    return (
        <Box sx={SX.main}>
            <NodeMainTitle label={label} info={info} db={database} renderActions={renderActions()}/>
            {renderBody()}
        </Box>
    )

    function renderBody() {
        if (!vaultId) return <ClusterNoPostgresVault/>
        if (database.host === "-") return <NoDatabaseError/>
        return body(connection)
    }

    function renderActions() {
        if (!vaultId) return null
        return (
            <Box sx={SX.inputs}>
                <AutocompleteFetch
                    value={dbSchema || null}
                    connection={connection}
                    useFetch={useRouterQuerySchemas}
                    placeholder={"Schema"}
                    variant={"outlined"}
                    padding={"3px"}
                    onUpdate={(v) => setDbSchema(v || undefined)}
                    disabled={!dbName}
                />
                <AutocompleteFetch
                    value={dbName || null}
                    connection={connection}
                    useFetch={useRouterQueryDatabase}
                    placeholder={"Database"}
                    variant={"outlined"}
                    padding={"3px"}
                    onUpdate={(v) => setDbName(v || undefined)}
                />
            </Box>
        )
    }
}
