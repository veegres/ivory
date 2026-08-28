import {Rocket} from "@mui/icons-material"
import {Box, Button, TextField} from "@mui/material"
import {useEffect, useState} from "react"

import {ErrorSmart} from "../../../../shared/component/box/ErrorSmart"
import {Logs} from "../../../../shared/component/box/Logs"
import {Note} from "../../../../shared/component/box/Note"
import {TitledBox} from "../../../../shared/component/box/TitledBox"
import {DialogButton} from "../../../../shared/component/button/DialogButton"
import {CodeField} from "../../../../shared/component/input/CodeField"
import {FieldRow} from "../../../../shared/component/input/FieldRow"
import {SkeletonGroup} from "../../../../shared/component/progress/SkeletonGroup"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {DeployPasswordMask, getShortUuid, interpolateCommand} from "../../../../shared/helper/HelperUtils"
import {Template, TemplateCommand} from "../../../deployment/api/DeploymentType"
import {DeploymentPreviewNote} from "../../../deployment/component/DeploymentPreviewNote"
import {useDeploymentTemplatePicker} from "../../../deployment/component/DeploymentTemplatePicker"
import {useRouterVault} from "../../../vault/api/VaultHook"
import {VaultType} from "../../../vault/api/VaultType"
import {useRouterNodeKeeperDeploy, useRouterNodeKeeperDeploySpec} from "../../api/NodeHook"
import {KeeperPlugin, PlatformPlugin, PlatformVaultConnection} from "../../api/NodeType"

const SX: SxPropsMap = {
    subContent: {display: "flex", flexDirection: "column"},
}

type Props = {
    connection: PlatformVaultConnection,
    plugin: KeeperPlugin,
    cluster: string,
    node: string,
    databaseId?: string,
    sshKeyId?: string,
}

