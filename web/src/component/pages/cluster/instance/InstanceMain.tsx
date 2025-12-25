import {Box, Link} from "@mui/material"

import {InstanceTab, InstanceTabType} from "../../../../api/instance/type"
import {ConnectionRequest, Database} from "../../../../api/postgres"
import {useRouterQueryDatabase, useRouterQuerySchemas} from "../../../../api/query/hook"
import {SxPropsMap} from "../../../../app/type"
import {getConnectionRequest} from "../../../../app/utils"
import {useStore, useStoreAction} from "../../../../provider/StoreProvider"
import {AutocompleteFetch} from "../../../view/autocomplete/AutocompleteFetch"
import {Chart} from "../../../widgets/chart/Chart"
import {ClusterNoPostgresPassword, NoDatabaseError} from "../overview/OverviewError"
import {InstanceMainQueries} from "./InstanceMainQueries"
import {InstanceMainTitle} from "./InstanceMainTitle"

const SX: SxPropsMap = {
    main: {flexGrow: 1, overflow: "auto", display: "flex", flexDirection: "column", gap: 1},
    inputs: {display: "flex", alignItems: "center", gap: 1, width: "300px"},
}

const Tabs: {[key in InstanceTabType]: InstanceTab} = {
    [InstanceTabType.CHART]: {
        label: "Charts",
        body: (connection: ConnectionRequest) => <Chart connection={connection}/>,
        info: <>
            Here you can check some basic charts about your overall database and each database separately
            by specifying database name in the input near by.
            If you have some proposal what can be added here, please, suggest
            it <Link href={"https://github.com/veegres/ivory/issues"} target={"_blank"}>here</Link>
        </>
    },
    [InstanceTabType.QUERY]: {
        label: "Queries",
        body: (connection: ConnectionRequest) => <InstanceMainQueries connection={connection}/>,
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
    tab: InstanceTabType,
    database: Database,
}

export function InstanceMain(props: Props) {
    const {tab, database} = props
    const activeCluster = useStore(s => s.activeCluster)
    const {dbName, dbSchema} = useStore(s => s.instance)
    const {setDbName, setDbSchema} = useStoreAction

    if (!activeCluster) return null
    const {cluster} = activeCluster
    const {label, info, body} = Tabs[tab]

    const credentialId = cluster.credentials.postgresId
    const db = {...database, name: dbName, schema: dbSchema} as Database
    const connection = getConnectionRequest(cluster, db)

    return (
        <Box sx={SX.main}>
            <InstanceMainTitle label={label} info={info} db={database} renderActions={renderActions()}/>
            {renderBody()}
        </Box>
    )

    function renderBody() {
        if (!credentialId) return <ClusterNoPostgresPassword/>
        if (database.host === "-") return <NoDatabaseError/>
        return body(connection)
    }

    function renderActions() {
        if (!credentialId) return null
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
