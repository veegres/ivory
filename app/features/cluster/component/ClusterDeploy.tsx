import {RocketLaunch} from "@mui/icons-material"
import {Box, Button, Checkbox, TextField, ToggleButton, ToggleButtonGroup} from "@mui/material"
import {useCallback, useEffect, useMemo, useState} from "react"

import {Options} from "../../../core/widgets/options/Options"
import {OptionsVault} from "../../../core/widgets/options/OptionsVault"
import {ErrorSmart} from "../../../shared/component/box/ErrorSmart"
import {Logs} from "../../../shared/component/box/Logs"
import {Note} from "../../../shared/component/box/Note"
import {SubContentBox} from "../../../shared/component/box/SubContentBox"
import {TitledBox} from "../../../shared/component/box/TitledBox"
import {DialogButton} from "../../../shared/component/button/DialogButton"
import {FieldRow} from "../../../shared/component/input/FieldRow"
import {SkeletonGroup} from "../../../shared/component/progress/SkeletonGroup"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {DeployPasswordMask, KeeperPluginOptions, VaultOptions} from "../../../shared/helper/HelperUtils"
import {Template} from "../../deployment/api/DeploymentType"
import {useDeploymentTemplatePicker} from "../../deployment/component/DeploymentTemplatePicker"
import {Feature} from "../../Feature"
import {ManageAccess} from "../../management/component/ManageAccess"
import {useRouterNodeKeeperDeploySpec} from "../../node/api/NodeHook"
import {KeeperPlugin, PlatformPlugin} from "../../node/api/NodeType"
import {DbPlugin} from "../../query/api/QueryType"
import {useRouterVault} from "../../vault/api/VaultHook"
import {VaultType} from "../../vault/api/VaultType"
import {useRouterClusterDeploy} from "../api/ClusterHook"
import {DeployCredentials, DeployNode, Options as ClusterOptions} from "../api/ClusterType"
import {ClusterDeployNode} from "./ClusterDeployNode"

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
    toggleButton: {padding: "0px 10px"},
}

// NOTE: a constant until a second platform exists - it was state with no
// setter, which reads as a choice nobody can make
const platform = PlatformPlugin.DOCKER

const InitialRequest = (keeper: KeeperPlugin, database: DbPlugin) => ({
    certs: {}, vaults: {}, tags: [],
    plugins: {database, keeper},
    tls: {keeper: false, database: false},
}) as ClusterOptions

type Props = {
    keeper: KeeperPlugin,
    database: DbPlugin,
    withLabel?: boolean,
    size?: number,
}

