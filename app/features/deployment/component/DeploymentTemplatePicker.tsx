import {Box} from "@mui/material"
import {useState} from "react"

import {Note} from "../../../shared/component/box/Note"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {KeeperPlugin, PlatformPlugin} from "../../node/api/NodeType"
import {useRouterDeploymentTemplateList} from "../api/DeploymentHook"
import {Template} from "../api/DeploymentType"
import {DeploymentTemplateForm} from "./DeploymentTemplateForm"
import {DeploymentTemplateList} from "./DeploymentTemplateList"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
}

// form holds the template being copied or edited; the absent template of a
// brand new one is exactly what tells the two apart
type FormState = {template?: Template, edit: boolean}

type Props = {
    keeper: KeeperPlugin,
    platform: PlatformPlugin,
    hint: string,
    onPick: (template: Template) => void,
}

// useDeploymentTemplatePicker is the "choose a template" step both deploy
// dialogs open on: the list, the three ways a template comes to exist, and the
// back navigation between them. It is a hook rather than a component because
// the dialog around it owns the back arrow, so it has to know whether the
// picker is showing a form.
export function useDeploymentTemplatePicker(props: Props) {
    const {keeper, platform, hint, onPick} = props
    const [form, setForm] = useState<FormState>()
    // NOTE: shares the list query's cache with DeploymentTemplateList, so the
    // name check costs no extra request
    const list = useRouterDeploymentTemplateList({keeper, platform})

    return {editing: !!form, render, back}

    function render() {
        return form ? renderForm(form) : renderList()
    }

    function renderList() {
        return (
            <Box sx={SX.box}>
                <Note center={true}>{hint}</Note>
                <DeploymentTemplateList
                    keeper={keeper}
                    platform={platform}
                    onOpen={onPick}
                    onCopy={(template) => setForm({template, edit: false})}
                    onEdit={(template) => setForm({template, edit: true})}
                    onNew={() => setForm({edit: false})}
                />
            </Box>
        )
    }

    function renderForm(current: FormState) {
        return (
            <DeploymentTemplateForm
                keeper={keeper}
                platform={platform}
                template={current.template}
                edit={current.edit}
                takenNames={list.data?.map(t => t.name) ?? []}
                onDone={close}
            />
        )
    }

    // NOTE: saving lands back on the list rather than on the deploy form - the
    // saved template is one row among the others, and picking it is what runs
    // it, so writing one never skips the step every other template takes
    function close() {
        setForm(undefined)
    }

    function back() {
        if (!form) return false
        close()
        return true
    }
}
