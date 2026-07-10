import {Edit, Preview, RocketLaunch} from "@mui/icons-material"
import {Box, Button, Checkbox, TextField, ToggleButton, ToggleButtonGroup, Tooltip} from "@mui/material"
import {useCallback, useEffect, useMemo, useState} from "react"

import {useRouterClusterDeploy} from "../../../../features/cluster/api/ClusterHook"
import {Options as ClusterOptions} from "../../../../features/cluster/api/ClusterType"
import {Feature} from "../../../../features/Feature"
import {ManageAccess} from "../../../../features/management/component/ManageAccess"
import {useRouterNodePlatformDeployOptions} from "../../../../features/node/api/NodeHook"
import {KeeperPlugin, PlatformDeployOptions} from "../../../../features/node/api/NodeType"
import {DbPlugin} from "../../../../features/query/api/QueryType"
import {VaultType} from "../../../../features/vault/api/VaultType"
import {Code} from "../../../../shared/component/box/Code"
import {ErrorSmart} from "../../../../shared/component/box/ErrorSmart"
import {List} from "../../../../shared/component/box/List"
import {ListItem} from "../../../../shared/component/box/ListItem"
import {Logs} from "../../../../shared/component/box/Logs"
import {SubContentBox} from "../../../../shared/component/box/SubContentBox"
import {TitledBox} from "../../../../shared/component/box/TitledBox"
import {DialogButton} from "../../../../shared/component/button/DialogButton"
import {SkeletonGroup} from "../../../../shared/component/progress/SkeletonGroup"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {
    getInterpolatedImageOptions,
    getNodeConfigs,
    InterpolatedOptionsKeys,
    VaultOptions,
} from "../../../../shared/helper/HelperUtils"
import {Options} from "../../../widgets/options/Options"
import {OptionsVault} from "../../../widgets/options/OptionsVault"
import {ListNodeInput} from "./ListNodeInput"

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

