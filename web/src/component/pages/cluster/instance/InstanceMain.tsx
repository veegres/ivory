import {Box, Link} from "@mui/material";

import {InstanceTab, InstanceTabType} from "../../../../api/instance/type";
import {useRouterQueryDatabase} from "../../../../api/query/hook";
import {Database, QueryConnection} from "../../../../api/query/type";
import {SxPropsMap} from "../../../../app/type";
import {getQueryConnection} from "../../../../app/utils";
import {useStore, useStoreAction} from "../../../../provider/StoreProvider";
import {AutocompleteFetch} from "../../../view/autocomplete/AutocompleteFetch";
import {Chart} from "../../../widgets/chart/Chart";
import {ClusterNoPostgresPassword, NoDatabaseError} from "../overview/OverviewError";
import {InstanceMainQueries} from "./InstanceMainQueries";
import {InstanceMainTitle} from "./InstanceMainTitle";

const SX: SxPropsMap = {
    main: {flexGrow: 1, overflow: "auto", display: "flex", flexDirection: "column", gap: 1},
}

const Tabs: {[key in InstanceTabType]: InstanceTab} = {
    [InstanceTabType.CHART]: {
        label: "Charts",
        body: (connection: QueryConnection) => <Chart connection={connection}/>,
        info: <>
            Here you can check some basic charts about your overall database and each database separately
            by specifying database name in the input near by.
            If you have some proposal what can be added here, please, suggest
            it <Link href={"https://github.com/veegres/ivory/issues"} target={"_blank"}>here</Link>
        </>
    },
    [InstanceTabType.QUERY]: {
        label: "Queries",
        body: (connection: QueryConnection) => <InstanceMainQueries connection={connection}/>,
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
    const {dbName} = useStore(s => s.instance)
    const {setDbName} = useStoreAction

    if (!activeCluster) return null
    const {cluster} = activeCluster
    const {label, info, body} = Tabs[tab]

    const credentialId = cluster.credentials.postgresId
    const db = {...database, name: dbName}
    const connection = getQueryConnection(cluster, db)

    return (
        <Box sx={SX.main}>
            <InstanceMainTitle label={label} info={info} db={db} renderActions={renderActions()}/>
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
            <Box width={200}>
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
