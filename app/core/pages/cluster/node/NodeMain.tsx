import {Box} from "@mui/material"

import {Cluster, Node} from "../../../../features/cluster/type"
import {NodeTabType} from "../../../../features/node/type"
import {useRouterQueryDatabase, useRouterQuerySchemas} from "../../../../features/query/hook"
import {AutocompleteFetch} from "../../../../shared/component/autocomplete/AutocompleteFetch"
import {SxPropsMap} from "../../../../shared/helper/type"
import {getQueryConnection} from "../../../../shared/helper/utils"
import {useStore, useStoreAction} from "../../../../shared/provider/StoreProvider"
import {NODE_TABS} from "./NodeMainTabs"
import {NodeMainTitle} from "./NodeMainTitle"

const SX: SxPropsMap = {
    main: {flexGrow: 1, overflow: "auto", display: "flex", flexDirection: "column", gap: 1},
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
        if (tab !== NodeTabType.DATABASE) return
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
}
