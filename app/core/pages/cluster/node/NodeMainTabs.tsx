import {Link} from "@mui/material"
import {ReactNode} from "react"

import {Cluster, Node} from "../../../../features/cluster/api/ClusterType"
import {NodeTabType} from "../../../../features/node/api/NodeType"
import {Container} from "../../../../features/node/component/container/Container"
import {Keeper} from "../../../../features/node/component/keeper/Keeper"
import {Platform} from "../../../../features/node/component/platform/Platform"
import {getPlatformConnection, getQueryConnection} from "../../../../shared/helper/HelperUtils"
import {NodeMainQueries} from "./NodeMainQueries"
import {NodeMainTools} from "./NodeMainTools"

interface NodeTab {
    label: string,
    body: (cluster: Cluster, node: Node) => ReactNode,
    info?: ReactNode,
    actions?: ReactNode,
}

export const NODE_TABS: { [key in NodeTabType]: NodeTab } = {
    [NodeTabType.CONTAINER]: {
        label: "Container",
        body: (c: Cluster, n: Node) => <Container connection={getPlatformConnection(c, n.config.host, n.config.sshPort)}/>,
        info: <>
            Here you can check some basic charts about your overall database and each database separately
            by specifying database name in the input near by.
            If you have some proposal what can be added here, please, suggest
            it <Link href={"https://github.com/veegres/ivory/issues"} target={"_blank"}>here</Link>
        </>
    },
    [NodeTabType.DATABASE]: {
        label: "Database",
        body: (c: Cluster, n: Node) => <NodeMainQueries connection={getQueryConnection(c, n.config.host, n.config.dbPort)}/>,
        info: <>
            Here you can run some queries to troubleshoot your postgres (<b>always use LIMIT in queries
            to reduce number of rows, it will help to render and execute query faster</b>). There are some default
            queries
            which are provided by the <i>system</i>. If manual queries are enabled, you can do such
            things as:
            <ul style={{margin: "0"}}>
                <li>create your own <i>custom</i> queries</li>
                <li>edit <i>system</i> or <i>custom</i> queries</li>
                <li>rollback these changes at anytime to default state (the first query that was saved)</li>
            </ul>
        </>
    },
    [NodeTabType.KEEPER]: {
        label: "Keeper",
        body: (c: Cluster, n: Node) => <Keeper node={n} cluster={c}/>,
        info: <>
            Here you can manipulate with your cluster management systems. Ivory calls it keeper.
            You can adjust your PostgreSQL configurations here, and any changes made will be applied to
            all cluster nodes. Instead of rewriting the entire configuration, it applies a patch
            update. If you wish to remove a specific setting, simply set it to <b>null</b>. Keep in mind that
            modifying certain parameters may necessitate restarting PostgreSQL. For further details on how
            this process functions, refer to
            the <Link href={"https://patroni.readthedocs.io/en/latest/rest_api.html#config-endpoint"}
                      target={"_blank"}>documentation</Link>.
        </>
    },
    [NodeTabType.TOOLS]: {
        label: "Tools",
        body: (c: Cluster, n: Node) => <NodeMainTools node={n} cluster={c}/>,
        info: <>
            Here, you can efficiently decrease the size of bloated tables and indexes without imposing
            heavy locks. This functionality is powered by
            the <Link href={"https://github.com/dataegret/pgcompacttable"} target={"_blank"}>pgcompacttable</Link> tool,
            which is seamlessly integrated with Ivory for streamlined usage. Ivory simplifies visualization
            and centralizes information about jobs and logs within each cluster, ensuring convenient access when
            needed. It's important to note that this tool can only be executed on the master node, and
            in the target database, the contrib module pgstattuple must be installed using the command
            "<b>CREATE EXTENSION IF NOT EXISTS pgstattuple;</b>". Ivory supports such features as:
            <ul>
                <li><b>Delay ratio</b> - A dynamic part of the delay between rounds is calculated as previous-round-time
                    * delay-ratio. By default 2.
                </li>
                <li><b>Min table size</b> - Tables smaller than the specified size (in megabytes) will be excluded from
                    processing.
                </li>
                <li><b>Max table size</b> - Tables larger than the specified size (in megabytes) will be excluded from
                    processing.
                </li>
                <li><b>Force</b> - Try to compact even those tables and indexes that do not meet minimal bloat
                    requirements.
                </li>
                <li><b>Routing vacuum</b> - Turn on the routine vacuum. By default all the vacuums are off.</li>
                <li><b>No initial vacuum</b> - Turn off initial vacuum before table processing.</li>
                <li><b>Initial reindex</b> - Perform an initial reindex of tables before processing.</li>
                <li><b>No reindex</b> - Turn off reindexing of tables after processing.</li>
            </ul>
        </>
    },
    [NodeTabType.PLATFORM]: {
        label: "Platform",
        body: (c: Cluster, n: Node) => <Platform connection={getPlatformConnection(c, n.config.host, n.config.sshPort)}/>,
        info: <>
            Here you can use see what is going on with your Platform, like Linux
        </>
    },
}
