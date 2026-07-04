import {Edit, Preview, Rocket} from "@mui/icons-material"
import {Box, Button, TextField, ToggleButton, ToggleButtonGroup, Tooltip} from "@mui/material"
import {useState} from "react"

import {Code} from "../../../../shared/component/box/Code"
import {Logs} from "../../../../shared/component/box/Logs"
import {TitledBox} from "../../../../shared/component/box/TitledBox"
import {DialogButton} from "../../../../shared/component/button/DialogButton"
import {SxPropsMap} from "../../../../shared/helper/type"
import {
    DatabaseImageOptions,
    getInterpolatedImageOptions,
    getShortUuid,
    InterpolatedOptionsKeys,
} from "../../../../shared/helper/utils"
import {Cluster} from "../../../cluster/api/type"
import {useRouterNodePlatformUp} from "../../api/hook"
import {PlatformVaultConnection} from "../../api/type"

const SX: SxPropsMap = {
    note: {
        display: "flex", justifyContent: "center", alignItems: "center",
        color: "text.disabled", fontSize: 12, flexWrap: "wrap", gap: 0.5,
    },
    between: {display: "flex", justifyContent: "space-between", alignItems: "center", gap: 1},
    subContent: {display: "flex", flexDirection: "column"},
    toggleButton: {padding: "0px 10px"},
    clusterInfo: {
        "& .MuiListItem-root": {padding: "2px 16px"},
        "& .MuiBox-root": {margin: "2px 0px"},
    },
}

type Props = {
    connection: PlatformVaultConnection,
    cluster: Cluster,
}

export function ContainerOverviewDeploy(props: Props) {
    const {connection, cluster} = props

    const [image, setImage] = useState(DatabaseImageOptions[cluster.plugins.keeper])
    const [options, setOptions] = useState(image.optionStr)
    const [preview, setPreview] = useState(true)

    const [dcs, setDcs] = useState<string>("")
    const [keeperPort, setKeeperPort] = useState<string>("")
    const [dbPort, setDbPort] = useState<string>("")

    const up = useRouterNodePlatformUp(connection)

    return (
        <DialogButton
            title={"DEPLOY CONTAINER"}
            variant={"outlined"}
            renderActions={renderActions()}
            icon={<Rocket fontSize={"small"}/>}
            back={!!up.data}
        >
            {up.data ? <Logs logs={up.data} height={600} auto={false}/> : (
                <Box sx={[SX.subContent, {gap: 2}]}>
                    {renderClusterInfo()}
                    {renderMandatoryFields()}
                    {renderImageOptions()}
                </Box>
            )}
        </DialogButton>
    )

    function renderActions() {
        return (
            <Button loading={up.isPending} onClick={handleAction}>
                Deploy
            </Button>
        )
    }

    function renderClusterInfo() {
        return (
            <TitledBox title={"Cluster"} island={true}>
                <Box sx={[SX.subContent, {gap: 1}]}>
                    <Box sx={SX.between}>
                        <TextField
                            fullWidth
                            size={"small"}
                            label={"Cluster Name"}
                            value={cluster.name}
                            disabled={true}
                        />
                    </Box>
                    <Box sx={SX.between}>
                        <TextField
                            fullWidth
                            size={"small"}
                            label={"Database Credentials"}
                            value={getShortUuid(cluster.vaults.databaseId ?? "none")}
                            disabled={true}
                        />
                        <TextField
                            fullWidth
                            size={"small"}
                            label={"SSH Credentials"}
                            value={getShortUuid(cluster.vaults.sshKeyId ?? "none")}
                            disabled={true}
                        />
                    </Box>
                </Box>
            </TitledBox>
        )
    }

    function renderMandatoryFields() {
        return (
            <TitledBox title={"Mandatory Options"} island={true}>
                <Box sx={[SX.subContent, {gap: 1}]}>
                    <Box sx={SX.between}>
                        <TextField
                            fullWidth
                            size={"small"}
                            type={"number"}
                            label={"Keeper Port"}
                            value={keeperPort}
                            onChange={e => setKeeperPort(e.target.value)}
                        />
                        <TextField
                            fullWidth
                            size={"small"}
                            type={"number"}
                            label={"Database Port"}
                            value={dbPort}
                            onChange={e => setDbPort(e.target.value)}
                        />
                    </Box>
                    <Box sx={SX.between}>
                        <TextField
                            fullWidth={true}
                            size={"small"}
                            label={"DCS (etcd, zookeper, etc)"}
                            helperText={"Example: etcd1:2379, etcd3:2379, etcd3:2379"}
                            value={dcs}
                            onChange={(e) => setDcs(e.target.value)}
                        />
                    </Box>
                </Box>
            </TitledBox>
        )
    }

    function renderImageOptions() {
        const interpolated = getInterpolatedImageOptions(options, {
            cluster: cluster.name,
            host: connection.host,
            keeperPort: Number(keeperPort),
            dbPort: Number(dbPort),
            dcs: dcs,
            dbUser: cluster.vaults.databaseId ? `{{vault:${getShortUuid(cluster.vaults.databaseId)}}}` : "",
            dbPass: cluster.vaults.databaseId ? "********" : "",
        })

        return (
            <TitledBox title={"Image Options"} island={true}>
                <Box sx={[SX.subContent, {gap: 2}]}>
                    <Box sx={SX.between}>
                        <TextField
                            fullWidth={true}
                            size={"small"}
                            label={"Image"}
                            value={image.uri}
                            onChange={v => setImage(prev => ({...prev, uri: v.target.value}))}
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
                    <TextField
                        fullWidth
                        multiline
                        minRows={5}
                        disabled={preview}
                        size={"small"}
                        label={"Options"}
                        value={preview ? interpolated : options}
                        onChange={v => setOptions(v.target.value)}
                    />
                </Box>
            </TitledBox>
        )
    }

    function handleAction() {
        up.mutate({
            connection,
            name: connection.host,
            image: image.uri,
            vaults: {
                databaseId: cluster.vaults.databaseId ?? "",
                sshKeyId: cluster.vaults.sshKeyId ?? "",
            },
            imageOptions: {
                cluster: cluster.name,
                dcs: dcs,
                keeperPort: Number(keeperPort),
                dbPort: Number(dbPort),
            },
            rawImageOptions: options,
        })
    }
}
