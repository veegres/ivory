import {Box, Button, Dialog, DialogActions, DialogContent, DialogTitle, TextField} from "@mui/material"
import {useCallback, useState} from "react"

import {Options as ClusterOptions} from "../../../../api/cluster/type"
import {Plugin as DbPlugin} from "../../../../api/database/type"
import {Feature} from "../../../../api/feature"
import {Plugin as KeeperPlugin} from "../../../../api/keeper/type"
import {VaultType} from "../../../../api/vault/type"
import {SxPropsMap} from "../../../../app/type"
import {VaultOptions} from "../../../../app/utils"
import {AlertCentered} from "../../../view/box/AlertCentered"
import {Code} from "../../../view/box/Code"
import {SubContentBox} from "../../../view/box/SubContentBox"
import {DeployIconButton} from "../../../view/button/IconButtons"
import {DynamicInputs} from "../../../view/input/DynamicInputs"
import {Access} from "../../../widgets/access/Access"
import {Options} from "../../../widgets/options/Options"
import {OptionsVault} from "../../../widgets/options/OptionsVault"

const SX: SxPropsMap = {
    dialog: {minWidth: "1010px"},
    content: {width: "590px", display: "flex", flexDirection: "column", gap: 1, padding: "5px 16px 5px 24px ", overflowY: "scroll"},
    center: {display: "flex", justifyContent: "center", gap: 3},
    subContent: {display: "flex", flexDirection: "column", gap: 1},
    inputTitle: {
        display: "flex", alignItems: "center", gap: 0.5, padding: "3px 8px 0px",
        color: "text.secondary", letterSpacing: "1px",
    },
    input: {height: "40px"},
}

const INITIAL_OPTIONS: ClusterOptions = {
    certs: {}, vaults: {}, tags: [],
    plugins: {database: DbPlugin.POSTGRES, keeper: KeeperPlugin.PATRONI},
    tls: {keeper: false, database: false},
}
const DEFAULT_IMAGES: {[key in KeeperPlugin]: {uri: string, volume: string, dbPort: number, keeperPort: number, restart: string, defaultCommonEnv: (password: string) => string[], defaultUniqueEnv: (node: string) => string[]}} = {
    [KeeperPlugin.PATRONI]: {
        uri: "ghcr.io/zalando/spilo-18:4.1-p2",
        volume: "/data/postgres:/home/postgres/pgdata",
        restart: "unless-stopped",
        dbPort: 5432,
        keeperPort: 8008,
        defaultCommonEnv: (password: string) => ([
            "SCOPE=\"pg-cluster\"",
            "ETCD_HOSTS=\"localhost:2379\"",
            `PGPASSWORD_SUPERUSER="${password}"`
        ]),
        defaultUniqueEnv: (node: string) => ([
            `PATRONI_NAME="${node}"`
        ]),
    },
    [KeeperPlugin.POSTGRES]: {
        uri: "postgres:18",
        volume: "/data/postgres:/var/lib/postgresql/data",
        restart: "unless-stopped",
        dbPort: 5432,
        keeperPort: 5432,
        defaultCommonEnv: (password: string) => ([`POSTGRES_PASSWORD="${password}"`]),
        defaultUniqueEnv: (_: string) => ([]),
    }
}

type Props = {
    size?: number,
    keeper: KeeperPlugin,
}

