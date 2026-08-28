import {Box, TextField} from "@mui/material"

import {InfoColorBox} from "../../../shared/component/box/InfoColorBox"
import {InfoColorBoxRow} from "../../../shared/component/box/InfoColorBoxRow"
import {SubContentBox} from "../../../shared/component/box/SubContentBox"
import {CodeField} from "../../../shared/component/input/CodeField"
import {FieldRow} from "../../../shared/component/input/FieldRow"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {interpolateCommand} from "../../../shared/helper/HelperUtils"
import {DeploymentPreviewNote} from "../../deployment/component/DeploymentPreviewNote"
import {DeployCredentials, DeployNode} from "../api/ClusterType"

const SX: SxPropsMap = {
    // NOTE: a frame with no heading - the fields name themselves and the nodes
    // are read in order, so numbering them was a label repeating the layout
    box: {
        display: "flex", flexDirection: "column", gap: 1,
        padding: 1, border: 1, borderColor: "divider", borderRadius: 2,
    },
    preview: {display: "flex", flexDirection: "column", gap: 1},
}

type Props = {
    node: DeployNode,
    cluster: string,
    withKeeperPort: boolean,
    // NOTE: a field is only marked red once the user has tried to deploy - an
    // untouched form should not open covered in errors
    showErrors: boolean,
    // NOTE: resolved by the dialog, which is where the credentials are chosen
    credentials: DeployCredentials,
    // NOTE: a name collision is the exception: it only exists once the user
    // has typed it, and it is the one error a red border cannot explain
    duplicate: boolean,
    onChange: (node: DeployNode) => void,
}

// ClusterDeployNode fills in one node of the template. The frame carries no
// heading: the nodes are read in the template's own order and the fields name
// themselves, so a number was a label restating the layout. Only the command,
// which is read-only here, collapses.
export function ClusterDeployNode(props: Props) {
    const {node, cluster, withKeeperPort, showErrors, duplicate, credentials, onChange} = props

    return (
        <Box sx={SX.box}>
            <FieldRow>
                <TextField
                    size={"small"}
                    label={"Name"}
                    placeholder={"etcd-1"}
                    value={node.name ?? ""}
                    error={duplicate || (showErrors && !node.name)}
                    helperText={duplicate && "another node already uses this name"}
                    onChange={(e) => onChange({...node, name: e.target.value})}
                />
                <TextField
                    size={"small"}
                    label={"Host"}
                    placeholder={"10.0.0.1"}
                    value={node.host}
                    error={showErrors && !node.host}
                    onChange={(e) => onChange({...node, host: e.target.value.toLowerCase()})}
                />
            </FieldRow>
            <FieldRow>
                {withKeeperPort && renderPort("Keeper Port", node.keeperPort, (v) => onChange({...node, keeperPort: v}))}
                {renderPort("Database Port", node.dbPort, (v) => onChange({...node, dbPort: v}))}
                {renderPort("SSH Port", node.sshPort, (v) => onChange({...node, sshPort: v}))}
            </FieldRow>
            {/* NOTE: the badge rides on the toggle's own row - it says what is
                inside the section, so it belongs to the line that opens it */}
            <SubContentBox label={"Preview"} renderActions={renderBadge()} dense={true}>
                {renderPreview()}
            </SubContentBox>
        </Box>
    )

    function renderBadge() {
        if (!node.postScript) return
        return (
            <InfoColorBoxRow>
                <InfoColorBox label={"post script"} title={"Runs inside the container once this node is up"}/>
            </InfoColorBoxRow>
        )
    }

    // NOTE: the hint sits above the code, not under it - it says how to read
    // what follows, which is no use once you have already read it
    function renderPreview() {
        return (
            <Box sx={SX.preview}>
                <DeploymentPreviewNote/>
                <CodeField
                    label={"Command"}
                    value={getPreview(node.command)}
                    editable={false}
                    minHeight={"120px"}
                />
                {node.postScript && (
                    <CodeField
                        label={"Post Script"}
                        hint={"runs in the container once this node is up"}
                        value={getPreview(node.postScript)}
                        editable={false}
                    />
                )}
            </Box>
        )
    }

    function renderPort(label: string, value: number | undefined, onPortChange: (value?: number) => void) {
        return (
            <TextField
                size={"small"}
                type={"number"}
                label={label}
                value={value ?? ""}
                onChange={(e) => onPortChange(Number(e.target.value) || undefined)}
            />
        )
    }

    function getPreview(text: string) {
        return interpolateCommand(text, {
            cluster,
            name: node.name,
            host: node.host,
            sshPort: node.sshPort,
            keeperPort: node.keeperPort,
            dbPort: node.dbPort,
            dbUser: credentials.user,
            dbPass: credentials.pass,
        })
    }
}
