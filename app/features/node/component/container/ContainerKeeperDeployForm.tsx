import {Box, Button, TextField, ToggleButton, ToggleButtonGroup} from "@mui/material"
import {useState} from "react"

import {DialogScreen} from "../../../../shared/component/box/DialogScreen"
import {InfoColorBox} from "../../../../shared/component/box/InfoColorBox"
import {InfoColorBoxRow} from "../../../../shared/component/box/InfoColorBoxRow"
import {Note} from "../../../../shared/component/box/Note"
import {SubContentBox} from "../../../../shared/component/box/SubContentBox"
import {TitledBox} from "../../../../shared/component/box/TitledBox"
import {CodeField} from "../../../../shared/component/input/CodeField"
import {FieldRow} from "../../../../shared/component/input/FieldRow"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {DeployPasswordMask, getShortUuid, interpolateCommand} from "../../../../shared/helper/HelperUtils"
import {Template} from "../../../deployment/api/DeploymentType"
import {DeploymentPreviewNote} from "../../../deployment/component/DeploymentPreviewNote"
import {useRouterVault} from "../../../vault/api/VaultHook"
import {VaultType} from "../../../vault/api/VaultType"
import {useRouterNodeKeeperDeploy} from "../../api/NodeHook"
import {KeeperDeploySpecResponse, KeeperPlugin, PlatformVaultConnection} from "../../api/NodeType"
import {ContainerKeeperDeployResponse} from "./ContainerKeeperDeployResponse"

const SX: SxPropsMap = {
    subContent: {display: "flex", flexDirection: "column"},
    // NOTE: the frame carries no heading, exactly as the cluster deploy's node
    // does - the fields name themselves, and there is only ever one node here
    node: {
        display: "flex", flexDirection: "column", gap: 1,
        padding: 1, border: 1, borderColor: "divider", borderRadius: 2,
    },
    preview: {display: "flex", flexDirection: "column", gap: 1},
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
// are the keeper plugin's own defaults, and the command belongs to the
// template - so every field here is a read-only account of the deployment.
export function ContainerKeeperDeployForm(props: Props) {
    const {connection, plugin, cluster, node, template, spec, keeperId, databaseId, sshKeyId, logs, onDeployed} = props
    // NOTE: which of the template's nodes runs on this host - the first one
    // unless the user says otherwise, since a template is written in node order
    const [index, setIndex] = useState(0)

    const nodeDeploy = useRouterNodeKeeperDeploy(connection, onDeployed)
    const keeperVaults = useRouterVault(VaultType.KEEPER_PASSWORD)
    const dbVaults = useRouterVault(VaultType.DATABASE_PASSWORD)

    const command = template.commands[index]
    const withKeeperCredentials = spec.keeperCredentials
    const withDbCredentials = spec.dbCredentials

    if (logs) return <ContainerKeeperDeployResponse logs={logs}/>

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
                    {renderPort("Keeper Port", spec.keeperPort)}
                    {renderPort("Database Port", spec.dbPort)}
                    {renderPort("SSH Port", connection.port)}
                </FieldRow>
                {/* NOTE: the badge rides on the toggle's own row - it says what
                    is inside the section, so it belongs to the line that opens it */}
                <SubContentBox label={"Preview"} renderActions={renderBadge()} defaultOpen={true} dense={true}>
                    {renderPreview()}
                </SubContentBox>
            </Box>
        )
    }

    function renderVaultField(label: string, vaultId?: string) {
        return <TextField fullWidth size={"small"} label={label} value={getShortUuid(vaultId ?? "none")} disabled={true}/>
    }

    function renderPort(label: string, value: number) {
        return <TextField size={"small"} type={"number"} label={label} value={value} disabled={true}/>
    }

    function renderBadge() {
        if (!command.postScript) return
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
                    value={interpolateCommand(command.command, getValues())}
                    editable={false}
                    minHeight={"120px"}
                />
                {command.postScript && (
                    <CodeField
                        label={"Post Script"}
                        hint={"runs in the container once this node is up"}
                        value={interpolateCommand(command.postScript, getValues())}
                        editable={false}
                    />
                )}
            </Box>
        )
    }

    function handleAction() {
        nodeDeploy.mutate({
            plugin,
            cluster,
            name: node,
            connection,
            command: command.command,
            postScript: command.postScript,
            keeperPort: spec.keeperPort,
            dbPort: spec.dbPort,
            vaults: {keeperId, databaseId, sshKeyId: sshKeyId ?? ""},
        })
    }

    function getValues() {
        return {
            cluster,
            name: node,
            host: connection.host,
            sshPort: connection.port,
            keeperPort: spec.keeperPort,
            dbPort: spec.dbPort,
            keeperUser: getKeeperVault()?.username,
            keeperPass: getKeeperVault() && DeployPasswordMask,
            dbUser: getDbVault()?.username,
            dbPass: getDbVault() && DeployPasswordMask,
        }
    }

    function getKeeperVault() {
        return keeperId ? keeperVaults.data?.[keeperId] : undefined
    }

    function getDbVault() {
        return databaseId ? dbVaults.data?.[databaseId] : undefined
    }
}
