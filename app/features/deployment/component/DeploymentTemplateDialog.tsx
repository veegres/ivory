import {SvgIconProps} from "@mui/material"
import {ReactElement, ReactNode, useState} from "react"

import {DialogButton} from "../../../shared/component/button/DialogButton"
import {KeeperPlugin, PlatformPlugin} from "../../node/api/NodeType"
import {DeployScreenProps, Template} from "../api/DeploymentType"
import {DeploymentTemplateForm} from "./DeploymentTemplateForm"
import {DeploymentTemplateList} from "./DeploymentTemplateList"

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
    platform: PlatformPlugin,
    title: string,
    icon: ReactElement<SvgIconProps>,
    hint: string,
    label?: string,
    variant?: "button" | "icon" | "button_label",
    size?: number,
    renderDeploy: (screen: DeployScreenProps) => ReactNode,
}

// DeploymentTemplateDialog is the dialog every deployment starts in: it owns
// which screen is up and the navigation between them, and nothing else. Three
// of its four screens are the same wherever it is opened - pick a template,
// write one, edit one - so the caller supplies only the fourth, the screen that
// runs the template it picked.
export function DeploymentTemplateDialog(props: Props) {
    const {keeper, platform, title, icon, hint, label, variant = "button", size, renderDeploy} = props
    const [step, setStep] = useState<Step>({kind: "list"})

    return (
        <DialogButton
            title={title}
            icon={icon}
            variant={variant}
            label={label}
            size={size}
            back={step.kind !== "list"}
            onBackClick={handleBack}
            onClose={handleReset}
        >
            {renderStep()}
        </DialogButton>
    )

    function renderStep() {
        switch (step.kind) {
            case "list":
                return (
                    <DeploymentTemplateList
                        keeper={keeper}
                        platform={platform}
                        hint={hint}
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
                return renderDeploy({
                    template: step.template,
                    logs: step.logs,
                    onDeployed: (logs) => setStep({kind: "deploy", template: step.template, logs}),
                })
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