// ClusterDeploy fills in a template: the template already says what to run and
// on how many nodes, so this asks only where each of those nodes lives. The
// node count comes from the template and is not editable here - a different
// size is a different template.
export function ClusterDeploy(props: Props) {
    const {keeper, database, withLabel = false, size} = props
    const [template, setTemplate] = useState<Template>()
    const [nodes, setNodes] = useState<DeployNode[]>([])
    const [cluster, setCluster] = useState("")
    const [options, setOptions] = useState<ClusterOptions>(InitialRequest(keeper, database))
    const [ssh, setSsh] = useState<"new" | "vault">("vault")
    const [db, setDb] = useState<"new" | "vault">("vault")
    const [sshCred, setSshCred] = useState({username: "", password: ""})
    const [dbCred, setDbCred] = useState({username: "", password: ""})
    const [parallel, setParallel] = useState(false)
    const [submitted, setSubmitted] = useState(false)
    const [response, setResponse] = useState<string[] | undefined>(undefined)

    const deploy = useRouterClusterDeploy(setResponse)
    const spec = useRouterNodeKeeperDeploySpec(options.plugins.keeper)
    const picker = useDeploymentTemplatePicker({
        keeper: options.plugins.keeper,
        platform,
        hint: "Pick a template to deploy a cluster, copy one to adjust it, or write a new one",
        onPick: setTemplate,
    })
    // NOTE: shares OptionsVault's query, so resolving the chosen vault's
    // username for the preview costs no extra request
    const dbVaults = useRouterVault(VaultType.DATABASE_PASSWORD)

    useEffect(handleEffectPluginProps, [keeper, database])
    useEffect(handleEffectSpec, [spec.data])
    // NOTE: the node list is seeded from the template exactly once; from then
    // on each node is the user's own, so re-seeding would discard their input
    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(handleEffectTemplate, [template])

    const withKeeperPort = spec.data?.keeperPort !== undefined
    const withDbCredentials = spec.data?.credentials ?? false

    const handleVaultUpdate = useCallback(handleCallVaultUpdate, [])
    const handleOptionsUpdate = useCallback(handleCallOptionsUpdate, [])

    const complete = nodes.length > 0 && nodes.every(n => !!n.name && !!n.host && !!n.command.trim())
    // NOTE: a duplicate name is a conflict the user should see while typing,
    // unlike an unfilled field, which waits for Deploy
    const duplicates = useMemo(handleMemoDuplicates, [nodes])

    return (
        <ManageAccess feature={Feature.ManageClusterCreate}>
            <DialogButton
                title={"DEPLOY CLUSTER"}
                renderActions={!picker.editing && template && renderActions()}
                icon={<RocketLaunch/>}
                variant={withLabel ? "button_label" : "button"}
                label={"Deploy"}
                size={size}
                back={!!response || !!template || picker.editing}
                onBackClick={handleBack}
            >
                {renderContent()}
            </DialogButton>
        </ManageAccess>
    )

    function renderContent() {
        if (response) return <Logs logs={response} height={570} auto={false}/>
        if (spec.isError) return <ErrorSmart error={spec.error}/>
        if (spec.isPending) return <SkeletonGroup count={3}/>
        if (!template) return picker.render()
        return renderBody(template)
    }

    function renderActions() {
        return (
            <Button fullWidth={true} loading={deploy.isPending} onClick={handleDeploy} disabled={!!response}>
                Deploy
            </Button>
        )
    }

    function renderBody(selected: Template) {
        return (
            <Box sx={SX.box}>
                {renderCluster(selected)}
                {renderSshInputs()}
                {renderDbInputs()}
                {renderNodes()}
            </Box>
        )
    }

    function renderCluster(selected: Template) {
        return (
            <TitledBox title={"Cluster"} renderActions={renderTemplateName(selected)} island={true}>
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
                    {renderClusterOptions()}
                </Box>
            </TitledBox>
        )
    }

    function renderTemplateName(selected: Template) {
        return <Box sx={SX.templateName}>{selected.name}</Box>
    }

    function renderParallel() {
        return (
            <Box sx={SX.toggle} onClick={() => setParallel(!parallel)}>
                <Box>
                    <Box>Parallel deployment</Box>
                    <Note>Some keepers, such as Patroni, need their nodes deployed one after another.</Note>
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
                withKeeperPort={withKeeperPort}
                showErrors={submitted}
                duplicate={!!node.name && duplicates.has(node.name)}
                credentials={getCredentials()}
                onChange={(updated) => handleNodeUpdate(index, updated)}
            />
        )
    }

    function renderDbInputs() {
        if (!withDbCredentials) return
        return (
            <TitledBox title={"Database Credentials"} renderActions={renderDbInputActions()} island={true}>
                {db === "new" ? (
                    <FieldRow>
                        <TextField
                            fullWidth
                            size={"small"}
                            label={"Username"}
                            value={dbCred.username}
                            disabled={!!spec.data?.dbUser}
                            error={submitted && db === "new" && !dbCred.username}
                            onChange={v => setDbCred({...dbCred, username: v.target.value})}
                        />
                        <TextField
                            fullWidth
                            size={"small"}
                            type={"password"}
                            label={"Password"}
                            value={dbCred.password}
                            error={submitted && db === "new" && !dbCred.password}
                            onChange={v => setDbCred({...dbCred, password: v.target.value})}
                        />
                    </FieldRow>
                ) : (
                    <OptionsVault
                        type={VaultType.DATABASE_PASSWORD}
                        selected={options.vaults.databaseId}
                        onUpdate={handleVaultUpdate}
                        username={spec.data?.dbUser || undefined}
                        error={submitted && db === "vault" && !options.vaults.databaseId}
                    />
                )}
            </TitledBox>
        )
    }

    function renderDbInputActions() {
        return (
            <ToggleButtonGroup size={"small"} exclusive={true} value={db} onChange={(_, v) => v && setDb(handleModeChange(VaultType.DATABASE_PASSWORD, v))}>
                <ToggleButton sx={SX.toggleButton} value={"new"}>NEW</ToggleButton>
                <ToggleButton sx={SX.toggleButton} value={"vault"}>VAULT</ToggleButton>
            </ToggleButtonGroup>
        )
    }

    function renderSshInputs() {
        return (
            <TitledBox title={"SSH Credentials"} renderActions={renderSshInputActions()} island={true}>
                {ssh === "new" ? (
                    <FieldRow>
                        <TextField
                            fullWidth
                            size={"small"}
                            label={"Username"}
                            value={sshCred.username}
                            error={submitted && ssh === "new" && !sshCred.username}
                            onChange={v => setSshCred({...sshCred, username: v.target.value})}
                        />
                        <TextField
                            fullWidth
                            size={"small"}
                            label={"Password"}
                            type={"password"}
                            value={sshCred.password}
                            error={submitted && ssh === "new" && !sshCred.password}
                            onChange={v => setSshCred({...sshCred, password: v.target.value})}
                        />
                    </FieldRow>
                ) : (
                    <OptionsVault
                        type={VaultType.SSH_KEY}
                        selected={options.vaults.sshKeyId}
                        onUpdate={handleVaultUpdate}
                        error={submitted && ssh === "vault" && !options.vaults.sshKeyId}
                    />
                )}
            </TitledBox>
        )
    }

    function renderSshInputActions() {
        return (
            <ToggleButtonGroup size={"small"} exclusive={true} value={ssh} onChange={(_, v) => v && setSsh(handleModeChange(VaultType.SSH_KEY, v))}>
                <ToggleButton sx={SX.toggleButton} value={"new"}>NEW</ToggleButton>
                <ToggleButton sx={SX.toggleButton} value={"vault"}>VAULT</ToggleButton>
            </ToggleButtonGroup>
        )
    }

    function renderClusterOptions() {
        return (
            <SubContentBox label={"Cluster Options"} dense={true}>
                <Options options={options} onUpdate={handleOptionsUpdate} disablePlugins={true}/>
            </SubContentBox>
        )
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

    // NOTE: the unselected mode's vault id is cleared too - leaving one behind
    // would send both a vault and a password for the same credential
    function handleModeChange(type: VaultType, mode: "new" | "vault") {
        if (mode === "new") handleCallVaultUpdate(type, undefined)
        return mode
    }

    function handleCallVaultUpdate(t: VaultType, s?: string) {
        setOptions(prev => ({...prev, vaults: {...prev.vaults, [VaultOptions[t].key]: s}}))
    }

    function handleBack() {
        if (response) return setResponse(undefined)
        if (picker.back()) return
        setTemplate(undefined)
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
            commonConfig: {cluster, ...getSshConfig(), ...getDbConfig()},
            clusterOptions: options,
        })
    }

    // NOTE: mirrors what cluster.Deploy rejects - ssh credentials are always
    // required, database ones only when the keeper plugin consumes them, which
    // is the same condition that renders their section
    function isReady() {
        return !!cluster && duplicates.size === 0 && !response && complete && isSshReady() && isDbReady()
    }

    // NOTE: the preview may show the username - the form knows it either way -
    // but never the password, only that one exists
    function getCredentials(): DeployCredentials {
        if (!withDbCredentials) return {}
        if (db === "new") {
            return {
                user: dbCred.username || undefined,
                pass: dbCred.password ? DeployPasswordMask : undefined,
            }
        }
        const vault = getDbVault()
        return {user: vault?.username, pass: vault && DeployPasswordMask}
    }

    function getDbVault() {
        const id = options.vaults.databaseId
        return id ? dbVaults.data?.[id] : undefined
    }

    // NOTE: a vault id and a username/password pair are two answers to one
    // question and the server rejects both together, so whichever the user did
    // not choose is left out entirely rather than sent alongside
    function getSshConfig() {
        if (ssh === "vault") return {sshUser: "", sshPass: ""}
        return {sshUser: sshCred.username, sshPass: sshCred.password}
    }

    function getDbConfig() {
        if (!withDbCredentials || db === "vault") return {dbUser: "", dbPass: ""}
        return {dbUser: dbCred.username, dbPass: dbCred.password}
    }

    function isSshReady() {
        if (ssh === "vault") return !!options.vaults.sshKeyId
        return !!sshCred.username && !!sshCred.password
    }

    function isDbReady() {
        if (!withDbCredentials) return true
        if (db === "vault") return !!options.vaults.databaseId
        return !!dbCred.username && !!dbCred.password
    }

    // NOTE: the plugin selectors are disabled inside the dialog, so the plugins
    // can only change through the cluster list filter; the filter also changes
    // the available templates, hence the reset
    function handleEffectPluginProps() {
        setOptions(prev => ({...prev, plugins: {keeper, database}}))
        setTemplate(undefined)
        setNodes([])
    }

    // NOTE: applies the engine-locked database username once per selected
    // keeper plugin, so a user edit is kept until the plugin changes
    function handleEffectSpec() {
        if (!spec.data) return
        setDbCred({username: spec.data.dbUser ?? "", password: ""})
    }

    // NOTE: the template's command count is the node count - a cluster of a
    // different size is a different template, so nodes are never added here
    function handleEffectTemplate() {
        setSubmitted(false)
        if (!template) return setNodes([])
        // NOTE: the shipped templates reference these exact names in their own
        // text (etcd-1, postgres-1, ...), so the prefill makes a copied
        // template deployable without editing its commands
        setNodes(template.commands.map((c, i) => ({
            name: `${KeeperPluginOptions[options.plugins.keeper].name}-${i + 1}`,
            host: "",
            keeperPort: spec.data?.keeperPort,
            dbPort: spec.data?.dbPort,
            sshPort: 22,
            command: c.command,
            postScript: c.postScript,
        })))
    }
}
