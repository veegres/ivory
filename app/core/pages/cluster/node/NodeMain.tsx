import {Box, Link} from "@mui/material"

import {NodeConfig, Options} from "../../../../features/cluster/type"
import {NodeTab, NodeTabType} from "../../../../features/node/type"
import {useRouterQueryDatabase, useRouterQuerySchemas} from "../../../../features/query/hook"
import {AutocompleteFetch} from "../../../../shared/component/autocomplete/AutocompleteFetch"
import {SxPropsMap} from "../../../../shared/helper/type"
import {getPlatformConnection, getQueryConnection} from "../../../../shared/helper/utils"
import {useStore, useStoreAction} from "../../../../shared/provider/StoreProvider"
import {Monitor} from "../../../widgets/monitor/Monitor"
import {NodeMainQueries} from "./NodeMainQueries"
import {NodeMainTitle} from "./NodeMainTitle"

const SX: SxPropsMap = {
    main: {flexGrow: 1, overflow: "auto", display: "flex", flexDirection: "column", gap: 1},
    inputs: {display: "flex", alignItems: "center", gap: 1, width: "300px"},
}

const Tabs: {[key in NodeTabType]: NodeTab} = {
    [NodeTabType.MONITOR]: {
        body: (o: Options, c: NodeConfig) => <Monitor connection={getPlatformConnection(o, c.host, c.sshPort)}/>,
        info: <>
            Here you can check some basic charts about your overall database and each database separately
            by specifying database name in the input near by.
            If you have some proposal what can be added here, please, suggest
            it <Link href={"https://github.com/veegres/ivory/issues"} target={"_blank"}>here</Link>
        </>
    },
    [NodeTabType.QUERY]: {
        body: (o: Options, c: NodeConfig) => <NodeMainQueries connection={getQueryConnection(o, c.host, c.dbPort)}/>,
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
    options: Options,
    config: NodeConfig,
}

export function NodeMain(props: Props) {
    const {options, config} = props
    const node = useStore(s => s.nodeState)
    const {dbName, dbSchema} = useStore(s => s.nodeState)
    const {setDbName, setDbSchema} = useStoreAction

    const tab = node.nodeTab
    const {info, body} = Tabs[tab]

    return (
        <Box sx={SX.main}>
            <NodeMainTitle info={info} tab={tab} renderActions={renderActions()}/>
            {body(options, config)}
        </Box>
    )

    function renderActions() {
        if (tab !== NodeTabType.QUERY) return
        const con = getQueryConnection(options, config.host, config.dbPort)
        if (!con) return
        return (
            <Box sx={SX.inputs}>
                <AutocompleteFetch
                    value={dbSchema || null}
                    connection={con}
                    useFetch={useRouterQuerySchemas}
                    placeholder={"Schema"}
                    variant={"outlined"}
                    padding={"3px"}
                    onUpdate={(v) => setDbSchema(v || undefined)}
                    disabled={!dbName}
                />
                <AutocompleteFetch
                    value={dbName || null}
                    connection={con}
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
