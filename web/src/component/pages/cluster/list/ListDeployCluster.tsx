import {Edit, Preview} from "@mui/icons-material"
import {
    Box, Button, Dialog, DialogActions, DialogTitle,
    TextField, ToggleButton, ToggleButtonGroup, Tooltip,
} from "@mui/material"
import {useCallback, useMemo, useState} from "react"

import {useRouterClusterDeploy} from "../../../../api/cluster/hook"
import {Options as ClusterOptions} from "../../../../api/cluster/type"
import {Plugin as DbPlugin} from "../../../../api/database/type"
import {Feature} from "../../../../api/feature"
import {Plugin as KeeperPlugin} from "../../../../api/keeper/type"
import {VaultType} from "../../../../api/vault/type"
import {SxPropsMap} from "../../../../app/type"
import {getInterpolatedImageOptions, getNodeConfigs, InterpolatedOptionsKeys, VaultOptions} from "../../../../app/utils"
import scroll from "../../../../style/scroll.module.css"
import {Code} from "../../../view/box/Code"
import {SubContentBox} from "../../../view/box/SubContentBox"
import {TitledBox} from "../../../view/box/TitledBox"
import {DeployIconButton} from "../../../view/button/IconButtons"
import {DynamicInputs} from "../../../view/input/DynamicInputs"
import {Access} from "../../../widgets/access/Access"
import {Options} from "../../../widgets/options/Options"
import {OptionsVault} from "../../../widgets/options/OptionsVault"

const SX: SxPropsMap = {
    dialog: {minWidth: "1010px"},
    content: {width: "590px", display: "flex", flexDirection: "column", gap: 1, padding: "5px 16px 5px 24px ", overflowY: "scroll"},
    center: {display: "flex", justifyContent: "center", gap: 3},
    note: {display: "flex", justifyContent: "center", alignItems: "center", color: "text.disabled", fontSize: 12, flexWrap: "wrap", gap: 0.5},
    between: {display: "flex", justifyContent: "space-between", alignItems: "center", gap: 1},
    subContent: {display: "flex", flexDirection: "column"},
    toggleButton: {padding: "0px 10px"},
    input: {height: "40px"},
    logs: {colorScheme: "dark", fontSize: "13px"},
    row: {"&:hover": {color: "primary.main"}},
}

const INITIAL_OPTIONS: ClusterOptions = {
    certs: {}, vaults: {}, tags: [],
    plugins: {database: DbPlugin.POSTGRES, keeper: KeeperPlugin.PATRONI},
    tls: {keeper: false, database: false},
}
const DEFAULT_IMAGES: {[key in KeeperPlugin]: {uri: string, optionStr: string, defaultValues: {[key: string]: string}}} = {
    [KeeperPlugin.PATRONI]: {
        uri: "ghcr.io/zalando/spilo-18:4.1-p2",
        defaultValues: {username: "postgres"},
        optionStr: `
          --name {{host}}
          --hostname {{host}}
          --restart unless-stopped
          -p {{keeperPort}}:{{keeperPort}}
          -p {{dbPort}}:{{dbPort}}
          -v /data/postgres:/home/postgres/pgdata
          -e SCOPE="{{cluster}}"
          -e PATRONI_NAME="{{host}}"
          -e ETCD3_HOSTS="{{dcs}}"
          -e PGPORT={{dbPort}}
          -e APIPORT={{keeperPort}}
          -e PGPASSWORD_SUPERUSER="{{password}}"
          -e RESTAPI_CONNECT_ADDRESS="{{host}}:{{keeperPort}}"
          -e SPILO_CONFIGURATION='{"postgresql":{"connect_address":"{{host}}:{{dbPort}}"},"bootstrap":{"dcs":{"primary_start_timeout":999}}}'
        `.replace(/\s{2,}/g, "\n").trim(),
    },
    [KeeperPlugin.POSTGRES]: {
        uri: "postgres:18",
        defaultValues: {},
        optionStr: `
          --name {{host}}
          --hostname {{host}}
          --restart unless-stopped
          -p {{dbPort}}:{{dbPort}}
          -v /data/postgres:/var/lib/postgresql/data
          -e PGPORT="{{dbPort}}"
          -e POSTGRES_USER="{{username}}"    
          -e POSTGRES_PASSWORD="{{password}}"
        `.replace(/\s{2,}/g, "\n").trim(),
    }
}

type Props = {
    size?: number,
    keeper: KeeperPlugin,
}

