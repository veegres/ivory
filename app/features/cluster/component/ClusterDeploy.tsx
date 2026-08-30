import {RocketLaunch} from "@mui/icons-material"

import {DeployScreenProps} from "../../deployment/api/DeploymentType"
import {DeploymentTemplateDialog} from "../../deployment/component/DeploymentTemplateDialog"
import {Feature} from "../../Feature"
import {ManageAccess} from "../../management/component/ManageAccess"
import {KeeperPlugin, PlatformPlugin} from "../../node/api/NodeType"
import {DbPlugin} from "../../query/api/QueryType"
import {ClusterDeployForm} from "./ClusterDeployForm"

// NOTE: a constant until a second platform exists - it was state with no
// setter, which reads as a choice nobody can make
const platform = PlatformPlugin.DOCKER

type Props = {
    keeper: KeeperPlugin,
    database: DbPlugin,
    withLabel?: boolean,
    size?: number,
}

// ClusterDeploy is the deploy dialog with a whole cluster as its last screen.
// Picking, writing and editing a template is the same everywhere, so all it
// adds to DeploymentTemplateDialog is the screen that fills one in.
export function ClusterDeploy(props: Props) {
    const {keeper, database, withLabel = false, size} = props

    return (
        <ManageAccess feature={Feature.ManageClusterCreate}>
            {/* NOTE: the plugin selectors are disabled inside the dialog, so
                the plugins can only change through the cluster list filter -
                which changes the templates on offer, hence the remount */}
            <DeploymentTemplateDialog
                key={`${keeper}:${database}`}
                keeper={keeper}
                platform={platform}
                title={"DEPLOY CLUSTER"}
                icon={<RocketLaunch/>}
                variant={withLabel ? "button_label" : "button"}
                label={"Deploy"}
                size={size}
                hint={"Pick a template to deploy a cluster, copy one to adjust it, or write a new one"}
                renderDeploy={renderDeploy}
            />
        </ManageAccess>
    )

    function renderDeploy(screen: DeployScreenProps) {
        return (
            <ClusterDeployForm
                keeper={keeper}
                database={database}
                template={screen.template}
                logs={screen.logs}
                onDeployed={screen.onDeployed}
            />
        )
    }
}
