import {Link} from "@mui/material"
import {ReactNode} from "react"

import {Cluster, Node} from "../../../../features/cluster/api/ClusterType"
import {NodeTabType} from "../../../../features/node/api/NodeType"
import {Container} from "../../../../features/node/component/container/Container"
import {Keeper} from "../../../../features/node/component/keeper/Keeper"
import {Platform} from "../../../../features/node/component/platform/Platform"
import {getKeeperOneRequest, getPlatformConnection, getQueryConnection} from "../../../../shared/helper/HelperUtils"
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
            Everything about the container this node runs in. The <i>overview</i> shows live
            resource usage and recent logs, and lets you start, restart, stop or remove the container,
            as well as deploy a new keeper container to this host. The <i>list</i> shows all
            containers that currently exist on the host, so you can check what else is running there.
        </>
    },
    [NodeTabType.DATABASE]: {
        label: "Database",
        body: (c: Cluster, n: Node) => <NodeMainQueries connection={getQueryConnection(c, n.config.host, n.config.dbPort)}/>,
        info: <>
            This tab is your window into PostgreSQL itself. The <i>console</i> lets you run arbitrary SQL,
            the <i>charts</i> section shows basic charts for the whole instance and for a particular
            database, and the remaining sections group ready-made queries by concern: activity,
            statistics, bloat, replication and other (<b>always use LIMIT in your queries to reduce the
            number of rows, it will help to render and execute the query faster</b>). Default queries are
            provided by the <i>system</i>. If manual queries are enabled, you can also:
            <ul style={{margin: "0"}}>
                <li>create your own <i>custom</i> queries</li>
                <li>edit <i>system</i> or <i>custom</i> queries</li>
                <li>rollback these changes at anytime to default state (the first query that was saved)</li>
            </ul>
        </>
    },
    [NodeTabType.KEEPER]: {
        label: "Keeper",
        body: (c: Cluster, n: Node) => (
            <Keeper
                request={getKeeperOneRequest(c, n.config.host, n.config.keeperPort)}
                cluster={c.name}
                candidates={c.nodes.map(node => node.host)}
                role={n.keeper.role}
            />
        ),
        info: <>
            Control over your cluster management system (e.g. Patroni) — Ivory calls it keeper.
            The actions let you reload or restart the keeper, reinitialise a replica, perform a switchover
            or failover, and schedule such operations for a later time. Below the actions you can adjust
            the cluster configuration: any change is applied to all cluster nodes as a patch update —
            instead of rewriting the entire configuration only the settings you provide are changed, and
            setting one to <b>null</b> removes it. Keep in mind that modifying certain parameters may
            necessitate restarting PostgreSQL. For further details on how this process functions, refer to
            the <Link href={"https://patroni.readthedocs.io/en/latest/rest_api.html#config-endpoint"}
                      target={"_blank"}>documentation</Link>.
        </>
    },
    [NodeTabType.TOOLS]: {
        label: "Tools",
        body: (c: Cluster, n: Node) => <NodeMainTools node={n} cluster={c}/>,
        info: <>
            External tools integrated into Ivory. Currently the only one
            is <i>pg_compacttable</i>: it efficiently decreases the size of bloated tables and indexes
            without imposing heavy locks (if you have a proposal for another tool, please, suggest
            it <Link href={"https://github.com/veegres/ivory/issues"} target={"_blank"}>here</Link>).
            This functionality is powered by
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
            A health view of the host (virtual machine) this node runs on. Ivory connects to it
            through the platform credentials and shows general host information, live CPU, memory and
            network usage charts, the list of running processes and system logs. It helps to spot
            problems at the host level — before digging into the container or the database itself.
        </>
    },
}
