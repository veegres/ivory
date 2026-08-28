import {Box} from "@mui/material"

import {CodeField} from "../../../shared/component/input/CodeField"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {TemplateCommand} from "../api/DeploymentType"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
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
            <CodeField
                label={"Command"}
                value={command.command}
                editable={editable}
                minHeight={"140px"}
                placeholder={"docker run -d --name {{name}} ..."}
                onUpdate={handleCommandChange}
            />
            {/* NOTE: the placeholder is a real post script rather than the
                word "optional" - it has to show that this one runs inside the
                container, unlike the command above it which starts one. It is
                deliberately plain shell: naming one engine's client made it
                read as the answer for that engine rather than an example. */}
            <CodeField
                label={"Post Script"}
                hint={"optional — runs in the container once every node is up"}
                value={command.postScript ?? ""}
                editable={editable}
                placeholder={`sh -c 'echo "{{name}} joined {{cluster}}"'`}
                onUpdate={handlePostScriptChange}
            />
        </Box>
    )

    function handleCommandChange(value: string) {
        onChange({...command, command: value})
    }

    function handlePostScriptChange(value: string) {
        onChange({...command, postScript: value})
    }
}
