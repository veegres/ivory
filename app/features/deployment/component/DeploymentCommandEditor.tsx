import {Box} from "@mui/material"

import {SubContentBox} from "../../../shared/component/box/SubContentBox"
import {CodeField} from "../../../shared/component/input/CodeField"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {DeployVar} from "../../node/api/NodeType"
import {CommandDefaults, TemplateCommand} from "../api/DeploymentType"
import {DeploymentDefaultsGrid} from "./DeploymentDefaultsGrid"

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
            {renderDefaults()}
            <CodeField
                label={"Command"}
                value={command.command}
                editable={editable}
                minHeight={"140px"}
                placeholder={"docker run -d --name {{name}} ..."}
                onUpdate={handleCommandChange}
            />
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

    function renderDefaults() {
        return (
            <SubContentBox
                label={"Variables & Defaults"}
                dense={true}
                collapsible={false}
            >
                <DeploymentDefaultsGrid
                    editable={editable}
                    fields={[
                        {
                            variable: DeployVar.Name, value: command.defaults?.name ?? "",
                            hint: "this node's identity - what a member list referencing it should use",
                            onChange: (v) => handleDefaultsChange({name: v}),
                        },
                        {
                            variable: DeployVar.Host, value: "", disabled: true,
                            hint: "the machine you deploy this node onto - never part of a template",
                            onChange: () => {},
                        },
                        {
                            variable: DeployVar.SshPort, value: getPortValue(command.defaults?.sshPort), numeric: true,
                            hint: "the port Ivory reaches this node over ssh",
                            onChange: (v) => handleDefaultsChange({sshPort: getPort(v)}),
                        },
                        {
                            variable: DeployVar.KeeperPort, value: getPortValue(command.defaults?.keeperPort), numeric: true,
                            hint: "the port this node's keeper listens on",
                            onChange: (v) => handleDefaultsChange({keeperPort: getPort(v)}),
                        },
                        {
                            variable: DeployVar.DbPort, value: getPortValue(command.defaults?.dbPort), numeric: true,
                            hint: "the port this node's database listens on",
                            onChange: (v) => handleDefaultsChange({dbPort: getPort(v)}),
                        },
                    ]}
                />
            </SubContentBox>
        )
    }

    function handleCommandChange(value: string) {
        onChange({...command, command: value})
    }

    function handlePostScriptChange(value: string) {
        const steps = value.split("\n").map((line) => line.trim()).filter((line) => line !== "")
        onChange({...command, postScripts: steps.length > 0 ? steps : undefined})
    }

    function handleDefaultsChange(defaults: CommandDefaults) {
        onChange({...command, defaults: {...command.defaults, ...defaults}})
    }

    function getPostScriptText() {
        return (command.postScripts ?? []).join("\n")
    }

    function getPort(value: string) {
        return value === "" ? undefined : Number(value)
    }

    function getPortValue(port?: number) {
        return port ? String(port) : ""
    }
}
