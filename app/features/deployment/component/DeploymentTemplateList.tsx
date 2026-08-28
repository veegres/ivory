import {Add} from "@mui/icons-material"
import {Box, Button} from "@mui/material"
import {SyntheticEvent} from "react"

import {ErrorSmart} from "../../../shared/component/box/ErrorSmart"
import {InfoColorBox} from "../../../shared/component/box/InfoColorBox"
import {InfoColorBoxRow} from "../../../shared/component/box/InfoColorBoxRow"
import {Note} from "../../../shared/component/box/Note"
import {CopyIconButton, DeleteIconButton, EditIconButton} from "../../../shared/component/button/IconButtons"
import {SkeletonGroup} from "../../../shared/component/progress/SkeletonGroup"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {PlatformPluginOptions} from "../../../shared/helper/HelperUtils"
import {Feature} from "../../Feature"
import {ManageAccess} from "../../management/component/ManageAccess"
import {KeeperPlugin, PlatformPlugin} from "../../node/api/NodeType"
import {useRouterDeploymentTemplateDelete, useRouterDeploymentTemplateList} from "../api/DeploymentHook"
import {CreationType, Template} from "../api/DeploymentType"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    row: {
        display: "flex", flexDirection: "column", gap: 0.5, padding: "10px 12px",
        border: 1, borderColor: "divider", borderRadius: 2,
        cursor: "pointer", "&:hover": {bgcolor: "action.hover", borderColor: "text.primary"},
    },
    title: {display: "flex", justifyContent: "space-between", alignItems: "center", gap: 1, overflow: "hidden"},
    name: {fontWeight: 600, whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis"},
    footer: {display: "flex", justifyContent: "space-between", alignItems: "flex-end", gap: 1, minHeight: "28px"},
    empty: {padding: "10px"},
    actions: {display: "flex", alignItems: "center"},
}

type Props = {
    keeper: KeeperPlugin,
    platform: PlatformPlugin,
    // onOpen runs a template: your own goes straight to the deploy form, a
    // shipped one opens as a copy first, since only your own are deployable
    onOpen: (template: Template) => void,
    onCopy: (template: Template) => void,
    onEdit: (template: Template) => void,
    onNew: () => void,
}

// DeploymentTemplateList shows the user's own templates and the shipped ones in
// one list. Picking a row runs it - system or not, they run the same way. The
// buttons cover the other two things you can do with one: edit your own, or
// copy any of them into one you own.
export function DeploymentTemplateList(props: Props) {
    const {keeper, platform, onOpen, onCopy, onEdit, onNew} = props
    const list = useRouterDeploymentTemplateList({keeper, platform})
    const remove = useRouterDeploymentTemplateDelete()

    return (
        <Box sx={SX.box}>
            {renderBody()}
            <ManageAccess feature={Feature.ManageDeploymentTemplateCreate}>
                <Button fullWidth={true} startIcon={<Add/>} onClick={onNew}>New template</Button>
            </ManageAccess>
        </Box>
    )

    function renderBody() {
        if (list.isError) return <ErrorSmart error={list.error}/>
        if (list.isPending) return <SkeletonGroup count={3}/>
        if (!list.data || list.data.length === 0) return renderEmpty()
        return list.data.map(renderRow)
    }

    function renderEmpty() {
        return <Box sx={SX.empty}><Note center={true}>No templates for this keeper yet</Note></Box>
    }

    function renderRow(template: Template) {
        return (
            <Box key={template.id} sx={SX.row} onClick={() => onOpen(template)}>
                <Box sx={SX.title}>
                    <Box sx={SX.name}>{template.name}</Box>
                    <InfoColorBoxRow>{renderLabels(template)}</InfoColorBoxRow>
                </Box>
                <Box sx={SX.footer}>
                    <Note>{template.description}</Note>
                    <Box sx={SX.actions}>{renderActions(template)}</Box>
                </Box>
            </Box>
        )
    }

    // NOTE: InfoColorBox is the app's own tag - the same one the cluster
    // overview uses - rather than a raw Chip that would look like a fourth
    // badge style. Every label shares one colour: they are all row metadata,
    // and colouring them differently made the row read as competing tags.
    function renderLabels(template: Template) {
        const platform = PlatformPluginOptions[template.platform]
        return (
            <>
                <InfoColorBox
                    label={platform ? platform.label.toLowerCase() : template.platform}
                    title={platform && `Deploys on ${platform.name ?? platform.label}`}
                />
                <InfoColorBox
                    label={isSystem(template) ? "system" : "manual"}
                    title={isSystem(template)
                        ? "Shipped with Ivory - copy it to make it yours"
                        : "Yours - editable and deletable"}
                />
            </>
        )
    }

    function renderCopyAction(template: Template) {
        return (
            <ManageAccess feature={Feature.ManageDeploymentTemplateCreate}>
                <CopyIconButton
                    tooltip={"Duplicate"}
                    onClick={(e) => handleAction(e, () => onCopy(template))}
                />
            </ManageAccess>
        )
    }

    // NOTE: a system template can be copied but never edited or deleted - it
    // is read-only by construction (it is computed, not stored), so those two
    // actions simply do not exist for it.
    function renderActions(template: Template) {
        if (isSystem(template)) return renderCopyAction(template)
        return (
            <>
                <ManageAccess feature={Feature.ManageDeploymentTemplateUpdate}>
                    <EditIconButton onClick={(e) => handleAction(e, () => onEdit(template))}/>
                </ManageAccess>
                {renderCopyAction(template)}
                <ManageAccess feature={Feature.ManageDeploymentTemplateDelete}>
                    <DeleteIconButton onClick={(e) => handleAction(e, () => remove.mutate(template.id))}/>
                </ManageAccess>
            </>
        )
    }

    // NOTE: the row itself opens the template, so an action inside it has to
    // stop the click from reaching the row
    function handleAction(event: Event | SyntheticEvent, action: () => void) {
        event.stopPropagation()
        action()
    }
}

function isSystem(template: Template) {
    return template.creation === CreationType.System
}
