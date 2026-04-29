import {Box, Link} from "@mui/material"

import {Connection as SshConnection, NodeTab, NodeTabType} from "../../../../api/node/type"
import {useRouterQueryDatabase, useRouterQuerySchemas} from "../../../../api/query/hook"
import {Connection as QueryConnection} from "../../../../api/query/type"
import {SxPropsMap} from "../../../../app/type"
import {useStore, useStoreAction} from "../../../../provider/StoreProvider"
import {AutocompleteFetch} from "../../../view/autocomplete/AutocompleteFetch"
import {Monitor} from "../../../widgets/monitor/Monitor"
import {ClusterNoPostgresVault, NoDatabaseError} from "../overview/OverviewError"
import {NodeMainQueries} from "./NodeMainQueries"
import {NodeMainTitle} from "./NodeMainTitle"

const SX: SxPropsMap = {
    main: {flexGrow: 1, overflow: "auto", display: "flex", flexDirection: "column", gap: 1},
    inputs: {display: "flex", alignItems: "center", gap: 1, width: "300px"},
}

const Tabs: {[key in NodeTabType]: NodeTab} = {
    [NodeTabType.MONITOR]: {
        label: "Monitor",
        body: (queryCon: QueryConnection, sshCon: SshConnection) => <Monitor queryCon={queryCon} sshCon={sshCon}/>,
        info: <>
            Here you can check some basic charts about your overall database and each database separately
            by specifying database name in the input near by.
            If you have some proposal what can be added here, please, suggest
            it <Link href={"https://github.com/veegres/ivory/issues"} target={"_blank"}>here</Link>
        </>
    },
    [NodeTabType.QUERY]: {
        label: "Queries",
        body: (queryCon: QueryConnection) => <NodeMainQueries connection={queryCon}/>,
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
    queryCon: QueryConnection,
    sshCon: SshConnection,
}

export function NodeMain(props: Props) {
    const {tab, queryCon: qc, sshCon: sc} = props
    const {dbName, dbSchema} = useStore(s => s.node)
    const {setDbName, setDbSchema} = useStoreAction

    const {label, info, body} = Tabs[tab]

    const db = {...qc.db, name: dbName, schema: dbSchema}
    const queryCon = {...qc, db}
    const sshCon = sc

    return (
        <Box sx={SX.main}>
            <NodeMainTitle label={label} info={info} db={db} renderActions={renderActions()}/>
            {renderBody()}
        </Box>
    )

    function renderBody() {
        if (!queryCon.vaultId) return <ClusterNoPostgresVault/>
        if (db.host === "-") return <NoDatabaseError/>
        return body(queryCon, sshCon)
    }

    function renderActions() {
        if (!queryCon.vaultId) return null
        return (
            <Box sx={SX.inputs}>
                <AutocompleteFetch
                    value={dbSchema || null}
                    connection={queryCon}
                    useFetch={useRouterQuerySchemas}
                    placeholder={"Schema"}
                    variant={"outlined"}
                    padding={"3px"}
                    onUpdate={(v) => setDbSchema(v || undefined)}
                    disabled={!dbName}
                />
                <AutocompleteFetch
                    value={dbName || null}
                    connection={queryCon}
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
