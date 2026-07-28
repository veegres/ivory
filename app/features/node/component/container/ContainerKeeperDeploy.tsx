import {Rocket} from "@mui/icons-material"
import {Box, Button, TextField} from "@mui/material"
import {useEffect, useMemo, useState} from "react"

import {ErrorSmart} from "../../../../shared/component/box/ErrorSmart"
import {Logs} from "../../../../shared/component/box/Logs"
import {TitledBox} from "../../../../shared/component/box/TitledBox"
import {WarningList} from "../../../../shared/component/box/WarningList"
import {DialogButton} from "../../../../shared/component/button/DialogButton"
import {DeployImageHeader} from "../../../../shared/component/input/DeployImageHeader"
import {SkeletonGroup} from "../../../../shared/component/progress/SkeletonGroup"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {getDeployPlaceholderKeys, getShortUuid} from "../../../../shared/helper/HelperUtils"
import {useDebounce} from "../../../../shared/hook/Debounce"
import {useKeeperDeployForm, useRouterNodeKeeperDeploy, useRouterNodeKeeperDeployPlan} from "../../api/NodeHook"
import {InterpolationVar, KeeperDeployPlanRequest, KeeperPlugin, PlatformVaultConnection} from "../../api/NodeType"

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
    plugin: KeeperPlugin,
    cluster: string,
    singleHost: boolean,
    databaseId?: string,
    sshKeyId?: string,
}

// ContainerKeeperDeploy deploys a keeper onto a single existing node of an
// already-configured cluster: it calls node's own KeeperDeploy action
// directly (no cluster endpoint involved), so it lives in the node/container
// feature rather than cluster - unlike ClusterDeploy, which batches nodes
// through cluster's own /cluster/deploy action.
export function ContainerKeeperDeploy(props: Props) {
    const {connection, plugin, cluster, singleHost, databaseId, sshKeyId} = props

    const [override, setOverride] = useState<string | undefined>(undefined)
    const [keeperPort, setKeeperPort] = useState<string>("")
    const [dbPort, setDbPort] = useState<string>("")

    const nodeDeploy = useRouterNodeKeeperDeploy(connection)
    const {
        deploySpec, image, imageUri, setImageUri, ready, preview, setPreview, inputs, renderField,
        withKeeperPort, withDbCredentials, mandatoryFields, autoFields,
    } = useKeeperDeployForm(plugin)

    useEffect(handleEffectImage, [image])

    const planRequest = useMemo(
        handleMemoPlanRequest,
        // eslint-disable-next-line react-hooks/exhaustive-deps
        [image, imageUri, plugin, cluster, singleHost, connection.host, keeperPort, dbPort, inputs, override]
    )
    const plan = useRouterNodeKeeperDeployPlan(useDebounce(planRequest, 300))
    const planNode = plan.data?.nodes[0]
    const planWarnings = plan.data?.warnings ?? []
    const planValues = plan.data?.values ?? {}

    return (
        <DialogButton
            title={"DEPLOY CONTAINER"}
            variant={"button"}
            renderActions={renderActions()}
            icon={<Rocket fontSize={"small"}/>}
            back={!!nodeDeploy.data}
        >
            {nodeDeploy.data ? <Logs logs={nodeDeploy.data} height={570} auto={false}/> : renderBody()}
        </DialogButton>
    )

    function renderBody() {
        if (deploySpec.isError) return <ErrorSmart error={deploySpec.error}/>
        if (deploySpec.isPending || !ready) return <SkeletonGroup count={3}/>
        return (
            <Box sx={[SX.subContent, {gap: 2}]}>
                {renderClusterInfo()}
                {renderMandatoryFields()}
                {renderImageOptions()}
            </Box>
        )
    }

    function renderActions() {
        const dbVaultMissing = withDbCredentials && !databaseId
        const planReady = !!plan.data && planWarnings.length === 0
        return (
            <Button loading={nodeDeploy.isPending} onClick={handleAction} disabled={!planReady || dbVaultMissing || !sshKeyId}>
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
                            value={cluster}
                            disabled={true}
                        />
                    </Box>
                    <Box sx={SX.between}>
                        {withDbCredentials && (
                            <TextField
                                fullWidth
                                size={"small"}
                                label={"Database Credentials"}
                                value={getShortUuid(databaseId ?? "none")}
                                disabled={true}
                            />
                        )}
                        <TextField
                            fullWidth
                            size={"small"}
                            label={"SSH Credentials"}
                            value={getShortUuid(sshKeyId ?? "none")}
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
                        {withKeeperPort && (
                            <TextField
                                fullWidth
                                size={"small"}
                                type={"number"}
                                label={"Keeper Port"}
                                value={keeperPort}
                                onChange={e => setKeeperPort(e.target.value)}
                            />
                        )}
                        <TextField
                            fullWidth
                            size={"small"}
                            type={"number"}
                            label={"Database Port"}
                            value={dbPort}
                            onChange={e => setDbPort(e.target.value)}
                        />
                    </Box>
                    {mandatoryFields.map(f => renderField(f, planValues))}
                    {image?.fields.fields.some(f => f.derived) && (
                        <Box sx={SX.note}>
                            The configuration is derived for a new cluster on this node,
                            joining an existing cluster is not supported yet
                        </Box>
                    )}
                </Box>
            </TitledBox>
        )
    }

    function renderImageOptions() {
        return (
            <TitledBox title={"Image Options"} island={true}>
                <Box sx={[SX.subContent, {gap: 2}]}>
                    <DeployImageHeader
                        imageUri={imageUri}
                        onImageUriChange={setImageUri}
                        preview={preview}
                        onPreviewChange={setPreview}
                        placeholderKeys={getPlaceholderKeys()}
                    />
                    <WarningList warnings={planWarnings}/>
                    {autoFields.map(f => renderField(f, planValues))}
                    <TextField
                        fullWidth
                        multiline
                        minRows={5}
                        disabled={preview}
                        size={"small"}
                        label={"Options"}
                        value={preview ? planNode?.optionsPreview ?? "" : override ?? planNode?.options ?? ""}
                        onChange={v => setOverride(v.target.value)}
                    />
                </Box>
            </TitledBox>
        )
    }

    function handleMemoPlanRequest(): KeeperDeployPlanRequest | undefined {
        if (!image) return undefined
        return {
            plugin,
            cluster,
            singleHost,
            image: imageUri,
            values: getValues(),
            nodes: [getNode()],
        }
    }

    function handleAction() {
        if (!plan.data) return
        nodeDeploy.mutate({
            plugin,
            cluster,
            singleHost,
            connection,
            node: getNode(),
            image: imageUri,
            values: getValues(),
            vaults: {databaseId, sshKeyId: sshKeyId ?? ""},
        })
    }

    // NOTE: the shared image hook only seeds the built-in field values it owns
    // itself; the keeper/db port defaults are specific to a single-node deploy
    function handleEffectImage() {
        if (!image) return
        setKeeperPort(image.fields.defaults[InterpolationVar.KeeperPort] ?? "")
        setDbPort(image.fields.defaults[InterpolationVar.DbPort] ?? "")
    }

    function getNode() {
        return {
            host: connection.host,
            keeperPort: Number(keeperPort) || undefined,
            dbPort: Number(dbPort) || undefined,
            options: override,
        }
    }

    function getValues() {
        return {...inputs}
    }

    function getPlaceholderKeys() {
        return getDeployPlaceholderKeys(image?.fields, withKeeperPort, withDbCredentials)
    }
}