export function ListDeployCluster(props: Props) {
    const {size, keeper} = props
    const [cluster, setCluster] = useState("")
    const [image, setImage] = useState(DEFAULT_IMAGES[keeper])
    const [options, setOptions] = useState(INITIAL_OPTIONS)
    const [imageOptions, setImageOptions] = useState<{[node: string]: string}>({})
    const [nodes, setNodes] = useState<string[]>([])
    const [ssh, setSsh] = useState<"new" | "vault">("vault")
    const [db, setDb] = useState<"new" | "vault">("vault")
    const [sshCred, setSshCred] = useState({username: "", password: ""})
    const [dbCred, setDbCred] = useState({username: image.defaultValues["username"] ?? "", password: image.defaultValues["password"] ?? ""})
    const [dcs, setDcs] = useState("")
    const [open, setOpen] = useState(false)
    const [preview, setPreview] = useState(true)
    const [response, setResponse] = useState<string[] | undefined>(undefined)

    const {mutate, isPending} = useRouterClusterDeploy(setResponse)

    const handleImageUpdate = useCallback(handleCallImageUpdate, [])
    const handleVaultUpdate = useCallback(handleCallVaultUpdate, [])
    const handleOptionsUpdate = useCallback(handleCallOptionsUpdate, [])
    const handleEnvUpdates = useCallback(handleCallEnvUpdates, [])
    const handleNodesUpdate = useCallback(handleCallNodesUpdate, [image.optionStr])

    const imageOptionEntries = useMemo(handleMemoImageOptionEntries, [imageOptions])
    const imageInterpolatedOptions = useMemo(
        handleMemoImageInterpolatedOptions,
        [imageOptionEntries, cluster, dbCred, dcs, sshCred, preview]
    )

    return (
        <Access feature={Feature.ManageClusterCreate}>
            <DeployIconButton tooltip={"Deploy Cluster"} size={size} onClick={() => setOpen(!open)}/>
            <Dialog sx={SX.dialog} open={open} onClose={() => setOpen(false)}>
                <DialogTitle sx={SX.center}>Deploy Cluster</DialogTitle>
                <Box sx={SX.content}>
                    {response ? (
                        <Box sx={SX.logs} className={scroll.small}>
                            {response.map((log, i) => (
                                <Box key={i} sx={SX.row}>{log}</Box>
                            ))}
                        </Box>
                    ) : (
                        <Box sx={SX.subContent} gap={1}>
                            {renderMandatoryFields()}
                            {renderDockerImage()}
                            {renderClusterOptions()}
                        </Box>
                    )}
                </Box>
                <DialogActions sx={SX.center}>
                    {response ? (
                        <>
                            <Button color={"inherit"} onClick={() => setResponse(undefined)}>Back</Button>
                            <Button loading={isPending} onClick={() => setOpen(false)} disabled={!cluster}>Ok</Button>
                        </>
                    ) : (
                        <>
                            <Button color={"inherit"} onClick={() => setOpen(false)}>Cancel</Button>
                            <Button loading={isPending} onClick={handleDeploy} disabled={!cluster}>Deploy</Button>
                        </>
                    )}
                </DialogActions>
            </Dialog>
        </Access>
    )

    function renderMandatoryFields() {
        return (
            <Box sx={SX.subContent} gap={1}>
                <TextField
                    fullWidth={true}
                    size={"small"}
                    label={"Cluster Name"}
                    value={cluster}
                    onChange={(e) => setCluster(e.target.value)}
                />
                <DynamicInputs
                    InputProps={SX.input}
                    minLength={3}
                    inputs={nodes}
                    onChange={handleNodesUpdate}
                    editable={true}
                    placeholder={"Node "}
                />
                {renderSshInputs()}
                {renderDbInputs()}
                {renderImageOptions()}
            </Box>
        )
    }

    function renderImageOptions() {
        if (image.defaultValues["dcs"]) return
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
                            disabled={!!image.defaultValues["username"]}
                            onChange={v => setDbCred({...dbCred, username: v.target.value})}
                        />
                        <TextField
                            fullWidth
                            size={"small"}
                            type={"password"}
                            label={"Password"}
                            value={dbCred.password}
                            disabled={!!image.defaultValues["password"]}
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

    function renderDockerImage() {
        return (
            <SubContentBox label={"Docker Options"} island={true}>
                <Box sx={SX.subContent} gap={2}>
                    <Box sx={SX.between}>
                        <TextField
                            fullWidth
                            size={"small"}
                            label={"Image"}
                            value={image.uri}
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
                        <Box>Use interpolated options to automatically populate values -</Box>
                        {InterpolatedOptionsKeys.map(k => (
                            <Code key={k} sx={{fontSize: "11px"}}>{`{{${k}}}`}</Code>
                        ))}
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
                cluster: cluster, host, keeperPort: Number(keeperPort), dbPort: Number(dbPort), dcs,
                dbUser: dbCred.username, dbPass: dbCred.password, sshPass: sshCred.password, sshUser: sshCred.username
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
                .map(node => prev[node] !== undefined ? [node, prev[node]] : [node, image.optionStr])
        ))
    }

    function handleCallImageUpdate(uri: string) {
        setImage(prev => ({...prev, uri}))
    }

    function handleCallOptionsUpdate(opt: ClusterOptions) {
        setOptions(prev => ({...prev, ...opt}))
    }

    function handleCallVaultUpdate(t: VaultType, s?: string) {
        setOptions(prev => ({...prev, vaults: {...prev.vaults, [VaultOptions[t].key]: s}}))
    }

    function handleDeploy() {
        const nodeConfigs = getNodeConfigs(nodes)
        const imageOptionsSmallKeys = Object.fromEntries(
            Object.entries(imageOptions).map(([nodeFull, opt]) => {
                const [host, keeperPort] = nodeFull.split(":")
                const node = keeperPort ? `${host}:${keeperPort}` : host
                return [node, opt]
            })
        )
        mutate({
            uri: image.uri,
            nodeRawOptions: imageOptionsSmallKeys,
            nodeConfig: nodeConfigs,
            commonConfig: {
                cluster, dcs,
                dbUser: dbCred.username, dbPass: dbCred.password,
                sshUser: sshCred.username, sshPass: sshCred.password,
            },
            clusterOptions: options,
        })
    }
}