export function ListDeployCluster(props: Props) {
    const {keeper, database, withLabel = false} = props
    const [cluster, setCluster] = useState("")
    const [image, setImage] = useState<PlatformDeployOptions>()
    const [imagePlugin, setImagePlugin] = useState<KeeperPlugin>()
    const [options, setOptions] = useState<ClusterOptions>(InitialRequest(keeper, database))
    const [imageOptions, setImageOptions] = useState<{[node: string]: string}>({})
    const [nodes, setNodes] = useState<string[]>([])
    const [ssh, setSsh] = useState<"new" | "vault">("vault")
    const [db, setDb] = useState<"new" | "vault">("vault")
    const [sshCred, setSshCred] = useState({username: "", password: ""})
    const [dbCred, setDbCred] = useState({username: "", password: ""})
    const [dcs, setDcs] = useState("")
    const [preview, setPreview] = useState(true)
    const [dev, setDev] = useState(false)
    const [parallel, setParallel] = useState(false)
    const [response, setResponse] = useState<string[] | undefined>(undefined)

    const deploy = useRouterClusterDeploy(setResponse)
    const deployOptions = useRouterNodePlatformDeployOptions(options.plugins.keeper)

    useEffect(handleEffectDeployOptions, [deployOptions.data, imagePlugin, options.plugins.keeper, dev])

    const imageOptionsStr = useMemo(() => (dev ? image?.optionsSingleHost : image?.options) ?? "", [dev, image])

    const handleImageUpdate = useCallback(handleCallImageUpdate, [])
    const handleVaultUpdate = useCallback(handleCallVaultUpdate, [])
    const handleOptionsUpdate = useCallback(handleCallOptionsUpdate, [])
    const handleEnvUpdates = useCallback(handleCallEnvUpdates, [])
    const handleSingleHostUpdate = useCallback(handleCallSingleHostUpdate, [])
    const handleNodesUpdate = useCallback(handleCallNodesUpdate, [imageOptionsStr])

    const imageOptionEntries = useMemo(handleMemoImageOptionEntries, [imageOptions])
    const imageInterpolatedOptions = useMemo(
        handleMemoImageInterpolatedOptions,
        [imageOptionEntries, cluster, dbCred, dcs, preview]
    )

    return (
        <ManageAccess feature={Feature.ManageClusterCreate}>
            <DialogButton
                title={"DEPLOY CLUSTER"}
                renderActions={renderActions()}
                icon={<RocketLaunch/>}
                variant={withLabel ? "button_label" : "icon"}
                label={"Deploy"}
                back={!!response}
                onBackClick={() => setResponse(undefined)}
            >
                {response ? <Logs logs={response} height={600} auto={false}/> : renderBody()}
            </DialogButton>
        </ManageAccess>
    )

    function renderBody() {
        if (deployOptions.isError) return <ErrorSmart error={deployOptions.error}/>
        if (deployOptions.isPending || !image || imagePlugin !== options.plugins.keeper) return <SkeletonGroup count={3}/>
        return (
            <Box sx={[SX.subContent, {gap: 1}]}>
                {renderMandatoryFields()}
                {renderImageOptions()}
                {renderClusterOptions()}
            </Box>
        )
    }

    function renderActions() {
        return (
            <Button fullWidth={true} loading={deploy.isPending} onClick={handleDeploy} disabled={!cluster || !image || !!response}>
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
        if (image?.defaultValues["dcs"]) return
        return (
            <TitledBox title={"Mandatory Options"} island={true}>
                <TextField
                    fullWidth={true}
                    size={"small"}
                    label={"DCS (etcd, zookeper, etc)"}
                    helperText={"Example: etcd1:2379, etcd3:2379, etcd3:2379"}
                    value={dcs}
                    onChange={(e) => setDcs(e.target.value)}
                />
            </TitledBox>
        )
    }

    function renderDbInputs() {
        return (
            <TitledBox title={"Database Credentials"} renderActions={renderDbInputActions()} island={true}>
                {db === "new" ? (
                    <Box sx={SX.between}>
                        <TextField
                            fullWidth
                            size={"small"}
                            label={"Username"}
                            value={dbCred.username}
                            disabled={!!image?.defaultValues["username"]}
                            onChange={v => setDbCred({...dbCred, username: v.target.value})}
                        />
                        <TextField
                            fullWidth
                            size={"small"}
                            type={"password"}
                            label={"Password"}
                            value={dbCred.password}
                            disabled={!!image?.defaultValues["password"]}
                            onChange={v => setDbCred({...dbCred, password: v.target.value})}
                        />
                    </Box>
                ) : (
                    <OptionsVault
                        type={VaultType.DATABASE_PASSWORD}
                        selected={options.vaults.databaseId}
                        onUpdate={handleVaultUpdate}
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
                            value={image?.uri ?? ""}
                            onChange={v => handleImageUpdate(v.target.value)}
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
                            {InterpolatedOptionsKeys.map(k => (
                                <Code key={k} sx={{fontSize: "11px"}}>{`{{${k}}}`}</Code>
                            ))}
                        </Box>
                    </Box>
                    {imageOptionEntries.length === 0 ? (
                        <Box sx={SX.note}>Start by adding nodes – image options will appear here</Box>
                    ) : imageOptionEntries.map(([nodeFull]) => (
                        <TextField
                            key={nodeFull}
                            fullWidth={true}
                            multiline={true}
                            disabled={preview}
                            size={"small"}
                            label={<Box>Node <Code>{nodeFull}</Code> options</Box>}
                            value={imageInterpolatedOptions[nodeFull]}
                            onChange={v => handleEnvUpdates(nodeFull, v.target.value)}
                        />
                    ))}
                </Box>
            </SubContentBox>
        )
    }

    function renderClusterOptions() {
        return (
            <SubContentBox label={"Cluster Options"} island={true}>
                <Options options={options} onUpdate={handleOptionsUpdate}/>
            </SubContentBox>
        )
    }

    function handleMemoImageOptionEntries() {
        return Object.entries(imageOptions)
    }

    function handleMemoImageInterpolatedOptions() {
        if (!preview) return Object.fromEntries(imageOptionEntries)
        return Object.fromEntries(imageOptionEntries.map(([nodeFull, opt]) => {
            const [host, keeperPort, dbPort] = nodeFull.split(":")
            return [nodeFull, getInterpolatedImageOptions(opt, {
                cluster: cluster, host: host, keeperPort: Number(keeperPort), dcs,
                dbPort: Number(dbPort), dbUser: dbCred.username, dbPass: dbCred.password,
            })]
        }))
    }

    function handleCallEnvUpdates(node: string, opt: string) {
        setImageOptions(prev => ({...prev, [node]: opt}))
    }

    function handleCallNodesUpdate(nodes: string[]) {
        setNodes(nodes)
        setImageOptions(prev => Object.fromEntries(
            nodes.filter(n => !!n)
                .map(node => prev[node] !== undefined ? [node, prev[node]] : [node, imageOptionsStr])
        ))
    }

    function handleCallSingleHostUpdate(checked: boolean) {
        setDev(checked)
    }

    function handleCallImageUpdate(uri: string) {
        setImage(prev => prev && ({...prev, uri}))
    }

    function handleCallOptionsUpdate(opt: ClusterOptions) {
        setOptions(prev => ({...prev, ...opt}))
    }

    function handleCallVaultUpdate(t: VaultType, s?: string) {
        setOptions(prev => ({...prev, vaults: {...prev.vaults, [VaultOptions[t].key]: s}}))
    }

    function handleDeploy() {
        if (!image) return
        const nodeConfigs = getNodeConfigs(nodes)
        const imageOptionsSmallKeys = Object.fromEntries(
            Object.entries(imageOptions).map(([nodeFull, opt]) => {
                const [host, keeperPort] = nodeFull.split(":")
                const node = keeperPort ? `${host}:${keeperPort}` : host
                return [node, opt]
            })
        )
        deploy.mutate({
            uri: image.uri,
            parallel: parallel,
            nodeRawImageOptions: imageOptionsSmallKeys,
            nodeConfig: nodeConfigs,
            commonConfig: {
                cluster, dcs,
                dbUser: dbCred.username, dbPass: dbCred.password,
                sshUser: sshCred.username, sshPass: sshCred.password,
            },
            clusterOptions: options,
        })
    }

    // NOTE: applies fetched defaults once per selected keeper plugin, so user
    // edits are kept until the plugin changes and per-node options are reseeded
    function handleEffectDeployOptions() {
        const data = deployOptions.data
        if (!data || imagePlugin === options.plugins.keeper) return
        setImagePlugin(options.plugins.keeper)
        setImage(data)
        setDbCred({username: data.defaultValues["username"] ?? "", password: data.defaultValues["password"] ?? ""})
        const template = dev ? data.optionsSingleHost : data.options
        setImageOptions(prev => Object.fromEntries(Object.keys(prev).map(node => [node, template])))
    }
}
