import {Box, TextField} from "@mui/material"

import {Note} from "../../../shared/component/box/Note"
import {CodeField} from "../../../shared/component/input/CodeField"
import {FieldRow} from "../../../shared/component/input/FieldRow"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {TemplateCommand, TemplateDefaults} from "../api/DeploymentType"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    defaults: {display: "flex", flexDirection: "column", gap: 0.5},
}

type Props = {
    command: TemplateCommand,
    editable?: boolean,
    onChange: (command: TemplateCommand) => void,
}

export function DeploymentCommandEditor(props: Props) {
    const {command, editable = true, onChange} = props

    return (
        <Box sx={SX.box}>
            {renderDefaults()}
            <CodeField
                label={"Command"}
                value={command.command}
                editable={editable}
                minHeight={"140px"}
                placeholder={"docker run -d --name {{name}} ..."}
                onUpdate={handleCommandChange}
            />
{/* NOTE: one command per line, because each line runs as its own
                execution rather than as a script - a distroless image has no
                shell to chain them with. The placeholder is a real step rather
                than the word "optional": it has to show that these run inside
                the container, unlike the command above which starts one. */}
            <CodeField
                label={"Post Script"}
                hint={"optional — one command per line, each runs in the container once every node is up"}
                value={getPostScriptText()}
                editable={editable}
                placeholder={"etcdctl --endpoints=http://localhost:{{dbPort}} auth enable"}
                onUpdate={handlePostScriptChange}
            />
        </Box>
    )

    // NOTE: these are what the deploy screen puts in this node's card before
    // anyone types - the one place a template can say that its second node
    // answers on a different port than its first, which the keeper plugin's
    // single set of defaults cannot express. Left empty they fall back to those.
    function renderDefaults() {
        return (
            <Box sx={SX.defaults}>
                <Note>Prefills this node on the deploy screen. Leave a field empty to use the engine's own default.</Note>
                <FieldRow>
                    <TextField
                        size={"small"}
                        label={"Node Name"}
                        disabled={!editable}
                        value={command.defaults?.name ?? ""}
                        onChange={(e) => handleDefaultsChange({name: e.target.value})}
                    />
                    <TextField
                        size={"small"}
                        type={"number"}
                        label={"Keeper Port"}
                        disabled={!editable}
                        value={command.defaults?.keeperPort ?? ""}
                        onChange={(e) => handleDefaultsChange({keeperPort: getPort(e.target.value)})}
                    />
                    <TextField
                        size={"small"}
                        type={"number"}
                        label={"Database Port"}
                        disabled={!editable}
                        value={command.defaults?.dbPort ?? ""}
                        onChange={(e) => handleDefaultsChange({dbPort: getPort(e.target.value)})}
                    />
                </FieldRow>
            </Box>
        )
    }

    function handleCommandChange(value: string) {
        onChange({...command, command: value})
    }

    function handlePostScriptChange(value: string) {
        const steps = value.split("\n").map((line) => line.trim()).filter((line) => line !== "")
        onChange({...command, postScripts: steps.length > 0 ? steps : undefined})
    }

    function handleDefaultsChange(defaults: TemplateDefaults) {
        onChange({...command, defaults: {...command.defaults, ...defaults}})
    }

    function getPostScriptText() {
        return (command.postScripts ?? []).join("\n")
    }

    function getPort(value: string) {
        return value === "" ? undefined : Number(value)
    }
}
