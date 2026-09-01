import {Box, Button, Checkbox, TextField} from "@mui/material"
import {useCallback, useMemo, useState} from "react"

import {DialogLogsScreen} from "../../../shared/component/box/DialogLogsScreen"
import {DialogScreen} from "../../../shared/component/box/DialogScreen"
import {Hint} from "../../../shared/component/box/Hint"
import {TitleBox} from "../../../shared/component/box/TitleBox"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {DeployPasswordMask, KeeperPluginOptions, VaultOptions} from "../../../shared/helper/HelperUtils"
import {useDeployVaultCredentials} from "../../deployment/api/DeploymentHook"
import {DeployCredentials, Template} from "../../deployment/api/DeploymentType"
import {KeeperPlugin} from "../../node/api/NodeType"
import {DbPlugin} from "../../query/api/QueryType"
import {VaultType} from "../../vault/api/VaultType"
import {useRouterClusterDeploy} from "../api/ClusterHook"
import {DeployNode, DeployPreviewCredentials, Options as ClusterOptions} from "../api/ClusterType"
import {ClusterDeployCredentials, Credential, CredentialMode} from "./ClusterDeployCredentials"
import {ClusterDeployNode} from "./ClusterDeployNode"
import {ClusterOptionsBox} from "./ClusterOptionsBox"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 2},
    column: {display: "flex", flexDirection: "column", gap: 1},
    // NOTE: bordered like the field above it and like a template row, so a
    // clickable row reads as one of them rather than as loose text; the radius
    // and the hover border are the outlined input's, not the island's, since
    // this sits inline with the fields rather than around them
    toggle: {
        display: "flex", justifyContent: "space-between", alignItems: "center", gap: 1,
        padding: "6px 10px", border: 1, borderColor: "divider", borderRadius: 1,
        cursor: "pointer", userSelect: "none",
        "&:hover": {bgcolor: "action.hover", borderColor: "text.primary"},
    },
    templateName: {fontFamily: "monospace", fontSize: "13px", color: "text.secondary"},
}

const InitialRequest = (keeper: KeeperPlugin, database: DbPlugin) => ({
    certs: {}, vaults: {}, tags: [],
    plugins: {database, keeper},
    tls: {keeper: false, database: false},
}) as ClusterOptions

type Props = {
    keeper: KeeperPlugin,
    database: DbPlugin,
    template: Template,
    logs?: string[],
    onDeployed: (logs: string[]) => void,
}

