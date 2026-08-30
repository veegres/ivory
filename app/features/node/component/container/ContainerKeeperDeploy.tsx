import {Rocket} from "@mui/icons-material"

import {DeployScreenProps} from "../../../deployment/api/DeploymentType"
import {DeploymentTemplateDialog} from "../../../deployment/component/DeploymentTemplateDialog"
import {KeeperPlugin, PlatformPlugin, PlatformVaultConnection} from "../../api/NodeType"
import {ContainerKeeperDeployForm} from "./ContainerKeeperDeployForm"

type Props = {
    connection: PlatformVaultConnection,
    plugin: KeeperPlugin,
    cluster: string,
    node: string,
    keeperId?: string,
    databaseId?: string,
    sshKeyId?: string,
}

// ContainerKeeperDeploy deploys a keeper onto a single existing node: it picks
// one command out of a template and runs it here, calling node's own
// KeeperDeploy directly - no cluster endpoint involved, which is why it lives
// in the node/container feature rather than cluster.
export function ContainerKeeperDeploy(props: Props) {
    const {connection, plugin, cluster, node, keeperId, databaseId, sshKeyId} = props

    return (
        <DeploymentTemplateDialog
            keeper={plugin}
            platform={connection.platform ?? PlatformPlugin.DOCKER}
            title={"DEPLOY CONTAINER"}
            icon={<Rocket fontSize={"small"}/>}
            hint={"Pick a template to deploy this node - you choose which of its nodes runs here"}
            renderDeploy={renderDeploy}
        />
    )

    function renderDeploy(screen: DeployScreenProps) {
        return (
            <ContainerKeeperDeployForm
                connection={connection}
                plugin={plugin}
                cluster={cluster}
                node={node}
                template={screen.template}
                keeperId={keeperId}
                databaseId={databaseId}
                sshKeyId={sshKeyId}
                logs={screen.logs}
                onDeployed={screen.onDeployed}
            />
        )
    }
}
