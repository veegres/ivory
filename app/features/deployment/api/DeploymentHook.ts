import {useQuery} from "@tanstack/react-query"
import {useMemo, useState} from "react"

import {getUnknownPlaceholders} from "../../../shared/helper/HelperUtils"
import {useMutationAdapter} from "../../../shared/hook/QueryCustom"
import {DeploymentApi} from "./DeploymentRouter"
import {Template, TemplateCommand, TemplateListRequest, TemplateRequest} from "./DeploymentType"

export function useRouterDeploymentTemplateList(request?: TemplateListRequest) {
    return useQuery({
        queryKey: DeploymentApi.template.list.key(request),
        queryFn: () => DeploymentApi.template.list.fn(request),
    })
}

export function useRouterDeploymentTemplateCreate(onSuccess?: (template: Template) => void) {
    return useMutationAdapter({
        mutationFn: DeploymentApi.template.create.fn,
        mutationKey: DeploymentApi.template.create.key(),
        successKeys: [DeploymentApi.template.list.keyCommon()],
        onSuccess: (_, data) => onSuccess ? onSuccess(data) : void 0,
    })
}

export function useRouterDeploymentTemplateUpdate(onSuccess?: (template: Template) => void) {
    return useMutationAdapter({
        mutationFn: DeploymentApi.template.update.fn,
        mutationKey: DeploymentApi.template.update.key(),
        successKeys: [DeploymentApi.template.list.keyCommon()],
        onSuccess: (_, data) => onSuccess ? onSuccess(data) : void 0,
    })
}

export function useRouterDeploymentTemplateDelete() {
    return useMutationAdapter({
        mutationFn: DeploymentApi.template.delete.fn,
        mutationKey: DeploymentApi.template.delete.key(),
        successKeys: [DeploymentApi.template.list.keyCommon()],
    })
}

// getUnknownCommandPlaceholders reports the placeholders one command uses that
// are outside the closed vocabulary - a validation error rather than a new
// variable. A command and its post script share one scope, so they are checked
// as one piece of text.
export function getUnknownCommandPlaceholders(command: TemplateCommand) {
    return getUnknownPlaceholders(command.command + " " + (command.postScript ?? ""))
}

// useTemplateForm owns the editable copy of a template: the list of commands
// plus its name and description. It is deliberately local state - a template
// is only written back when the user saves it.
export function useTemplateForm(initial: TemplateRequest) {
    const [template, setTemplate] = useState<TemplateRequest>(initial)

    // NOTE: kept per command rather than flattened - a template has several,
    // and "{{foo}} is not a known variable" is no help without saying where
    const unknown = useMemo(handleMemoUnknown, [template.commands])
    const valid = !!template.name.trim()
        && template.commands.length > 0
        && unknown.every(commandUnknown => commandUnknown.length === 0)

    return {template, setTemplate, unknown, valid, updateCommand, addCommand, removeCommand}

    function handleMemoUnknown() {
        return template.commands.map(getUnknownCommandPlaceholders)
    }

    function updateCommand(index: number, command: TemplateCommand) {
        setTemplate(prev => ({...prev, commands: prev.commands.map((c, i) => i === index ? command : c)}))
    }

    function addCommand(source?: TemplateCommand) {
        setTemplate(prev => ({...prev, commands: [...prev.commands, source ?? {command: ""}]}))
    }

    function removeCommand(index: number) {
        setTemplate(prev => ({...prev, commands: prev.commands.filter((_, i) => i !== index)}))
    }
}
