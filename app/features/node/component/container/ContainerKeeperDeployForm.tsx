import {Box, Button, TextField, ToggleButton, ToggleButtonGroup} from "@mui/material"
import {useMemo, useState} from "react"

import {DialogLogsScreen} from "../../../../shared/component/box/DialogLogsScreen"
import {DialogScreen} from "../../../../shared/component/box/DialogScreen"
import {Hint} from "../../../../shared/component/box/Hint"
import {PaperBlue} from "../../../../shared/component/box/PaperBlue"
import {TitleBox} from "../../../../shared/component/box/TitleBox"
import {FieldRow} from "../../../../shared/component/input/FieldRow"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {DeployValues, getShortUuid} from "../../../../shared/helper/HelperUtils"
import {useDeployVaultCredentials} from "../../../deployment/api/DeploymentHook"
import {Template} from "../../../deployment/api/DeploymentType"
import {DeploymentCommandPreview} from "../../../deployment/component/DeploymentCommandPreview"
import {useRouterNodeKeeperDeploy} from "../../api/NodeHook"
import {DeployVar, KeeperPlugin, PlatformVaultConnection} from "../../api/NodeType"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    column: {display: "flex", flexDirection: "column", gap: 1, marginTop: "5px"},
    hint: {textTransform: "uppercase"},
    node: {
        display: "flex", flexDirection: "column", gap: 1,
        padding: 1, border: 1, borderColor: "divider", borderRadius: 2,
    },
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
    keeperId?: string,
    databaseId?: string,
    sshKeyId?: string,
    logs?: string[],
    onDeployed: (logs: string[]) => void,
}

export function ContainerKeeperDeployForm(props: Props) {
    const {connection, plugin, cluster, node, template, keeperId, databaseId, sshKeyId, logs, onDeployed} = props
    const [index, setIndex] = useState(0)
    const [dcs, setDcs] = useState(template.defaults?.dcs ?? "")

    const nodeDeploy = useRouterNodeKeeperDeploy(connection, onDeployed)
    const credentials = useDeployVaultCredentials(keeperId, databaseId)

    const command = template.commands[index]
    const withKeeperCredentials = !!template.defaults?.keeperUser
    const withDbCredentials = !!template.defaults?.dbUser
    const keeperPort = command.defaults?.keeperPort
    const dbPort = command.defaults?.dbPort
    const name = command.defaults?.name?.trim() || node
    const withDcs = useMemo(handleMemoWithDcs, [command])

    if (logs) return <DialogLogsScreen logs={logs}/>

    return (
        <DialogScreen renderActions={renderActions()}>
            <Box sx={SX.box}>
                {renderClusterInfo()}
                {renderNodeChooser()}
                {renderNode()}
            </Box>
        </DialogScreen>
    )

    function renderActions() {
        const portsMissing = !keeperPort || !dbPort
        const dcsMissing = withDcs && !dcs.trim()
        return (
            <Button
                loading={nodeDeploy.isPending}
                onClick={handleAction}
                disabled={portsMissing || dcsMissing || !sshKeyId || !command.command.trim()}
            >
                Deploy
            </Button>
        )
    }

    function renderNodeChooser() {
        if (template.commands.length === 1) return
        return (
            <PaperBlue sx={SX.chooser}>
                <Hint sx={SX.hint} center={true}>Select the node you wish to deploy from the template</Hint>
                <ToggleButtonGroup fullWidth={true} exclusive={true} value={index} onChange={(_, v) => setIndex(v ?? index)}>
                    {template.commands.map((_, i) => (
                        <ToggleButton key={i} sx={SX.toggleButton} value={i}>Node {i + 1}</ToggleButton>
                    ))}
                </ToggleButtonGroup>
            </PaperBlue>
        )
    }

    function renderClusterInfo() {
        return (
            <PaperBlue>
                <TitleBox label={"Cluster"} island={true} collapsible={false}>
                    <Box sx={SX.column}>
                        <TextField fullWidth label={"Cluster Name"} value={cluster} disabled={true}/>
                        <FieldRow>
                            {withKeeperCredentials && renderVaultField("Keeper Credentials", keeperId)}
                            {withDbCredentials && renderVaultField("Database Credentials", databaseId)}
                            {renderVaultField("SSH Credentials", sshKeyId)}
                        </FieldRow>
                        {withDcs && renderDcs()}
                    </Box>
                </TitleBox>
            </PaperBlue>
        )
    }

    function renderNode() {
        return (
            <PaperBlue sx={SX.node}>
                <FieldRow>
                    <TextField label={"Name"} value={name} disabled={true}/>
                    <TextField label={"Host"} value={connection.host} disabled={true}/>
                </FieldRow>
                <FieldRow>
                    {renderPort("Keeper Port", keeperPort)}
                    {renderPort("Database Port", dbPort)}
                    {renderPort("SSH Port", connection.port)}
                </FieldRow>
                <DeploymentCommandPreview
                    command={command.command}
                    postScripts={command.postScripts}
                    values={getValues()}
                    defaultOpen={true}
                />
            </PaperBlue>
        )
    }

    function renderDcs() {
        return (
            <TextField
                fullWidth={true}
                label={"DCS"}
                placeholder={"etcd1:2379, etcd2:2379, etcd3:2379"}
                helperText={"This cluster relies on the Distributed Consensus Store (DCS)"}
                value={dcs}
                onChange={(e) => setDcs(e.target.value)}
            />
        )
    }

    function renderVaultField(label: string, vaultId?: string) {
        return <TextField fullWidth label={label} value={getShortUuid(vaultId ?? "none")} disabled={true}/>
    }

    function renderPort(label: string, value?: number) {
        return <TextField type={"number"} label={label} value={value ?? ""} disabled={true}/>
    }

    function handleMemoWithDcs() {
        return [command.command, ...(command.postScripts ?? [])].some(text => text.includes(DeployVar.Dcs))
    }

    function handleAction() {
        if (!keeperPort || !dbPort) return
        nodeDeploy.mutate({
            plugin,
            cluster,
            dcs,
            name,
            connection: {...connection, platform: template.platform},
            command: command.command,
            postScripts: command.postScripts,
            keeperPort,
            dbPort,
            vaults: {keeperId, databaseId, sshKeyId: sshKeyId ?? ""},
        })
    }

    function getValues(): DeployValues {
        return {
            cluster,
            dcs,
            name,
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
