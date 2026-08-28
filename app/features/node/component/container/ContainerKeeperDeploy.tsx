import {Rocket} from "@mui/icons-material"
import {useState} from "react"

import {DialogScreen} from "../../../../shared/component/box/DialogScreen"
import {ErrorSmart} from "../../../../shared/component/box/ErrorSmart"
import {DialogButton} from "../../../../shared/component/button/DialogButton"
import {SkeletonGroup} from "../../../../shared/component/progress/SkeletonGroup"
import {Template} from "../../../deployment/api/DeploymentType"
import {DeploymentTemplateForm} from "../../../deployment/component/DeploymentTemplateForm"
import {DeploymentTemplateList} from "../../../deployment/component/DeploymentTemplateList"
import {useRouterNodeKeeperDeploySpec} from "../../api/NodeHook"
import {KeeperPlugin, PlatformPlugin, PlatformVaultConnection} from "../../api/NodeType"
import {ContainerKeeperDeployForm} from "./ContainerKeeperDeployForm"

// NOTE: create covers both a blank template and a copy of another - they are
// the same POST, and the source is only what the editor opens with; update is
// the only one that needs a template to write back to
type Step =
    | {kind: "list"}
    | {kind: "create", source?: Template}
    | {kind: "update", template: Template}
    | {kind: "deploy", template: Template, logs?: string[]}

type Props = {
    connection: PlatformVaultConnection,
    plugin: KeeperPlugin,
    cluster: string,
    node: string,
    databaseId?: string,
    sshKeyId?: string,
}

// ContainerKeeperDeploy deploys a keeper onto a single existing node: it picks
// one command out of a template and runs it here, calling node's own
// KeeperDeploy directly - no cluster endpoint involved, which is why it lives
// in the node/container feature rather than cluster.
export function ContainerKeeperDeploy(props: Props) {
    const {connection, plugin, cluster, node, databaseId, sshKeyId} = props
    const [step, setStep] = useState<Step>({kind: "list"})
    const spec = useRouterNodeKeeperDeploySpec(plugin)

    const platform = connection.platform ?? PlatformPlugin.DOCKER

    return (
        <DialogButton
            title={"DEPLOY CONTAINER"}
            variant={"button"}
            icon={<Rocket fontSize={"small"}/>}
            back={step.kind !== "list"}
            onBackClick={handleBack}
            onClose={handleReset}
        >
            {renderStep()}
        </DialogButton>
    )

    function renderStep() {
        if (spec.isError) return <DialogScreen><ErrorSmart error={spec.error}/></DialogScreen>
        if (spec.isPending) return <DialogScreen><SkeletonGroup count={3}/></DialogScreen>
        switch (step.kind) {
            case "list":
                return (
                    <DeploymentTemplateList
                        keeper={plugin}
                        platform={platform}
                        hint={"Pick a template to deploy this node - you choose which of its nodes runs here"}
                        onOpen={(template) => setStep({kind: "deploy", template})}
                        onCopy={(source) => setStep({kind: "create", source})}
                        onEdit={(template) => setStep({kind: "update", template})}
                        onNew={() => setStep({kind: "create"})}
                    />
                )
            case "create":
                return (
                    <DeploymentTemplateForm
                        keeper={plugin}
                        platform={platform}
                        edit={false}
                        source={step.source}
                        onDone={handleTemplateDone}
                    />
                )
            case "update":
                return (
                    <DeploymentTemplateForm
                        keeper={plugin}
                        platform={platform}
                        edit={true}
                        template={step.template}
                        onDone={handleTemplateDone}
                    />
                )
            case "deploy":
                return (
                    <ContainerKeeperDeployForm
                        connection={connection}
                        plugin={plugin}
                        cluster={cluster}
                        node={node}
                        template={step.template}
                        spec={spec.data}
                        databaseId={databaseId}
                        sshKeyId={sshKeyId}
                        logs={step.logs}
                        onDeployed={(logs) => setStep({kind: "deploy", template: step.template, logs})}
                    />
                )
        }
    }

    // NOTE: saving lands back on the list rather than on the deploy form - the
    // saved template is one row among the others, and picking it is what runs
    // it, so writing one never skips the step every other template takes
    function handleTemplateDone() {
        setStep({kind: "list"})
    }

    function handleReset() {
        setStep({kind: "list"})
    }

    function handleBack() {
        if (step.kind === "deploy" && step.logs) return setStep({kind: "deploy", template: step.template})
        setStep({kind: "list"})
    }
}
