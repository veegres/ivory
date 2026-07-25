import {Edit, Preview, RocketLaunch} from "@mui/icons-material"
import {Box, Button, Checkbox, TextField, ToggleButton, ToggleButtonGroup, Tooltip} from "@mui/material"
import {useCallback, useEffect, useMemo, useState} from "react"

import {ListNodeInput} from "../../../core/pages/cluster/list/ListNodeInput"
import {Options} from "../../../core/widgets/options/Options"
import {OptionsVault} from "../../../core/widgets/options/OptionsVault"
import {Code} from "../../../shared/component/box/Code"
import {ErrorSmart} from "../../../shared/component/box/ErrorSmart"
import {List} from "../../../shared/component/box/List"
import {ListItem} from "../../../shared/component/box/ListItem"
import {Logs} from "../../../shared/component/box/Logs"
import {SubContentBox} from "../../../shared/component/box/SubContentBox"
import {TitledBox} from "../../../shared/component/box/TitledBox"
import {WarningList} from "../../../shared/component/box/WarningList"
import {DialogButton} from "../../../shared/component/button/DialogButton"
import {TypedField} from "../../../shared/component/input/TypedField"
import {SkeletonGroup} from "../../../shared/component/progress/SkeletonGroup"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {getDeployPlaceholderKeys, getNodeConfigs, NodeInputFormat, VaultOptions} from "../../../shared/helper/HelperUtils"
import {useDebounce} from "../../../shared/hook/Debounce"
import {Feature} from "../../Feature"
import {ManageAccess} from "../../management/component/ManageAccess"
import {useKeeperDeployForm, useRouterNodeKeeperDeployPlan} from "../../node/api/NodeHook"
import {DeployFieldResponse, InterpolationVar, KeeperDeployPlanRequest, KeeperPlugin} from "../../node/api/NodeType"
import {DbPlugin} from "../../query/api/QueryType"
import {VaultType} from "../../vault/api/VaultType"
import {useRouterClusterDeploy} from "../api/ClusterHook"
import {Options as ClusterOptions} from "../api/ClusterType"

const SX: SxPropsMap = {
    note: {
        display: "flex", justifyContent: "center", alignItems: "center",
        color: "text.disabled", fontSize: 12, flexWrap: "wrap", gap: 0.5,
    },
    between: {display: "flex", justifyContent: "space-between", alignItems: "center", gap: 1},
    subContent: {display: "flex", flexDirection: "column"},
    toggleButton: {padding: "0px 10px"},
    input: {height: "40px"},
    logs: {colorScheme: "dark", fontSize: "13px"},
    row: {"&:hover": {color: "primary.main"}},
}

const InitialRequest = (keeper: KeeperPlugin, database: DbPlugin) => ({
    certs: {}, vaults: {}, tags: [],
    plugins: {database, keeper},
    tls: {keeper: false, database: false},
}) as ClusterOptions

type Props = {
    keeper: KeeperPlugin,
    database: DbPlugin,
    withLabel?: boolean,
}

