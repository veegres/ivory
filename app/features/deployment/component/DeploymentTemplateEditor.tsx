import {Add} from "@mui/icons-material"
import {Box, Button, TextField} from "@mui/material"

import {SubContentBox} from "../../../shared/component/box/SubContentBox"
import {CopyIconButton, DeleteIconButton} from "../../../shared/component/button/IconButtons"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {TemplateCommand, TemplateRequest} from "../api/DeploymentType"
import {DeploymentCommandEditor} from "./DeploymentCommandEditor"
import {DeploymentVariables} from "./DeploymentVariables"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 2},
}

type Props = {
    template: TemplateRequest,
    editable?: boolean,
    onChange: (template: TemplateRequest) => void,
    onCommandChange: (index: number, command: TemplateCommand) => void,
    onCommandAdd: (source?: TemplateCommand) => void,
    onCommandRemove: (index: number) => void,
}

// DeploymentTemplateEditor edits the ordered command list: one command per
// node, each written independently. A command has no identity beyond its
// position - the node it lands on is chosen at deploy time.
export function DeploymentTemplateEditor(props: Props) {
    const {template, editable = true, onChange, onCommandChange, onCommandAdd, onCommandRemove} = props

    return (
        <Box sx={SX.box}>
            {renderInfo()}
            <DeploymentVariables/>
            {template.commands.map(renderCommand)}
            {editable && renderAdd()}
        </Box>
    )

    function renderInfo() {
        return (
            <Box sx={SX.box}>
                <TextField
                    fullWidth={true}
                    size={"small"}
                    label={"Name"}
                    disabled={!editable}
                    value={template.name}
                    onChange={(e) => onChange({...template, name: e.target.value})}
                />
                <TextField
                    fullWidth={true}
                    size={"small"}
                    label={"Description"}
                    disabled={!editable}
                    value={template.description ?? ""}
                    onChange={(e) => onChange({...template, description: e.target.value})}
                />
            </Box>
        )
    }

    function renderCommand(command: TemplateCommand, index: number) {
        return (
            <SubContentBox
                key={index}
                label={`Node ${index + 1}`}
                renderActions={editable && renderCommandActions(command, index)}
                island={true}
                defaultOpen={index === 0}
            >
                <DeploymentCommandEditor
                    command={command}
                    editable={editable}
                    onChange={(updated) => onCommandChange(index, updated)}
                />
            </SubContentBox>
        )
    }

    function renderCommandActions(command: TemplateCommand, index: number) {
        return (
            <>
                <CopyIconButton
                    tooltip={"Duplicate node"}
                    onClick={() => onCommandAdd(command)}
                />
                <DeleteIconButton
                    disabled={template.commands.length <= 1}
                    onClick={() => onCommandRemove(index)}
                />
            </>
        )
    }

    function renderAdd() {
        return (
            <Button fullWidth={true} startIcon={<Add/>} onClick={() => onCommandAdd()}>
                Add node command
            </Button>
        )
    }
}
