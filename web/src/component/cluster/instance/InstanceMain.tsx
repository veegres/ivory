import {InstanceMainTitle} from "./InstanceMainTitle";
import {ClusterNoPostgresPassword} from "../overview/OverviewError";
import {Box, Link} from "@mui/material";
import {Database, SxPropsMap} from "../../../type/common";
import {InstanceTab, InstanceTabType} from "../../../type/Instance";
import {Chart} from "../../shared/chart/Chart";
import {InstanceMainQueries} from "./InstanceMainQueries";
import {useState} from "react";
import {queryApi} from "../../../app/api";
import {AutocompleteFetch} from "../../view/././autocomplete/AutocompleteFetch";
import {useStore} from "../../../provider/StoreProvider";
import {getDomain} from "../../../app/utils";

const SX: SxPropsMap = {
    main: {flexGrow: 1, overflow: "auto", display: "flex", flexDirection: "column", gap: 1},
}

const Tabs: {[key in InstanceTabType]: InstanceTab} = {
    [InstanceTabType.CHART]: {
        label: "Charts",
        body: (cluster: string, db: Database) => <Chart cluster={cluster} db={db}/>,
        info: <>
            Here you can check some basic charts about your overall database and each database separatly
            by specifing database name in the input near by.
            If you have some proposal what can be added here, please, suggest
            it <Link href={"https://github.com/veegres/ivory/issues"} target={"_blank"}>here</Link>
        </>
    },
    [InstanceTabType.QUERY]: {
        label: "Queries",
        body: (cluster: string, db: Database) => <InstanceMainQueries cluster={cluster} db={db}/>,
        info: <>
            Here you can run some queries to troubleshoot your postgres. There are some default queries
            which are provided by the <i>system</i>. You can do such thing here as:
            <ul>
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
    const {store: { activeCluster }} = useStore()
    const [dbName, setDbName] = useState("")

    if (!activeCluster) return null

    const cluster = activeCluster?.cluster.name
    const password = activeCluster?.cluster.credentials.postgresId
    const {label, info, body} = Tabs[tab]

    return (
        <Box sx={SX.main}>
            <InstanceMainTitle label={label} info={info} renderActions={renderActions()}/>
            {password ? body(cluster, {...database, database: dbName}) : <ClusterNoPostgresPassword/>}
        </Box>
    )

    function renderActions() {
        return (
            <Box width={200}>
                <AutocompleteFetch
                    keys={["query", "databases", cluster, getDomain(database), database.database ?? ""]}
                    onFetch={(v) => queryApi.databases({clusterName: cluster, db: database, name: v})}
                    placeholder={" Database"}
                    variant={"standard"}
                    onUpdate={(v) => setDbName(v || "")}
                />
            </Box>
        )
    }
}
