import {AutoFixHigh} from "@mui/icons-material"
import {Box, Button, Divider, TextField} from "@mui/material"
import {useEffect, useState} from "react"

import {Options} from "../../../core/widgets/options/Options"
import {DialogButton} from "../../../shared/component/button/DialogButton"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {getKeeperDefaultPort} from "../../../shared/helper/HelperUtils"
import {Feature} from "../../Feature"
import {ManageAccess} from "../../management/component/ManageAccess"
import {useRouterNodeKeeperDeploySpec} from "../../node/api/NodeHook"
import {KeeperPlugin} from "../../node/api/NodeType"
import {DbPlugin} from "../../query/api/QueryType"
import {useRouterClusterCreateAuto} from "../api/ClusterHook"
import {AutoRequest} from "../api/ClusterType"

const SX: SxPropsMap = {
    node: {display: "flex", flexWrap: "wrap", gap: 2},
}

const InitialRequest = (keeper: KeeperPlugin, database: DbPlugin) => ({
    name: "", certs: {}, vaults: {}, tags: [],
    plugins: {database, keeper},
    tls: {keeper: false, database: false},
    host: "", port: 0,
}) as AutoRequest

type Props = {
    keeper: KeeperPlugin,
    database: DbPlugin,
    withLabel?: boolean,
    size?: number,
}

export function ClusterDetect(props: Props) {
    const {keeper, database, withLabel = false, size} = props
    const [request, setRequest] = useState(InitialRequest(keeper, database))
    const [portPlugin, setPortPlugin] = useState<KeeperPlugin>()
    const updateCluster = useRouterClusterCreateAuto(handleSuccessUpdate)
    const deploySpec = useRouterNodeKeeperDeploySpec(keeper)

    useEffect(handleEffectDeploySpec, [deploySpec.data, portPlugin, keeper, database])

    return (
        <ManageAccess feature={Feature.ManageClusterCreate}>
            <DialogButton
                title={"DETECT CLUSTER"}
                renderActions={renderActions()}
                icon={<AutoFixHigh/>}
                variant={withLabel ? "button_label" : "button"}
                label={"Detect"}
                size={size}
            >
                <TextField
                    size={"small"}
                    label={"Name"}
                    required
                    value={request.name}
                    onChange={(e) => handleNameUpdate(e.target.value)}
                />
                <Box sx={SX.node}>
                    <TextField
                        fullWidth
                        size={"small"}
                        label={"Domain"}
                        required
                        value={request.host}
                        onChange={(e) => handleHostUpdate(e.target.value)}
                    />
                    <TextField
                        type={"number"}
                        size={"small"}
                        label={"Port"}
                        required
                        value={request.port || ""}
                        onChange={(e) => handlePortUpdate(parseInt(e.target.value))}
                    />
                </Box>
                <Divider variant={"middle"}/>
                <Options options={request} onUpdate={(opt) => setRequest({...request, ...opt})} disablePlugins={true}/>
            </DialogButton>
        </ManageAccess>
    )

    function renderActions() {
        return (
            <Button
                fullWidth={true}
                loading={updateCluster.isPending}
                onClick={() => updateCluster.mutate(request)}
                disabled={!request.name || !request.host || !request.port}
            >
                Detect
            </Button>
        )
    }
    function handleNameUpdate(v: string) {
        setRequest(c => ({...c, name: v}))
    }

    function handleHostUpdate(v: string) {
        setRequest(c => ({...c, host: v}))
    }

    function handlePortUpdate(v: number) {
        setRequest(c => ({...c, port: isNaN(v) ? 0 : v}))
    }

    function handleSuccessUpdate() {
        setRequest(InitialRequest(keeper, database))
        setPortPlugin(undefined)
    }

    // NOTE: seeds the keeper API port default once per selected keeper plugin,
    // the plugin selectors are disabled inside the dialog, so the plugins can
    // only change through the cluster list filter
    function handleEffectDeploySpec() {
        const data = deploySpec.data
        if (!data || portPlugin === keeper) return
        setPortPlugin(keeper)
        setRequest(prev => ({...prev, plugins: {keeper, database}, port: getKeeperDefaultPort(data)}))
    }
}
