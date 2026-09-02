import {AutoFixHigh} from "@mui/icons-material"
import {Box, Button, TextField} from "@mui/material"
import {useCallback, useEffect, useState} from "react"

import {AlertCentered} from "../../../shared/component/box/AlertCentered"
import {DialogScreen} from "../../../shared/component/box/DialogScreen"
import {PaperBlue} from "../../../shared/component/box/PaperBlue"
import {TitleBox} from "../../../shared/component/box/TitleBox"
import {DialogButton} from "../../../shared/component/button/DialogButton"
import {FieldRow} from "../../../shared/component/input/FieldRow"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {Feature} from "../../Feature"
import {ManageAccess} from "../../management/component/ManageAccess"
import {KeeperPlugin} from "../../node/api/NodeType"
import {DbPlugin} from "../../query/api/QueryType"
import {useRouterClusterCreateAuto} from "../api/ClusterHook"
import {AutoRequest, Options as ClusterOptions} from "../api/ClusterType"
import {ClusterOptionsBox} from "./ClusterOptionsBox"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 2},
    column: {display: "flex", flexDirection: "column", gap: 1},
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
    const [submitted, setSubmitted] = useState(false)
    const updateCluster = useRouterClusterCreateAuto(handleSuccessUpdate)

    const handleOptionsUpdate = useCallback(handleCallOptionsUpdate, [])

    useEffect(handleEffectPlugins, [keeper, database])

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
                        {renderInfo()}
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

    function renderInfo() {
        return (
            <PaperBlue>
                <AlertCentered text={
                    "Host, name, and ports come from the keeper, not your platform - if ports are mapped " +
                    "or forwarded, they may differ from what you'd expect"
                }/>
            </PaperBlue>
        )
    }

    function renderCluster() {
        return (
            <PaperBlue>
                <TitleBox label={"Cluster"} island={true} collapsible={false}>
                    <Box sx={SX.column}>
                        <TextField
                            fullWidth={true}
                            label={"Name"}
                            value={request.name}
                            error={submitted && !request.name}
                            onChange={(e) => handleNameUpdate(e.target.value)}
                        />
                        <ClusterOptionsBox options={request} onUpdate={handleOptionsUpdate}/>
                    </Box>
                </TitleBox>
            </PaperBlue>
        )
    }

    // NOTE: one node is the whole of it, unlike the deploy form's card per
    // template command - Ivory asks the keeper on it for the rest of the cluster
    function renderNode() {
        return (
            <PaperBlue>
                <TitleBox
                    label={"Node"}
                    hint={"Any node of the cluster - the others are discovered through its keeper"}
                    island={true}
                    collapsible={false}
                >
                    <Box sx={SX.column}>
                        <FieldRow>
                            <TextField
                                label={"Host"}
                                placeholder={"10.0.0.1"}
                                value={request.host}
                                error={submitted && !request.host}
                                onChange={(e) => handleHostUpdate(e.target.value)}
                            />
                            <TextField
                                type={"number"}
                                label={"Port"}
                                value={request.port || ""}
                                error={submitted && !request.port}
                                onChange={(e) => handlePortUpdate(parseInt(e.target.value))}
                            />
                        </FieldRow>
                    </Box>
                </TitleBox>
            </PaperBlue>
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
        setSubmitted(false)
    }

    // NOTE: the plugin selectors are disabled inside the dialog, so the pair can
    // only change through the cluster list filter behind it - this is what
    // carries that change into a dialog mounted before it happened. The port is
    // not seeded alongside them: it is the one the keeper actually answers on,
    // which only the operator knows.
    function handleEffectPlugins() {
        setRequest(prev => ({...prev, plugins: {keeper, database}}))
    }
}
