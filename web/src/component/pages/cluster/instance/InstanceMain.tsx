import {InstanceMainTitle} from "./InstanceMainTitle";
import {ClusterNoPostgresPassword, NoDatabaseError} from "../overview/OverviewError";
import {Box, Link} from "@mui/material";
import {Database, SxPropsMap} from "../../../../type/common";
import {InstanceTab, InstanceTabType} from "../../../../type/instance";
import {Chart} from "../../../shared/chart/Chart";
import {InstanceMainQueries} from "./InstanceMainQueries";
import {QueryApi} from "../../../../app/api";
import {AutocompleteFetch} from "../../../view/autocomplete/AutocompleteFetch";
import {useStore, useStoreAction} from "../../../../provider/StoreProvider";
import {getDomain} from "../../../../app/utils";

const SX: SxPropsMap = {
    main: {flexGrow: 1, overflow: "auto", display: "flex", flexDirection: "column", gap: 1},
}

const Tabs: {[key in InstanceTabType]: InstanceTab} = {
    [InstanceTabType.CHART]: {
        label: "Charts",
        body: (id: string, db: Database) => <Chart credentialId={id} db={db}/>,
        info: <>
            Here you can check some basic charts about your overall database and each database separately
            by specifying database name in the input near by.
            If you have some proposal what can be added here, please, suggest
            it <Link href={"https://github.com/veegres/ivory/issues"} target={"_blank"}>here</Link>
        </>
    },
    [InstanceTabType.QUERY]: {
        label: "Queries",
        body: (id: string, db: Database) => <InstanceMainQueries credentialId={id} db={db}/>,
        info: <>
            Here you can run some queries to troubleshoot your postgres. There are some default queries
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
    const {activeCluster, instance: {dbName}} = useStore()
    const {setDbName} = useStoreAction()

    if (!activeCluster) return null

    const postgresId = activeCluster?.cluster.credentials.postgresId
    const {label, info, body} = Tabs[tab]

    const db = {...database, name: dbName}

    return (
        <Box sx={SX.main}>
            <InstanceMainTitle label={label} info={info} db={db} renderActions={renderActions()}/>
            {renderBody()}
        </Box>
    )

    function renderBody() {
        if (!postgresId) return <ClusterNoPostgresPassword/>
        if (database.host === "-") return <NoDatabaseError/>

        return body(postgresId, db)
    }

    function renderActions() {
        if (!postgresId) return null

        return (
            <Box width={200}>
                <AutocompleteFetch
                    value={dbName || null}
                    keys={["query", "databases", getDomain(database), database.name ?? "postgres", postgresId]}
                    onFetch={(v) => QueryApi.databases({credentialId: postgresId, db: database, name: v})}
                    placeholder={"Database"}
                    variant={"outlined"}
                    padding={"3px"}
                    onUpdate={(v) => setDbName(v || undefined)}
                />
            </Box>
        )
    }
}
