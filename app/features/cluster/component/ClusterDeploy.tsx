import {RocketLaunch} from "@mui/icons-material"
import {useEffect, useState} from "react"

import {DialogScreen} from "../../../shared/component/box/DialogScreen"
import {ErrorSmart} from "../../../shared/component/box/ErrorSmart"
import {DialogButton} from "../../../shared/component/button/DialogButton"
import {SkeletonGroup} from "../../../shared/component/progress/SkeletonGroup"
import {Template} from "../../deployment/api/DeploymentType"
import {DeploymentTemplateForm} from "../../deployment/component/DeploymentTemplateForm"
import {DeploymentTemplateList} from "../../deployment/component/DeploymentTemplateList"
import {Feature} from "../../Feature"
import {ManageAccess} from "../../management/component/ManageAccess"
import {useRouterNodeKeeperDeploySpec} from "../../node/api/NodeHook"
import {KeeperPlugin, PlatformPlugin} from "../../node/api/NodeType"
import {DbPlugin} from "../../query/api/QueryType"
import {ClusterDeployForm} from "./ClusterDeployForm"

// NOTE: a constant until a second platform exists - it was state with no
// setter, which reads as a choice nobody can make
const platform = PlatformPlugin.DOCKER

// NOTE: create covers both a blank template and a copy of another - they are
// the same POST, and the source is only what the editor opens with; update is
// the only one that needs a template to write back to
type Step =
    | {kind: "list"}
    | {kind: "create", source?: Template}
    | {kind: "update", template: Template}
    | {kind: "deploy", template: Template, logs?: string[]}

type Props = {
    keeper: KeeperPlugin,
    database: DbPlugin,
    withLabel?: boolean,
    size?: number,
}

// ClusterDeploy is the dialog every cluster deployment starts in: it owns which
// screen is up and the navigation between them, and nothing else. Each screen
// brings its own content and its own action bar.
export function ClusterDeploy(props: Props) {
    const {keeper, database, withLabel = false, size} = props
    const [step, setStep] = useState<Step>({kind: "list"})
    const spec = useRouterNodeKeeperDeploySpec(keeper)

    useEffect(handleEffectPluginProps, [keeper, database])

    return (
        <ManageAccess feature={Feature.ManageClusterCreate}>
            <DialogButton
                title={"DEPLOY CLUSTER"}
                icon={<RocketLaunch/>}
                variant={withLabel ? "button_label" : "button"}
                label={"Deploy"}
                size={size}
                back={step.kind !== "list"}
                onBackClick={handleBack}
                onClose={handleReset}
            >
                {renderStep()}
            </DialogButton>
        </ManageAccess>
    )

    function renderStep() {
        if (spec.isError) return <DialogScreen><ErrorSmart error={spec.error}/></DialogScreen>
        if (spec.isPending) return <DialogScreen><SkeletonGroup count={3}/></DialogScreen>
        switch (step.kind) {
            case "list":
                return (
                    <DeploymentTemplateList
                        keeper={keeper}
                        platform={platform}
                        hint={"Pick a template to deploy a cluster, copy one to adjust it, or write a new one"}
                        onOpen={(template) => setStep({kind: "deploy", template})}
                        onCopy={(source) => setStep({kind: "create", source})}
                        onEdit={(template) => setStep({kind: "update", template})}
                        onNew={() => setStep({kind: "create"})}
                    />
                )
            case "create":
                return (
                    <DeploymentTemplateForm
                        keeper={keeper}
                        platform={platform}
                        edit={false}
                        source={step.source}
                        onDone={handleTemplateDone}
                    />
                )
            case "update":
                return (
                    <DeploymentTemplateForm
                        keeper={keeper}
                        platform={platform}
                        edit={true}
                        template={step.template}
                        onDone={handleTemplateDone}
                    />
                )
            case "deploy":
                return (
                    <ClusterDeployForm
                        keeper={keeper}
                        database={database}
                        template={step.template}
                        spec={spec.data}
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

    // NOTE: the plugin selectors are disabled inside the dialog, so the plugins
    // can only change through the cluster list filter; the filter also changes
    // the available templates, hence the reset
    function handleEffectPluginProps() {
        setStep({kind: "list"})
    }
}
