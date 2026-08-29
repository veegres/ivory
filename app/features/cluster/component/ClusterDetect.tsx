import {AutoFixHigh} from "@mui/icons-material"
import {Box, Button, TextField} from "@mui/material"
import {useCallback, useEffect, useState} from "react"

import {DialogScreen} from "../../../shared/component/box/DialogScreen"
import {Note} from "../../../shared/component/box/Note"
import {TitledBox} from "../../../shared/component/box/TitledBox"
import {DialogButton} from "../../../shared/component/button/DialogButton"
import {FieldRow} from "../../../shared/component/input/FieldRow"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {Feature} from "../../Feature"
import {ManageAccess} from "../../management/component/ManageAccess"
import {useRouterNodeKeeperDeploySpec} from "../../node/api/NodeHook"
import {KeeperPlugin} from "../../node/api/NodeType"
import {DbPlugin} from "../../query/api/QueryType"
import {useRouterClusterCreateAuto} from "../api/ClusterHook"
import {AutoRequest, Options as ClusterOptions} from "../api/ClusterType"
import {ClusterOptionsBox} from "./ClusterOptionsBox"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 2},
    column: {display: "flex", flexDirection: "column", gap: 1},
    // NOTE: the same indent the section title carries, so the sentence starts
    // under the heading it belongs to rather than flush against the frame
    note: {paddingX: "10px"},
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

// ClusterDetect is the deploy dialog's counterpart for a cluster that already
// runs: the same sections filled in the same way, except that only one node is
// asked for, because the keeper on it reports the rest.
export function ClusterDetect(props: Props) {
    const {keeper, database, withLabel = false, size} = props
    const [request, setRequest] = useState(InitialRequest(keeper, database))
    const [portPlugin, setPortPlugin] = useState<KeeperPlugin>()
    const [submitted, setSubmitted] = useState(false)
    const updateCluster = useRouterClusterCreateAuto(handleSuccessUpdate)
    const deploySpec = useRouterNodeKeeperDeploySpec(keeper)

    const handleOptionsUpdate = useCallback(handleCallOptionsUpdate, [])

    useEffect(handleEffectDeploySpec, [deploySpec.data, portPlugin, keeper, database])

    return (
        <ManageAccess feature={Feature.ManageClusterCreate}>
            <DialogButton
                title={"DETECT CLUSTER"}
                icon={<AutoFixHigh/>}
                variant={withLabel ? "button_label" : "button"}
                label={"Detect"}
                size={size}
            >
                <DialogScreen renderActions={renderActions()}>
                    <Box sx={SX.box}>
                        {renderCluster()}
                        {renderNode()}
                    </Box>
                </DialogScreen>
            </DialogButton>
        </ManageAccess>
    )

    function renderActions() {
        return (
            <Button fullWidth={true} loading={updateCluster.isPending} onClick={handleDetect}>
                Detect
            </Button>
        )
    }

    function renderCluster() {
        return (
            <TitledBox title={"Cluster"} island={true}>
                <Box sx={SX.column}>
                    <TextField
                        fullWidth={true}
                        size={"small"}
                        label={"Name"}
                        value={request.name}
                        error={submitted && !request.name}
                        onChange={(e) => handleNameUpdate(e.target.value)}
                    />
                    <ClusterOptionsBox options={request} onUpdate={handleOptionsUpdate}/>
                </Box>
            </TitledBox>
        )
    }

    // NOTE: one node is the whole of it, unlike the deploy form's card per
    // template command - Ivory asks the keeper on it for the rest of the cluster
    function renderNode() {
        return (
            <TitledBox title={"Node"} island={true}>
                <Box sx={SX.column}>
                    <Box sx={SX.note}>
                        <Note>Any node of the cluster - the others are discovered through its keeper.</Note>
                    </Box>
                    <FieldRow>
                        <TextField
                            size={"small"}
                            label={"Host"}
                            placeholder={"10.0.0.1"}
                            value={request.host}
                            error={submitted && !request.host}
                            onChange={(e) => handleHostUpdate(e.target.value)}
                        />
                        <TextField
                            size={"small"}
                            type={"number"}
                            label={"Port"}
                            value={request.port || ""}
                            error={submitted && !request.port}
                            onChange={(e) => handlePortUpdate(parseInt(e.target.value))}
                        />
                    </FieldRow>
                </Box>
            </TitledBox>
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

    function handleCallOptionsUpdate(opt: ClusterOptions) {
        setRequest(prev => ({...prev, ...opt}))
    }

    // NOTE: Detect stays clickable while fields are missing, exactly as Deploy
    // does - clicking it is what asks for the errors to be shown, and a
    // disabled button could never explain itself
    function handleDetect() {
        setSubmitted(true)
        if (!request.name || !request.host || !request.port) return
        updateCluster.mutate(request)
    }

    function handleSuccessUpdate() {
        setRequest(InitialRequest(keeper, database))
        setPortPlugin(undefined)
        setSubmitted(false)
    }

    // NOTE: seeds the keeper API port default once per selected keeper plugin,
    // the plugin selectors are disabled inside the dialog, so the plugins can
    // only change through the cluster list filter
    function handleEffectDeploySpec() {
        const data = deploySpec.data
        if (!data || portPlugin === keeper) return
        setPortPlugin(keeper)
        setRequest(prev => ({...prev, plugins: {keeper, database}, port: data.keeperPort}))
    }
}
