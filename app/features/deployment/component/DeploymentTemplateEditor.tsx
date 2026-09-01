import {Add} from "@mui/icons-material"
import {Box, Button, TextField} from "@mui/material"

import {PaperBlue} from "../../../shared/component/box/PaperBlue"
import {TitleBox} from "../../../shared/component/box/TitleBox"
import {CopyIconButton, DeleteIconButton} from "../../../shared/component/button/IconButtons"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {DeployVar} from "../../node/api/NodeType"
import {TemplateCommand, TemplateDefaults, TemplateRequest} from "../api/DeploymentType"
import {DeploymentCommandEditor} from "./DeploymentCommandEditor"
import {DeploymentDefaultsGrid} from "./DeploymentDefaultsGrid"

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
            {renderDefaults()}
            {template.commands.map(renderNode)}
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

    // NOTE: usernames only, and one set for the whole template: every node of a
    // deployment authenticates the same way, so a per-command copy could only
    // disagree with its neighbours. A username here means the deployment ends
    // up with that account, so the deploy screen opens on it; naming none opens
    // that credential switched off. A password is never stored in a template.
    function renderDefaults() {
        return (
            <PaperBlue>
                <TitleBox
                    label={"Variables & Defaults"}
                    hint={"cluster-wide values available to every command in this template as {{variable}} - some are disabled because their value is always set at deploy time"}
                    island={true}
                    collapsible={false}
                >
                    <DeploymentDefaultsGrid
                        editable={editable}
                        fields={[
                            {
                                variable: DeployVar.Cluster, value: "", disabled: true,
                                hint: "set when you deploy - the same value reaches every node's command",
                                onChange: () => {},
                            },
                            undefined,
                            {
                                variable: DeployVar.KeeperUser, value: template.defaults?.keeperUser ?? "",
                                hint: "the account the keeper command creates - shown as a suggestion on the deploy screen, editable there",
                                onChange: (v) => handleDefaultsChange({keeperUser: v})},
                            {
                                variable: DeployVar.KeeperPass, value: "", disabled: true,
                                hint: "resolved from the keeper vault when you deploy - a password is never stored in a template",
                                onChange: () => {},
                            },
                            {
                                variable: DeployVar.DbUser, value: template.defaults?.dbUser ?? "",
                                hint: "the account the database command creates - shown as a suggestion on the deploy screen, editable there",
                                onChange: (v) => handleDefaultsChange({dbUser: v})},
                            {
                                variable: DeployVar.DbPass, value: "", disabled: true,
                                hint: "resolved from the database vault when you deploy - a password is never stored in a template",
                                onChange: () => {},
                            },
                        ]}
                    />
                </TitleBox>
            </PaperBlue>
        )
    }

    function renderNode(command: TemplateCommand, index: number) {
        return (
            <PaperBlue key={index}>
                <TitleBox
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
                </TitleBox>
            </PaperBlue>
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

    function handleDefaultsChange(defaults: TemplateDefaults) {
        onChange({...template, defaults: {...template.defaults, ...defaults}})
    }
}