export function ClusterDeploy(props: Props) {
    const {keeper, database, withLabel = false} = props
    const [cluster, setCluster] = useState("")
    const [options, setOptions] = useState<ClusterOptions>(InitialRequest(keeper, database))
    const [overrides, setOverrides] = useState<{[node: string]: string}>({})
    const [nodes, setNodes] = useState<string[]>([])
    const [ssh, setSsh] = useState<"new" | "vault">("vault")
    const [db, setDb] = useState<"new" | "vault">("vault")
    const [sshCred, setSshCred] = useState({username: "", password: ""})
    const [dbCred, setDbCred] = useState({username: "", password: ""})
    const [dev, setDev] = useState(false)
    const [parallel, setParallel] = useState(false)
    const [response, setResponse] = useState<string[] | undefined>(undefined)

    const deploy = useRouterClusterDeploy(setResponse)
    const {
        deploySpec, image, imageUri, setImageUri, ready, preview, setPreview, inputs, updateInput,
        withKeeperPort, withDbCredentials, mandatoryFields, autoFields,
    } = useKeeperDeployForm(options.plugins.keeper)

    useEffect(handleEffectPluginProps, [keeper, database])
    useEffect(handleEffectImage, [image])

    const fields = image?.fields

    const handleVaultUpdate = useCallback(handleCallVaultUpdate, [])
    const handleOptionsUpdate = useCallback(handleCallOptionsUpdate, [])
    const handleNodesUpdate = useCallback(handleCallNodesUpdate, [])
    const handleSingleHostUpdate = useCallback(handleCallSingleHostUpdate, [])

    const nodeFormat = useMemo(handleMemoNodeFormat, [fields, withKeeperPort])
    const placeholderKeys = getDeployPlaceholderKeys(fields, withKeeperPort, withDbCredentials)
    const planRequest = useMemo(
        handleMemoPlanRequest,
        // NOTE: getValues is intentionally not a dependency, its input
        // (inputs) is listed instead
        // eslint-disable-next-line react-hooks/exhaustive-deps
        [image, imageUri, options.plugins.keeper, cluster, dev, nodes, nodeFormat, inputs, overrides, fields]
    )
    const plan = useRouterNodeKeeperDeployPlan(useDebounce(planRequest, 300))
    // NOTE: the plan query keeps the previous response while a new one is
    // fetching (and when it is disabled), so removed nodes are filtered out
    // by the current request instead of waiting for the next response
    const planNodes = useMemo(handleMemoPlanNodes, [plan.data, planRequest])
    const planWarnings = planRequest ? plan.data?.warnings ?? [] : []
    const planValues = planRequest ? plan.data?.values ?? {} : {}

    return (
        <ManageAccess feature={Feature.ManageClusterCreate}>
            <DialogButton
                title={"DEPLOY CLUSTER"}
                renderActions={renderActions()}
                icon={<RocketLaunch/>}
                variant={withLabel ? "button_label" : "button"}
                label={"Deploy"}
                back={!!response}
                onBackClick={() => setResponse(undefined)}
            >
                {response ? <Logs logs={response} height={600} auto={false}/> : renderBody()}
            </DialogButton>
        </ManageAccess>
    )

    function renderBody() {
        if (deploySpec.isError) return <ErrorSmart error={deploySpec.error}/>
        if (deploySpec.isPending || !ready) return <SkeletonGroup count={3}/>
        return (
            <Box sx={[SX.subContent, {gap: 1}]}>
                {renderMandatoryFields()}
                {renderImageOptions()}
                {renderClusterOptions()}
            </Box>
        )
    }

    function renderActions() {
        const planReady = !!planRequest && !!plan.data && planWarnings.length === 0
        return (
            <Button fullWidth={true} loading={deploy.isPending} onClick={handleDeploy} disabled={!cluster || !planReady || !!response}>
                Deploy
            </Button>
        )
    }

    function renderMandatoryFields() {
        return (
            <Box sx={[SX.subContent, {gap: 1}]}>
                <List>
                    <ListItem
                        title={"Parallel Deployment"}
                        description={`
                            Some services, such as Patroni, require sequential deployment to establish a cluster
                            correctly and allow all nodes to connect to each other. Be cautious.
                        `}
                        button={<Checkbox
                            size={"small"}
                            checked={parallel}
                            color={"default"}
                            onChange={(_, c) => setParallel(c)}
                        />}
                    />
                    <ListItem
                        title={"Single-host mode"}
                        description={`
                            Single-VM setup: services communicate via localhost. Local DNS resolution may
                            require additional /etc/hosts entries. Designed for development and testing.
                        `}
                        button={<Checkbox
                            size={"small"}
                            checked={dev}
                            color={"default"}
                            onChange={(_, c) => handleSingleHostUpdate(c)}
                        />}
                    />
                </List>
                <TitledBox title={"Cluster"} island={true}>
                    <Box sx={[SX.subContent, {gap: 1}]}>
                        <TextField
                            fullWidth={true}
                            size={"small"}
                            label={"Name"}
                            value={cluster}
                            onChange={(e) => setCluster(e.target.value)}
                        />
                        <ListNodeInput
                            InputProps={SX.input}
                            minLength={4}
                            inputs={nodes}
                            onChange={handleNodesUpdate}
                            editable={true}
                            withKeeperPort={withKeeperPort}
                        />
                    </Box>
                </TitledBox>
                {renderSshInputs()}
                {renderDbInputs()}
                {renderMandatoryOptions()}
            </Box>
        )
    }

    function renderMandatoryOptions() {
        if (mandatoryFields.length === 0) return
        return (
            <TitledBox title={"Mandatory Options"} island={true}>
                <Box sx={[SX.subContent, {gap: 1}]}>
                    {mandatoryFields.map(renderField)}
                </Box>
            </TitledBox>
        )
    }

    function renderField(field: DeployFieldResponse) {
        return (
            <TypedField
                key={field.name}
                label={field.label}
                example={field.example}
                type={field.type}
                value={inputs[field.name] ?? planValues[field.name] ?? field.default ?? ""}
                disabled={autoFields.includes(field) && preview}
                onChange={(v) => updateInput(field.name, v)}
            />
        )
    }

    function renderDbInputs() {
        if (!withDbCredentials) return
        return (
            <TitledBox title={"Database Credentials"} renderActions={renderDbInputActions()} island={true}>
                {db === "new" ? (
                    <Box sx={SX.between}>
                        <TextField
                            fullWidth
                            size={"small"}
                            label={"Username"}
                            value={dbCred.username}
                            disabled={!!fields?.defaults[InterpolationVar.DbUser]}
                            onChange={v => setDbCred({...dbCred, username: v.target.value})}
                        />
                        <TextField
                            fullWidth
                            size={"small"}
                            type={"password"}
                            label={"Password"}
                            value={dbCred.password}
                            onChange={v => setDbCred({...dbCred, password: v.target.value})}
                        />
                    </Box>
                ) : (
                    <OptionsVault
                        type={VaultType.DATABASE_PASSWORD}
                        selected={options.vaults.databaseId}
                        onUpdate={handleVaultUpdate}
                        username={fields?.defaults[InterpolationVar.DbUser] || undefined}
                    />
                )}
            </TitledBox>
        )
    }

    function renderDbInputActions() {
        return (
            <ToggleButtonGroup size={"small"} exclusive={true} value={db} onChange={(_, v) => setDb(v)}>
                <ToggleButton sx={SX.toggleButton} value={"new"}>NEW</ToggleButton>
                <ToggleButton sx={SX.toggleButton} value={"vault"}>VAULT</ToggleButton>
            </ToggleButtonGroup>
        )
    }

    function renderSshInputs() {
        return (
            <TitledBox title={"SSH Credentials"} renderActions={renderSshInputActions()} island={true}>
                {ssh === "new" ? (
                    <Box sx={SX.between}>
                        <TextField
                            fullWidth
                            size={"small"}
                            label={"Username"}
                            value={sshCred.username}
                            onChange={v => setSshCred({...sshCred, username: v.target.value})}
                        />
                        <TextField
                            fullWidth
                            size={"small"}
                            label={"Password"}
                            type={"password"}
                            value={sshCred.password}
                            onChange={v => setSshCred({...sshCred, password: v.target.value})}
                        />
                    </Box>
                ) : (
                    <OptionsVault
                        type={VaultType.SSH_KEY}
                        selected={options.vaults.sshKeyId}
                        onUpdate={handleVaultUpdate}
                    />
                )}
            </TitledBox>
        )
    }

    function renderSshInputActions() {
        return (
            <ToggleButtonGroup size={"small"} exclusive={true} value={ssh} onChange={(_, v) => setSsh(v)}>
                <ToggleButton sx={SX.toggleButton} value={"new"}>NEW</ToggleButton>
                <ToggleButton sx={SX.toggleButton} value={"vault"}>VAULT</ToggleButton>
            </ToggleButtonGroup>
        )
    }

    function renderImageOptions() {
        return (
            <SubContentBox label={"Image Options"} island={true}>
                <Box sx={[SX.subContent, {gap: 2}]}>
                    <Box sx={SX.between}>
                        <TextField
                            fullWidth={true}
                            size={"small"}
                            label={"Image"}
                            value={imageUri}
                            onChange={v => setImageUri(v.target.value)}
                        />
                        <ToggleButtonGroup value={preview} exclusive={true} size={"small"} onChange={(_, v) => setPreview(v)}>
                            <Tooltip title={"Preview"} placement={"top"}>
                                <ToggleButton value={true}><Preview/></ToggleButton>
                            </Tooltip>
                            <Tooltip title={"Edit"} placement={"top"}>
                                <ToggleButton value={false}><Edit/></ToggleButton>
                            </Tooltip>
                        </ToggleButtonGroup>
                    </Box>
                    <Box sx={SX.note}>
                        <Box>Use interpolated options to automatically populate values</Box>
                        <Box sx={SX.note}>
                            {placeholderKeys.map(k => (
                                <Code key={k} sx={{fontSize: "11px"}}>{k}</Code>
                            ))}
                        </Box>
                    </Box>
                    <WarningList warnings={planWarnings}/>
                    {autoFields.length > 0 && planNodes.length > 0 && (
                        <Box sx={[SX.subContent, {gap: 1}]}>
                            {autoFields.map(renderField)}
                        </Box>
                    )}
                    {planNodes.length === 0 ? (
                        <Box sx={SX.note}>Start by adding nodes – image options will appear here</Box>
                    ) : planNodes.map(node => {
                        const key = `${node.host}:${node.keeperPort}`
                        return (
                            <TextField
                                key={key}
                                fullWidth={true}
                                multiline={true}
                                disabled={preview}
                                size={"small"}
                                label={<Box>Node <Code>{getNodeLabel(node.host, node.keeperPort, node.dbPort, node.sshPort)}</Code> options</Box>}
                                value={preview ? node.optionsPreview : overrides[key] ?? node.options}
                                onChange={v => handleOverrideUpdate(key, v.target.value)}
                            />
                        )
                    })}
                </Box>
            </SubContentBox>
        )
    }

    function renderClusterOptions() {
        return (
            <SubContentBox label={"Cluster Options"} island={true}>
                <Options options={options} onUpdate={handleOptionsUpdate} disablePlugins={true}/>
            </SubContentBox>
        )
    }

    function handleMemoPlanNodes() {
        if (!planRequest) return []
        const keys = new Set(planRequest.nodes.map(n => `${n.host}:${n.keeperPort}`))
        return (plan.data?.nodes ?? []).filter(n => keys.has(`${n.host}:${n.keeperPort}`))
    }

    function handleMemoNodeFormat(): NodeInputFormat | undefined {
        if (!fields) return undefined
        return {
            withKeeperPort,
            defaults: {
                keeperPort: Number(fields.defaults[InterpolationVar.KeeperPort]) || undefined,
                dbPort: Number(fields.defaults[InterpolationVar.DbPort]) || undefined,
                sshPort: 22,
            },
        }
    }

    function handleMemoPlanRequest(): KeeperDeployPlanRequest | undefined {
        if (!image || !nodeFormat) return undefined
        const activeNodes = nodes.filter(n => !!n)
        if (activeNodes.length === 0) return undefined
        return {
            plugin: options.plugins.keeper,
            cluster,
            singleHost: dev,
            image: imageUri,
            values: getValues(),
            nodes: getNodeConfigs(activeNodes, nodeFormat).map(config => (
                {...config, options: overrides[`${config.host}:${config.keeperPort}`]}
            )),
        }
    }

    function handleOverrideUpdate(node: string, opt: string) {
        setOverrides(prev => ({...prev, [node]: opt}))
    }

    function handleCallNodesUpdate(nodes: string[]) {
        setNodes(nodes)
    }

    function handleCallSingleHostUpdate(checked: boolean) {
        setDev(checked)
        // NOTE: the rendered template differs between the modes, edits made
        // for one mode don't apply to the other
        setOverrides({})
    }

    function handleCallOptionsUpdate(opt: ClusterOptions) {
        setOptions(prev => ({...prev, ...opt}))
    }

    function handleCallVaultUpdate(t: VaultType, s?: string) {
        setOptions(prev => ({...prev, vaults: {...prev.vaults, [VaultOptions[t].key]: s}}))
    }

    function handleDeploy() {
        if (!planRequest) return
        deploy.mutate({
            parallel: parallel,
            singleHost: dev,
            image: imageUri,
            nodes: planRequest.nodes,
            values: getValues(),
            commonConfig: {
                cluster,
                dbUser: dbCred.username, dbPass: dbCred.password,
                sshUser: sshCred.username, sshPass: sshCred.password,
            },
            clusterOptions: options,
        })
    }

    // NOTE: the plugin selectors are disabled inside the dialog, so the plugins
    // can only change through the cluster list filter; the filter also changes
    // the node input format, hence the entered nodes are reset (the mandatory
    // field inputs are reset by useKeeperDeployForm itself once the new
    // plugin's spec loads)
    function handleEffectPluginProps() {
        setOptions(prev => ({...prev, plugins: {keeper, database}}))
        setNodes([])
        setOverrides({})
    }

    // NOTE: applies the fetched default database username once per selected
    // keeper plugin, so a user edit is kept until the plugin changes
    function handleEffectImage() {
        if (!image) return
        setDbCred({username: image.fields.defaults[InterpolationVar.DbUser] ?? "", password: ""})
    }

    // NOTE: credentials are never part of the values — they are resolved
    // from the vault at execution time and stay visible as {{dbUser}}/
    // {{dbPass}} in the preview
    function getValues() {
        return {...inputs}
    }

    function getNodeLabel(host: string, keeperPort: number, dbPort: number, sshPort: number) {
        return withKeeperPort ? `${host}:${keeperPort}:${dbPort}:${sshPort}` : `${host}:${dbPort}:${sshPort}`
    }
}