export function ListDeployCluster(props: Props) {
    const {size, keeper} = props
    const [name, setName] = useState("")
    const [image, setImage] = useState(DEFAULT_IMAGES[keeper])
    const [options, setOptions] = useState(INITIAL_OPTIONS)
    const [commonEnv, setCommonEnv] = useState(image.defaultCommonEnv("SUPERUSER_PASSWORD"))
    const [uniqueEnv, setUniqueEnv] = useState<{[node: string]: string[]}>({})
    const [nodes, setNodes] = useState([""])
    const [open, setOpen] = useState(false)

    const handleImageChange = useCallback(handleCallImageUpdate, [])
    const handleVaultUpdate = useCallback(handleCallVaultUpdate, [])
    const handleOptionsUpdate = useCallback(handleCallOptionsUpdate, [])
    const handleEnvUpdates = useCallback(handleCallEnvUpdates, [])
    const handleNodesUpdate = useCallback(handleCallNodesUpdate, [image])

    return (
        <Access feature={Feature.ManageClusterCreate}>
            <DeployIconButton tooltip={"Deploy Cluster"} size={size} onClick={() => setOpen(!open)}/>
            <Dialog sx={SX.dialog} open={open} onClose={() => setOpen(false)}>
                <DialogTitle sx={SX.center}>Deploy Cluster</DialogTitle>
                <DialogContent sx={SX.content}>
                    <AlertCentered text={`
                        You can deploy the cluster from scratch here, just providing list of virtual machines with
                        same ssh credentials or generate ssh key in vaults and choose it here. When you provide
                        ssh credentials Ivory will automatically generate ssh key in vault and add it to authorised
                        keys in your VMs.
                        
                        All is preconfigured for you, but you can always change all configs.
                    `}/>
                    {renderRequiredFields()}
                    {renderDockerImage()}
                    {renderClusterOptions()}
                </DialogContent>
                <DialogActions sx={SX.center}>
                    <Button color={"inherit"} onClick={() => setOpen(false)}>Cancel</Button>
                    <Button loading={false} onClick={() => void 0} disabled={!name}>Deploy</Button>
                </DialogActions>
            </Dialog>
        </Access>
    )

    function renderRequiredFields() {
        return (
            <SubContentBox label={"Mandatory Fields"} defaultOpen={true}>
                <Box sx={SX.subContent}>
                    <TextField
                        size={"small"}
                        label={"Name"}
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                    />
                    <Box sx={SX.inputTitle}>Nodes</Box>
                    <DynamicInputs
                        InputProps={SX.input}
                        minLength={3}
                        inputs={nodes}
                        onChange={handleNodesUpdate}
                        editable={true}
                        placeholder={"Node "}
                    />
                    <OptionsVault
                        type={VaultType.SSH_KEY}
                        selected={options.vaults.sshKeyId}
                        onUpdate={handleVaultUpdate}
                    />
                </Box>
            </SubContentBox>
        )
    }

    function renderDockerImage() {
        return (
            <SubContentBox label={"Docker Image"}>
                <Box sx={SX.subContent}>
                    <TextField
                        fullWidth
                        size={"small"}
                        label={"Image"}
                        value={image.uri}
                        onChange={v => handleImageChange("uri", v.target.value)}
                    />
                    <TextField
                        fullWidth
                        size={"small"}
                        label={"Volume"}
                        value={image.volume}
                        onChange={v => handleImageChange("volume", v.target.value)}
                    />
                    <TextField
                        fullWidth
                        size={"small"}
                        label={"Restart"}
                        value={image.restart}
                        onChange={v => handleImageChange("restart", v.target.value)}
                    />
                    <Box sx={SX.inputTitle}>Common environment variables</Box>
                    <DynamicInputs
                        InputProps={SX.input}
                        InputSize={"542px"}
                        inputs={commonEnv}
                        onChange={setCommonEnv}
                        editable={true}
                        placeholder={"Env "}
                    />
                    {Object.entries(uniqueEnv).map(([node, envs]) => (
                        <>
                            <Box sx={SX.inputTitle}>Node <Code sx={{fontSize: "13px"}}>{node}</Code> environment variables</Box>
                            <DynamicInputs
                                InputProps={SX.input}
                                InputSize={"266px"}
                                minLength={2}
                                inputs={envs}
                                onChange={v => handleEnvUpdates(node, v)}
                                editable={true}
                                placeholder={"Env "}
                            />
                        </>
                    ))}
                </Box>
            </SubContentBox>
        )
    }

    function renderClusterOptions() {
        return (
            <SubContentBox label={"Cluster Options"}>
                <Options options={options} onUpdate={handleOptionsUpdate}/>
            </SubContentBox>
        )
    }

    function handleCallEnvUpdates(node: string, envs: string[]) {
        setUniqueEnv(prev => ({...prev, [node]: envs}))
    }

    function handleCallNodesUpdate(nodes: string[]) {
        setNodes(nodes)
        setUniqueEnv(prev => Object.fromEntries(
            nodes.filter(n => !!n)
                .map(node => prev[node] !== undefined ? [node, prev[node]] : [node, image.defaultUniqueEnv(node)])
        ))
    }

    function handleCallImageUpdate(key: string, value: string) {
        setImage(prev => ({...prev, [key]: value}))
    }

    function handleCallOptionsUpdate(opt: ClusterOptions) {
        setOptions(prev => ({...prev, ...opt}))
    }

    function handleCallVaultUpdate(t: VaultType, s?: string) {
        setOptions(prev => ({...prev, vaults: {...prev.vaults, [VaultOptions[t].key]: s}}))
    }
}