// ClusterDeployForm fills in a template: the template already says what to run
// and on how many nodes, so this asks only where each of those nodes lives. The
// node count comes from the template and is not editable here - a different
// size is a different template.
export function ClusterDeployForm(props: Props) {
    const {keeper, database, template, logs, onDeployed} = props
    const [nodes, setNodes] = useState<DeployNode[]>(getInitialNodes)
    const [cluster, setCluster] = useState("")
    const [options, setOptions] = useState<ClusterOptions>(InitialRequest(keeper, database))
    const [sshMode, setSshMode] = useState<CredentialMode>("vault")
    // NOTE: a template that names no username for a pair is saying its
    // deployment ends up with no such account, so that section opens switched
    // off rather than demanding an answer the deployment has no use for
    const [keeperMode, setKeeperMode] = useState<CredentialMode>(getInitialMode(template.defaults?.keeperUser))
    const [dbMode, setDbMode] = useState<CredentialMode>(getInitialMode(template.defaults?.dbUser))
    const [sshCred, setSshCred] = useState<Credential>({username: "", password: ""})
    // NOTE: seeded with the usernames the template names, and read-only where
    // it names one: that account is what its commands create, so the deploy
    // only asks for the password to give it
    const [keeperCred, setKeeperCred] = useState<Credential>({username: template.defaults?.keeperUser ?? "", password: ""})
    const [dbCred, setDbCred] = useState<Credential>({username: template.defaults?.dbUser ?? "", password: ""})
    const [parallel, setParallel] = useState(false)
    const [submitted, setSubmitted] = useState(false)

    const deploy = useRouterClusterDeploy(onDeployed)
    const vaultCredentials = useDeployVaultCredentials(options.vaults.keeperId, options.vaults.databaseId)

    const handleVaultUpdate = useCallback(handleCallVaultUpdate, [])
    const handleOptionsUpdate = useCallback(handleCallOptionsUpdate, [])

    const complete = nodes.length > 0 && nodes.every(isNodeComplete)
    // NOTE: a duplicate name is a conflict the user should see while typing,
    // unlike an unfilled field, which waits for Deploy
    const duplicates = useMemo(handleMemoDuplicates, [nodes])

    // NOTE: the response replaces the fields from inside this component rather
    // than from the dialog above it, so going back from the logs of a deploy
    // that failed on its third node returns to the form still filled in
    if (logs) return <DialogLogsScreen logs={logs}/>

    return (
        <DialogScreen renderActions={renderActions()}>
            <Box sx={SX.box}>
                {renderCluster()}
                {renderNodes()}
            </Box>
        </DialogScreen>
    )

    function renderActions() {
        return (
            <Button fullWidth={true} loading={deploy.isPending} onClick={handleDeploy}>
                Deploy
            </Button>
        )
    }

    function renderCluster() {
        return (
            <TitleBox label={"Cluster"} renderActions={renderTemplateName()} island={true} collapsible={false}>
                <Box sx={SX.column}>
                    <TextField
                        fullWidth={true}
                        size={"small"}
                        label={"Name"}
                        value={cluster}
                        error={submitted && !cluster}
                        onChange={(e) => setCluster(e.target.value)}
                    />
                    {renderParallel()}
                    {renderSshCredentials()}
                    {renderKeeperCredentials()}
                    {renderDbCredentials()}
                    {renderClusterOptions()}
                </Box>
            </TitleBox>
        )
    }

    function renderTemplateName() {
        return <Box sx={SX.templateName}>{template.name}</Box>
    }

    function renderParallel() {
        return (
            <Box sx={SX.toggle} onClick={() => setParallel(!parallel)}>
                <Box>
                    <Box>Parallel deployment</Box>
                    <Hint>Some keepers, such as Patroni, need their nodes deployed one after another.</Hint>
                </Box>
                {/* NOTE: the row toggles on click, so the checkbox has to stop
                    its own click bubbling or the two would cancel out */}
                <Checkbox
                    size={"small"}
                    color={"default"}
                    checked={parallel}
                    onClick={(e) => e.stopPropagation()}
                    onChange={(_, c) => setParallel(c)}
                />
            </Box>
        )
    }

    // NOTE: no wrapping box - a node box is a top-level section here exactly as
    // it is in the template editor, so the two screens show a node the same way
    function renderNodes() {
        return nodes.map(renderNode)
    }

    function renderNode(node: DeployNode, index: number) {
        return (
            <ClusterDeployNode
                key={index}
                node={node}
                cluster={cluster}
                showErrors={submitted}
                duplicate={!!node.name && duplicates.has(node.name)}
                credentials={getPreviewCredentials()}
                onChange={(updated) => handleNodeUpdate(index, updated)}
            />
        )
    }

    function renderSshCredentials() {
        return (
            <ClusterDeployCredentials
                title={"SSH Credentials"}
                type={VaultType.SSH_KEY}
                mode={sshMode}
                credential={sshCred}
                vaultId={options.vaults.sshKeyId}
                showErrors={submitted}
                onModeChange={setSshMode}
                onCredentialChange={setSshCred}
                onVaultChange={handleVaultUpdate}
            />
        )
    }

    // NOTE: an engine that is its own keeper renders this and the database
    // section both, and the user points them at the same vault entry
    function renderKeeperCredentials() {
        return (
            <ClusterDeployCredentials
                title={"Keeper Credentials"}
                type={VaultType.KEEPER_PASSWORD}
                mode={keeperMode}
                credential={keeperCred}
                vaultId={options.vaults.keeperId}
                locked={!!template.defaults?.keeperUser}
                optional={true}
                showErrors={submitted}
                onModeChange={setKeeperMode}
                onCredentialChange={setKeeperCred}
                onVaultChange={handleVaultUpdate}
            />
        )
    }

    function renderDbCredentials() {
        return (
            <ClusterDeployCredentials
                title={"Database Credentials"}
                type={VaultType.DATABASE_PASSWORD}
                mode={dbMode}
                credential={dbCred}
                vaultId={options.vaults.databaseId}
                locked={!!template.defaults?.dbUser}
                optional={true}
                showErrors={submitted}
                onModeChange={setDbMode}
                onCredentialChange={setDbCred}
                onVaultChange={handleVaultUpdate}
            />
        )
    }

    function renderClusterOptions() {
        return <ClusterOptionsBox options={options} onUpdate={handleOptionsUpdate}/>
    }

    function handleMemoDuplicates() {
        const names = nodes.map(n => n.name).filter(name => !!name)
        return new Set(names.filter((name, i) => names.indexOf(name) !== i))
    }

    function handleNodeUpdate(index: number, node: DeployNode) {
        setNodes(prev => prev.map((n, i) => i === index ? node : n))
    }

    function handleCallOptionsUpdate(opt: ClusterOptions) {
        setOptions(prev => ({...prev, ...opt}))
    }

    function handleCallVaultUpdate(t: VaultType, s?: string) {
        setOptions(prev => ({...prev, vaults: {...prev.vaults, [VaultOptions[t].key]: s}}))
    }

    // NOTE: Deploy stays clickable while fields are missing - clicking it is
    // what asks for the errors to be shown, and a disabled button could never
    // explain itself
    function handleDeploy() {
        setSubmitted(true)
        if (!isReady()) return
        deploy.mutate({
            parallel,
            nodes,
            commonConfig: {cluster, ...getSshConfig(), ...getKeeperConfig(), ...getDbConfig()},
            clusterOptions: options,
        })
    }

    // NOTE: mirrors what cluster.Deploy rejects - ssh is the one credential
    // that can never be answered with nothing, since it is how Ivory reaches
    // the host at all; the other two are complete the moment the user says
    // they are not needed
    function isReady() {
        return !!cluster && duplicates.size === 0 && complete
            && isCredentialReady(sshMode, sshCred, options.vaults.sshKeyId)
            && isCredentialReady(keeperMode, keeperCred, options.vaults.keeperId)
            && isCredentialReady(dbMode, dbCred, options.vaults.databaseId)
    }

    function isCredentialReady(mode: CredentialMode, credential: Credential, vaultId?: string) {
        if (mode === "none") return true
        if (mode === "vault") return !!vaultId
        return !!credential.username && !!credential.password
    }

    function isNodeComplete(node: DeployNode) {
        return !!node.name && !!node.host && !!node.command.trim()
            && !!node.keeperPort && !!node.dbPort && !!node.sshPort
    }

    // NOTE: the preview may show a username - the form knows it either way -
    // but never a password, only that one exists
    function getPreviewCredentials(): DeployPreviewCredentials {
        return {
            keeper: getCredentialPreview(keeperMode, keeperCred, vaultCredentials.keeper),
            database: getCredentialPreview(dbMode, dbCred, vaultCredentials.database),
        }
    }

    function getCredentialPreview(mode: CredentialMode, credential: Credential, vault: DeployCredentials): DeployCredentials {
        if (mode === "none") return {}
        if (mode === "new") {
            return {
                user: credential.username || undefined,
                pass: credential.password ? DeployPasswordMask : undefined,
            }
        }
        return vault
    }

    // NOTE: a vault id and a username/password pair are two answers to one
    // question and the server rejects both together, so whichever the user did
    // not choose is left out entirely rather than sent alongside
    function getSshConfig() {
        if (sshMode === "vault") return {sshUser: "", sshPass: ""}
        return {sshUser: sshCred.username, sshPass: sshCred.password}
    }

    function getKeeperConfig() {
        if (keeperMode !== "new") return {keeperUser: "", keeperPass: ""}
        return {keeperUser: keeperCred.username, keeperPass: keeperCred.password}
    }

    function getDbConfig() {
        if (dbMode !== "new") return {dbUser: "", dbPass: ""}
        return {dbUser: dbCred.username, dbPass: dbCred.password}
    }

    function getInitialMode(defaultUser?: string): CredentialMode {
        return defaultUser ? "vault" : "none"
    }

    // NOTE: the template's command count is the node count - a cluster of a
    // different size is a different template, so nodes are never added here.
    // A card is filled in from its own command and from nowhere else: the ports
    // a node answers on belong to the command that writes them, so a template
    // that states none leaves the field empty for the user rather than
    // borrowing the engine's, which on a single host would put every node on
    // one port. Host is usually empty for the same reason - a multi-host
    // template can't know the actual machine - but a single-host template
    // names one (its commands all land on the same machine), and that default
    // still stays editable like every other field here.
    function getInitialNodes(): DeployNode[] {
        return template.commands.map((c, i) => ({
            name: c.defaults?.name || `${KeeperPluginOptions[keeper].name}${i + 1}`,
            host: c.defaults?.host || "",
            sshPort: c.defaults?.sshPort,
            keeperPort: c.defaults?.keeperPort,
            dbPort: c.defaults?.dbPort,
            command: c.command,
            postScripts: c.postScripts,
        }))
    }
}
