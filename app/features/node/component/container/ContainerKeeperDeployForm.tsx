import {Box, Button, TextField, ToggleButton, ToggleButtonGroup} from "@mui/material"
import {useState} from "react"

import {DialogLogsScreen} from "../../../../shared/component/box/DialogLogsScreen"
import {DialogScreen} from "../../../../shared/component/box/DialogScreen"
import {Note} from "../../../../shared/component/box/Note"
import {TitledBox} from "../../../../shared/component/box/TitledBox"
import {FieldRow} from "../../../../shared/component/input/FieldRow"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {DeployValues, getShortUuid} from "../../../../shared/helper/HelperUtils"
import {useDeployVaultCredentials} from "../../../deployment/api/DeploymentHook"
import {Template} from "../../../deployment/api/DeploymentType"
import {DeploymentCommandPreview} from "../../../deployment/component/DeploymentCommandPreview"
import {useRouterNodeKeeperDeploy} from "../../api/NodeHook"
import {KeeperDeploySpecResponse, KeeperPlugin, PlatformVaultConnection} from "../../api/NodeType"

const SX: SxPropsMap = {
    subContent: {display: "flex", flexDirection: "column"},
    // NOTE: the frame carries no heading, exactly as the cluster deploy's node
    // does - the fields name themselves, and there is only ever one node here
    node: {
        display: "flex", flexDirection: "column", gap: 1,
        padding: 1, border: 1, borderColor: "divider", borderRadius: 2,
    },
    // NOTE: the same frame as the node card below it - it is a section of the
    // screen in its own right, not a caption floating between two boxes
    chooser: {
        display: "flex", flexDirection: "column", gap: 0.5,
        padding: 1, border: 1, borderColor: "divider", borderRadius: 2,
    },
    toggleButton: {padding: "0px 10px"},
}

type Props = {
    connection: PlatformVaultConnection,
    plugin: KeeperPlugin,
    cluster: string,
    node: string,
    template: Template,
    spec: KeeperDeploySpecResponse,
    keeperId?: string,
    databaseId?: string,
    sshKeyId?: string,
    logs?: string[],
    onDeployed: (logs: string[]) => void,
}

// ContainerKeeperDeployForm shows what deploying this node will run. It has
// nothing to fill in: the node is the one the dialog was opened on, its ports
// are the ones the chosen command states, and the command belongs to the
// template - so every field here is a read-only account of the deployment.
export function ContainerKeeperDeployForm(props: Props) {
    const {connection, plugin, cluster, node, template, spec, keeperId, databaseId, sshKeyId, logs, onDeployed} = props
    // NOTE: which of the template's nodes runs on this host - the first one
    // unless the user says otherwise, since a template is written in node order
    const [index, setIndex] = useState(0)

    const nodeDeploy = useRouterNodeKeeperDeploy(connection, onDeployed)
    const credentials = useDeployVaultCredentials(keeperId, databaseId)

    const command = template.commands[index]
    const withKeeperCredentials = spec.keeperCredentials
    const withDbCredentials = spec.dbCredentials
    // NOTE: the ports belong to the chosen node, not to the engine - node 2 of
    // a single-host template answers on its own pair, and switching the toggle
    // above has to move them with it
    const keeperPort = command.defaults?.keeperPort ?? spec.keeperPort
    const dbPort = command.defaults?.dbPort ?? spec.dbPort

    if (logs) return <DialogLogsScreen logs={logs}/>

    return (
        <DialogScreen renderActions={renderActions()}>
            <Box sx={[SX.subContent, {gap: 2}]}>
                {renderClusterInfo()}
                {renderNodeChooser()}
                {renderNode()}
            </Box>
        </DialogScreen>
    )

    function renderActions() {
        const keeperVaultMissing = withKeeperCredentials && !keeperId
        const dbVaultMissing = withDbCredentials && !databaseId
        return (
            <Button
                loading={nodeDeploy.isPending}
                onClick={handleAction}
                disabled={keeperVaultMissing || dbVaultMissing || !sshKeyId || !command.command.trim()}
            >
                Deploy
            </Button>
        )
    }

    // NOTE: a template describes a whole cluster, this screen deploys one host
    // of it - so which of its nodes that is has to be the user's answer, not an
    // assumption. It sits above the card it changes, and disappears for a
    // one-node template, where there is nothing to choose.
    function renderNodeChooser() {
        if (template.commands.length === 1) return
        return (
            <Box sx={SX.chooser}>
                <Note center={true}>Pick which of the template's nodes is deployed here</Note>
                <ToggleButtonGroup fullWidth={true} size={"small"} exclusive={true} value={index} onChange={(_, v) => setIndex(v ?? index)}>
                    {template.commands.map((_, i) => (
                        <ToggleButton key={i} sx={SX.toggleButton} value={i}>Node {i + 1}</ToggleButton>
                    ))}
                </ToggleButtonGroup>
            </Box>
        )
    }

    function renderClusterInfo() {
        return (
            <TitledBox title={"Cluster"} island={true}>
                <Box sx={[SX.subContent, {gap: 1}]}>
                    <TextField fullWidth size={"small"} label={"Cluster Name"} value={cluster} disabled={true}/>
                    <FieldRow>
                        {withKeeperCredentials && renderVaultField("Keeper Credentials", keeperId)}
                        {withDbCredentials && renderVaultField("Database Credentials", databaseId)}
                        {renderVaultField("SSH Credentials", sshKeyId)}
                    </FieldRow>
                </Box>
            </TitledBox>
        )
    }

    // NOTE: the same card the cluster deploy shows a node in - a bordered
    // frame, the fields, then the command folded away behind a toggle. Here it
    // opens expanded: this screen has one node and nothing to fill in, so the
    // command is the only thing on it left to read.
    function renderNode() {
        return (
            <Box sx={SX.node}>
                <FieldRow>
                    <TextField size={"small"} label={"Name"} value={node} disabled={true}/>
                    <TextField size={"small"} label={"Host"} value={connection.host} disabled={true}/>
                </FieldRow>
                <FieldRow>
                    {renderPort("Keeper Port", keeperPort)}
                    {renderPort("Database Port", dbPort)}
                    {renderPort("SSH Port", connection.port)}
                </FieldRow>
                <DeploymentCommandPreview
                    command={command.command}
                    postScript={command.postScript}
                    values={getValues()}
                    defaultOpen={true}
                />
            </Box>
        )
    }

    function renderVaultField(label: string, vaultId?: string) {
        return <TextField fullWidth size={"small"} label={label} value={getShortUuid(vaultId ?? "none")} disabled={true}/>
    }

    function renderPort(label: string, value: number) {
        return <TextField size={"small"} type={"number"} label={label} value={value} disabled={true}/>
    }

    function handleAction() {
        nodeDeploy.mutate({
            plugin,
            cluster,
            name: node,
            connection,
            command: command.command,
            postScript: command.postScript,
            keeperPort,
            dbPort,
            vaults: {keeperId, databaseId, sshKeyId: sshKeyId ?? ""},
        })
    }

    function getValues(): DeployValues {
        return {
            cluster,
            name: node,
            host: connection.host,
            sshPort: connection.port,
            keeperPort,
            dbPort,
            keeperUser: credentials.keeper.user,
            keeperPass: credentials.keeper.pass,
            dbUser: credentials.database.user,
            dbPass: credentials.database.pass,
        }
    }
}
