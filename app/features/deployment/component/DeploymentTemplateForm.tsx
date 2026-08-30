import {Box, Button} from "@mui/material"

import {DialogScreen} from "../../../shared/component/box/DialogScreen"
import {ErrorSmart} from "../../../shared/component/box/ErrorSmart"
import {WarningList} from "../../../shared/component/box/WarningList"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {Feature} from "../../Feature"
import {ManageAccess} from "../../management/component/ManageAccess"
import {KeeperPlugin, PlatformPlugin} from "../../node/api/NodeType"
import {
    useRouterDeploymentTemplateCreate,
    useRouterDeploymentTemplateList,
    useRouterDeploymentTemplateUpdate,
    useTemplateForm,
} from "../api/DeploymentHook"
import {Template, TemplateRequest} from "../api/DeploymentType"
import {DeploymentTemplateEditor} from "./DeploymentTemplateEditor"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 2},
}

// NOTE: source and template are not the same field under two names - a source
// is copied from and then forgotten, a template is written back to and is the
// only one that needs an id, which is why edit discriminates them
type Props = {
    keeper: KeeperPlugin,
    platform: PlatformPlugin,
    onDone: (template: Template) => void,
} & ({edit: false, source?: Template} | {edit: true, template: Template})

// DeploymentTemplateForm covers the three ways a template comes to exist: a new
// one written from scratch, a copy of another (yours or a shipped one), and an
// edit of your own. All three end the same way - the saved template is handed
// back so it can be run.
export function DeploymentTemplateForm(props: Props) {
    const {keeper, platform, edit, onDone} = props
    const form = useTemplateForm(getInitialTemplate())
    // NOTE: shares the list query's cache, so the name check costs no request
    const list = useRouterDeploymentTemplateList({keeper, platform})
    const create = useRouterDeploymentTemplateCreate(onDone)
    const update = useRouterDeploymentTemplateUpdate(onDone)

    const action = edit ? update : create

    return (
        <DialogScreen renderActions={renderActions()}>
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
            </Box>
        </DialogScreen>
    )

    // NOTE: error={true} on purpose - a user without the permission must see
    // why they cannot proceed, instead of a form with no button and no reason
    function renderActions() {
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
        if (props.edit && name === props.template.name) return false
        return (list.data ?? []).some(template => template.name === name)
    }

    function handleAction() {
        if (props.edit) update.mutate({id: props.template.id, template: form.template})
        else create.mutate(form.template)
    }

    // NOTE: a template's platform is fixed once it exists, and a new one takes
    // the dialog's - there is nothing to choose while only one platform exists
    function getInitialTemplate(): TemplateRequest {
        const origin = props.edit ? props.template : props.source
        if (!origin) {
            return {name: "", description: "", keeper, platform, commands: [{command: ""}]}
        }
        return {
            name: props.edit ? origin.name : `${origin.name} (copy)`,
            description: origin.description,
            keeper: origin.keeper,
            platform: origin.platform,
            defaults: origin.defaults,
            commands: origin.commands,
        }
    }
}