// ContainerKeeperDeploy deploys a keeper onto a single existing node: it picks
// one command out of a template and runs it here, calling node's own
// KeeperDeploy directly - no cluster endpoint involved, which is why it lives
// in the node/container feature rather than cluster.
export function ContainerKeeperDeploy(props: Props) {
    const {connection, plugin, cluster, node, databaseId, sshKeyId} = props

    const [template, setTemplate] = useState<Template>()
    const [command, setCommand] = useState<TemplateCommand>()
    const [name, setName] = useState("")
    const [keeperPort, setKeeperPort] = useState<string>("")
    const [dbPort, setDbPort] = useState<string>("")
    const [submitted, setSubmitted] = useState(false)

    const platform = connection.platform ?? PlatformPlugin.LINUX
    const nodeDeploy = useRouterNodeKeeperDeploy(connection)
    const spec = useRouterNodeKeeperDeploySpec(plugin)
    const picker = useDeploymentTemplatePicker({
        keeper: plugin,
        platform,
        hint: "Pick a template to deploy this node - it uses the template's first command",
        onPick: setTemplate,
    })
    const dbVaults = useRouterVault(VaultType.DATABASE_PASSWORD)

    useEffect(handleEffectSpec, [spec.data, node])
    useEffect(handleEffectTemplate, [template])

    const withKeeperPort = spec.data?.keeperPort !== undefined
    const withDbCredentials = spec.data?.credentials ?? false

    return (
        <DialogButton
            title={"DEPLOY CONTAINER"}
            variant={"button"}
            renderActions={!picker.editing && command && renderActions()}
            icon={<Rocket fontSize={"small"}/>}
            back={!!nodeDeploy.data || !!template || picker.editing}
            onBackClick={handleBack}
        >
            {renderContent()}
        </DialogButton>
    )

    function renderContent() {
        if (nodeDeploy.data) return <Logs logs={nodeDeploy.data} height={570} auto={false}/>
        if (spec.isError) return <ErrorSmart error={spec.error}/>
        if (spec.isPending) return <SkeletonGroup count={3}/>
        if (!template || !command) return picker.render()
        return renderBody(command)
    }

    function renderActions() {
        const dbVaultMissing = withDbCredentials && !databaseId
        return (
            <Button
                loading={nodeDeploy.isPending}
                onClick={handleAction}
                disabled={dbVaultMissing || !sshKeyId}
            >
                Deploy
            </Button>
        )
    }

    function renderBody(current: TemplateCommand) {
        return (
            <Box sx={[SX.subContent, {gap: 2}]}>
                {renderClusterInfo()}
                {renderNodeFields()}
                {renderPreview(current)}
                {template && template.commands.length > 1 && renderCommandNote()}
            </Box>
        )
    }

    // NOTE: a note, not a picker - this screen deploys one node, and which of
    // a template's commands that is has no answer beyond "the first"
    function renderCommandNote() {
        return (
            <Note center={true}>
                This template has {template?.commands.length} commands - the first one is used here
            </Note>
        )
    }

    function renderClusterInfo() {
        return (
            <TitledBox title={"Cluster"} island={true}>
                <Box sx={[SX.subContent, {gap: 1}]}>
                    <TextField fullWidth size={"small"} label={"Cluster Name"} value={cluster} disabled={true}/>
                    <FieldRow>
                        {withDbCredentials && (
                            <TextField
                                fullWidth
                                size={"small"}
                                label={"Database Credentials"}
                                value={getShortUuid(databaseId ?? "none")}
                                disabled={true}
                            />
                        )}
                        <TextField
                            fullWidth
                            size={"small"}
                            label={"SSH Credentials"}
                            value={getShortUuid(sshKeyId ?? "none")}
                            disabled={true}
                        />
                    </FieldRow>
                </Box>
            </TitledBox>
        )
    }

    function renderNodeFields() {
        return (
            <TitledBox title={"Node"} island={true}>
                <Box sx={[SX.subContent, {gap: 1}]}>
                    <TextField
                        fullWidth
                        size={"small"}
                        label={"Name"}
                        value={name}
                        error={submitted && !name.trim()}
                        onChange={e => setName(e.target.value)}
                    />
                    <FieldRow>
                        {withKeeperPort && (
                            <TextField
                                size={"small"}
                                type={"number"}
                                label={"Keeper Port"}
                                value={keeperPort}
                                onChange={e => setKeeperPort(e.target.value)}
                            />
                        )}
                        <TextField
                            size={"small"}
                            type={"number"}
                            label={"Database Port"}
                            value={dbPort}
                            onChange={e => setDbPort(e.target.value)}
                        />
                    </FieldRow>
                </Box>
            </TitledBox>
        )
    }

    // NOTE: the command belongs to the template and is read-only here - what
    // this screen supplies is the node config it interpolates, so it is shown
    // filled in rather than editable. The hint sits above the code: it says how
    // to read what follows, which is no use once you have already read it.
    function renderPreview(current: TemplateCommand) {
        return (
            <TitledBox title={"Preview"} island={true}>
                <Box sx={[SX.subContent, {gap: 1}]}>
                    <DeploymentPreviewNote/>
                    <CodeField
                        label={"Command"}
                        value={interpolateCommand(current.command, getValues())}
                        editable={false}
                        minHeight={"120px"}
                    />
                    {current.postScript && (
                        <CodeField
                            label={"Post Script"}
                            hint={"runs in the container once this node is up"}
                            value={interpolateCommand(current.postScript, getValues())}
                            editable={false}
                        />
                    )}
                </Box>
            </TitledBox>
        )
    }

    function getValues() {
        return {
            cluster,
            name,
            host: connection.host,
            sshPort: connection.port,
            keeperPort: Number(keeperPort) || undefined,
            dbPort: Number(dbPort) || undefined,
            dbUser: getDbVault()?.username,
            dbPass: getDbVault() && DeployPasswordMask,
        }
    }

    function getDbVault() {
        return databaseId ? dbVaults.data?.[databaseId] : undefined
    }

    function isReady() {
        return !!command?.command.trim() && !!name.trim()
    }

    function handleBack() {
        if (nodeDeploy.data) return
        if (picker.back()) return
        setTemplate(undefined)
        setCommand(undefined)
    }

    // NOTE: Deploy stays clickable while a field is missing - clicking it is
    // what asks for the errors to be shown
    function handleAction() {
        setSubmitted(true)
        if (!command || !isReady()) return
        nodeDeploy.mutate({
            plugin,
            cluster,
            name,
            connection,
            command: command.command,
            postScript: command.postScript,
            keeperPort: Number(keeperPort) || undefined,
            dbPort: Number(dbPort) || undefined,
            vaults: {databaseId, sshKeyId: sshKeyId ?? ""},
        })
    }

    // NOTE: the ports come from the plugin's own defaults, and the name from
    // the node this dialog belongs to - it is what the platform addresses its
    // container by
    function handleEffectSpec() {
        if (!spec.data) return
        setKeeperPort(spec.data.keeperPort?.toString() ?? "")
        setDbPort(spec.data.dbPort.toString())
        setName(prev => prev || node)
    }

    // NOTE: a template's command is copied in once; from then on the edit is
    // this node's own and never touches the template
    function handleEffectTemplate() {
        setSubmitted(false)
        if (!template) return
        setCommand(template.commands[0])
    }
}
