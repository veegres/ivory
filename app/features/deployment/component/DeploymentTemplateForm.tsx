import {Box, Button} from "@mui/material"
import {createPortal} from "react-dom"

import {ErrorSmart} from "../../../shared/component/box/ErrorSmart"
import {WarningList} from "../../../shared/component/box/WarningList"
import {useDialogFooter} from "../../../shared/component/button/DialogButton"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {Feature} from "../../Feature"
import {ManageAccess} from "../../management/component/ManageAccess"
import {KeeperPlugin, PlatformPlugin} from "../../node/api/NodeType"
import {useRouterDeploymentTemplateCreate, useRouterDeploymentTemplateUpdate, useTemplateForm} from "../api/DeploymentHook"
import {Template, TemplateRequest} from "../api/DeploymentType"
import {DeploymentTemplateEditor} from "./DeploymentTemplateEditor"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 2},
}

type Props = {
    keeper: KeeperPlugin,
    platform: PlatformPlugin,
    // template being edited, or the one being copied; absent means a brand new
    // one written from scratch
    template?: Template,
    edit?: boolean,
    takenNames: string[],
    onDone: (template: Template) => void,
}

// DeploymentTemplateForm covers the three ways a template comes to exist: a new
// one written from scratch, a copy of another (yours or a shipped one), and an
// edit of your own. All three end the same way - the saved template is handed
// back so it can be run.
export function DeploymentTemplateForm(props: Props) {
    const {keeper, platform: initialPlatform, template, edit = false, takenNames, onDone} = props
    // NOTE: a template's platform is fixed once it exists, and a new one takes
    // the dialog's - there is nothing to choose while only one platform exists
    const platform = template?.platform ?? initialPlatform
    const footer = useDialogFooter()
    const form = useTemplateForm(getInitialTemplate())
    const create = useRouterDeploymentTemplateCreate(onDone)
    const update = useRouterDeploymentTemplateUpdate(onDone)

    const action = edit ? update : create

    return (
        <Box sx={SX.box}>
            {action.isError && <ErrorSmart error={action.error}/>}
            <WarningList warnings={getWarnings()}/>
            <DeploymentTemplateEditor
                template={form.template}
                onChange={form.setTemplate}
                onCommandChange={form.updateCommand}
                onCommandAdd={form.addCommand}
                onCommandRemove={form.removeCommand}
            />
            {renderAction()}
        </Box>
    )

    // NOTE: the button belongs in the dialog's action bar, next to where every
    // other dialog keeps its confirm - so it is rendered there rather than at
    // the end of a form the user has to scroll through to reach it
    function renderAction() {
        if (footer === undefined) return renderButton()
        if (footer === null) return
        return createPortal(renderButton(), footer)
    }

    // NOTE: error={true} on purpose - a user without the permission must see
    // why they cannot proceed, instead of a form with no button and no reason
    function renderButton() {
        return (
            <ManageAccess feature={getFeature()} error={true}>
                <Button
                    fullWidth={true}
                    loading={action.isPending}
                    disabled={!form.valid || isNameTaken()}
                    onClick={handleAction}
                >
                    {edit ? "Save" : "Create"}
                </Button>
            </ManageAccess>
        )
    }

    function getFeature() {
        return edit ? Feature.ManageDeploymentTemplateUpdate : Feature.ManageDeploymentTemplateCreate
    }

    // NOTE: the node is part of the warning - a collapsed command hides its
    // own, and a disabled Create button with no visible reason is worse than a
    // slightly longer sentence
    function getWarnings() {
        const warnings = form.unknown.flatMap((names, index) => (
            names.map(name => `Node ${index + 1}: ${name} is not a known variable`)
        ))
        if (isNameTaken()) warnings.push(`The name "${form.template.name.trim()}" is already taken`)
        return warnings
    }

    // NOTE: keeping its own name is not a collision with itself
    function isNameTaken() {
        const name = form.template.name.trim()
        if (edit && name === template?.name) return false
        return takenNames.includes(name)
    }

    function handleAction() {
        const request = {...form.template, platform}
        if (edit && template) update.mutate({id: template.id, template: request})
        else create.mutate(request)
    }

    function getInitialTemplate(): TemplateRequest {
        if (!template) {
            return {name: "", description: "", keeper, platform: initialPlatform, commands: [{command: ""}]}
        }
        return {
            name: edit ? template.name : `${template.name} (copy)`,
            description: template.description,
            keeper: template.keeper,
            platform: template.platform,
            commands: template.commands,
        }
    }
}
