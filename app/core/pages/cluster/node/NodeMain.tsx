import {Box} from "@mui/material"

import {Cluster, Node} from "../../../../features/cluster/api/ClusterType"
import {NodeApi} from "../../../../features/node/api/NodeRouter"
import {NodeTabType} from "../../../../features/node/api/NodeType"
import {useRouterQueryDatabase, useRouterQuerySchemas} from "../../../../features/query/api/QueryHook"
import {AutocompleteFetch} from "../../../../shared/component/autocomplete/AutocompleteFetch"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {getPlatformConnection, getQueryConnection} from "../../../../shared/helper/HelperUtils"
import {useStore, useStoreAction} from "../../../../shared/provider/StoreProvider"
import {Refresher} from "../../../widgets/browser/Refresher"
import {NODE_TABS} from "./NodeMainTabs"
import {NodeMainTitle} from "./NodeMainTitle"

const SX: SxPropsMap = {
    main: {
        flexGrow: 1, overflow: "auto", display: "flex", flexDirection: "column",
        gap: 1, backgroundImage: "inherit", backgroundColor: "inherit",
    },
    inputs: {display: "flex", alignItems: "center", gap: 1, width: "300px"},
}

type Props = {
    cluster: Cluster,
    node: Node,
}

export function NodeMain(props: Props) {
    const {cluster, node} = props
    const nodeState = useStore(s => s.nodeState)
    const {dbName, dbSchema} = useStore(s => s.nodeState)
    const {setDbName, setDbSchema} = useStoreAction

    const tab = nodeState.nodeTab
    const {info, body} = NODE_TABS[tab]

    return (
        <Box sx={SX.main}>
            <NodeMainTitle info={info} tab={tab} renderActions={renderActions()}/>
            {body(cluster, node)}
        </Box>
    )

    function renderActions() {
        if (tab === NodeTabType.DATABASE) return renderDatabaseActions()
        if (tab === NodeTabType.PLATFORM) return renderPlatformActions()
        if (tab === NodeTabType.CONTAINER) return renderContainerActions()
        return
    }

    function renderDatabaseActions() {
        const con = getQueryConnection(cluster, node.config.host, node.config.dbPort)
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

    function renderPlatformActions() {
        const con = getPlatformConnection(cluster, node.config.host, node.config.sshPort)
        if (!con) return
        const queryKeys = [NodeApi.metrics.key(con.host), NodeApi.processes.key(con.host)]
        return <Refresher queryKeys={queryKeys} defaultPeriod={["3s", 3000]} size={32}/>
    }

    function renderContainerActions() {
        const con = getPlatformConnection(cluster, node.config.host, node.config.sshPort)
        if (!con) return
        const queryKeys = [NodeApi.deployment.metrics.key({connection: con, name: con.host})]
        return <Refresher queryKeys={queryKeys} defaultPeriod={["3s", 3000]} size={32}/>
    }
}
