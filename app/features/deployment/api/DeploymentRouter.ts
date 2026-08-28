import {api} from "../../Api"
import {R} from "../../management/api/ManagementType"
import {Template, TemplateListRequest, TemplateRequest} from "./DeploymentType"

export const DeploymentApi = {
    template: {
        // NOTE: one list returns both the user's own templates and the shipped
        // defaults, so the UI never has to reconcile two sources
        list: {
            key: (request?: TemplateListRequest) => ["deployment", "template", "list", request],
            // NOTE: refetch matches on a common prefix, so a write has to
            // invalidate with this and not key() - key() appends the filter
            // (or undefined) and would never match a filtered list
            keyCommon: () => ["deployment", "template", "list"],
            fn: (request?: TemplateListRequest) => api.get<R<Template[]>>("/deployment/template", {params: request})
                .then((response) => response.data.response),
        },
        create: {
            key: () => ["deployment", "template", "create"],
            fn: (request: TemplateRequest) => api.post<R<Template>>("/deployment/template", request)
                .then((response) => response.data.response),
        },
        update: {
            key: () => ["deployment", "template", "update"],
            fn: (request: {id: string, template: TemplateRequest}) => api.put<R<Template>>(`/deployment/template/${request.id}`, request.template)
                .then((response) => response.data.response),
        },
        delete: {
            key: () => ["deployment", "template", "delete"],
            fn: (id: string) => api.delete<R<string>>(`/deployment/template/${id}`)
                .then((response) => response.data.response),
        },
    },
}
