import {AutoFixHigh} from "@mui/icons-material"
import {Box, Button, Divider, TextField} from "@mui/material"
import {useState} from "react"

import {useRouterClusterCreateAuto} from "../../../../api/cluster/hook"
import {AutoRequest} from "../../../../api/cluster/type"
import {Plugin as DbPlugin} from "../../../../api/database/type"
import {Feature} from "../../../../api/feature"
import {Plugin as KeeperPlugin} from "../../../../api/keeper/type"
import {SxPropsMap} from "../../../../app/type"
import {DialogButton} from "../../../view/button/DialogButton"
import {Access} from "../../../widgets/access/Access"
import {Options} from "../../../widgets/options/Options"

const SX: SxPropsMap = {
    dialog: {minWidth: "1010px"},
    content: {display: "flex", flexDirection: "column", gap: 1, padding: "0px 24px"},
    center: {display: "flex", justifyContent: "center", gap: 3},
    node: {display: "flex", gap: 2},
}

const InitialRequest: AutoRequest = {
    name: "", certs: {}, vaults: {}, tags: [],
    plugins: {database: DbPlugin.POSTGRES, keeper: KeeperPlugin.PATRONI},
    tls: {keeper: false, database: false},
    host: "", port: 8008,
}

type Props = {
    keeper: KeeperPlugin,
}

export function ListDetectCluster(_: Props) {
    const [request, setRequest] = useState(InitialRequest)
    const updateCluster = useRouterClusterCreateAuto(handleSuccessUpdate)

    return (
        <Access feature={Feature.ManageClusterCreate}>
            <DialogButton title={"DETECT CLUSTER"} renderActions={renderActions()} icon={<AutoFixHigh/>}>
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
                <Options options={request} onUpdate={(opt) => setRequest({...request, ...opt})}/>
            </DialogButton>
        </Access>
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
        setRequest(InitialRequest)
    }
}
