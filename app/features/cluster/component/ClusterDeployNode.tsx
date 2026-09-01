import {TextField} from "@mui/material"

import {PaperBlue} from "../../../shared/component/box/PaperBlue"
import {FieldRow} from "../../../shared/component/input/FieldRow"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {DeployValues} from "../../../shared/helper/HelperUtils"
import {DeploymentCommandPreview} from "../../deployment/component/DeploymentCommandPreview"
import {DeployNode, DeployPreviewCredentials} from "../api/ClusterType"

const SX: SxPropsMap = {
    // NOTE: a frame with no heading - the fields name themselves and the nodes
    // are read in order, so numbering them was a label repeating the layout
    box: {
        display: "flex", flexDirection: "column", gap: 1,
        padding: 1, border: 1, borderColor: "divider", borderRadius: 2,
    },
}

type Props = {
    node: DeployNode,
    cluster: string,
    // NOTE: a field is only marked red once the user has tried to deploy - an
    // untouched form should not open covered in errors
    showErrors: boolean,
    // NOTE: resolved by the dialog, which is where the credentials are chosen
    credentials: DeployPreviewCredentials,
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
    const {node, cluster, showErrors, duplicate, credentials, onChange} = props

    return (
        <PaperBlue sx={SX.box}>
            <FieldRow>
                <TextField
                    size={"small"}
                    label={"Name"}
                    placeholder={"etcd-1"}
                    value={node.name}
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
                {renderPort("Keeper Port", node.keeperPort, (v) => onChange({...node, keeperPort: v}))}
                {renderPort("Database Port", node.dbPort, (v) => onChange({...node, dbPort: v}))}
                {renderPort("SSH Port", node.sshPort, (v) => onChange({...node, sshPort: v}))}
            </FieldRow>
            <DeploymentCommandPreview command={node.command} postScripts={node.postScripts} values={getValues()}/>
        </PaperBlue>
    )

    function renderPort(label: string, value: number | undefined, onPortChange: (value?: number) => void) {
        return (
            <TextField
                size={"small"}
                type={"number"}
                label={label}
                value={value ?? ""}
                error={showErrors && !value}
                onChange={(e) => onPortChange(Number(e.target.value) || undefined)}
            />
        )
    }

    function getValues(): DeployValues {
        return {
            cluster,
            name: node.name,
            host: node.host,
            sshPort: node.sshPort,
            keeperPort: node.keeperPort,
            dbPort: node.dbPort,
            keeperUser: credentials.keeper.user,
            keeperPass: credentials.keeper.pass,
            dbUser: credentials.database.user,
            dbPass: credentials.database.pass,
        }
    }
}
