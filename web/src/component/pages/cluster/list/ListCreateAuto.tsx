import {Box, Button, Dialog, DialogActions, DialogContent, DialogTitle, Divider, TextField} from "@mui/material"
import {useState} from "react"

import {useRouterClusterCreateAuto} from "../../../../api/cluster/hook"
import {AutoRequest} from "../../../../api/cluster/type"
import {Plugin as DbPlugin} from "../../../../api/database/type"
import {Feature} from "../../../../api/feature"
import {Plugin as KeeperPlugin} from "../../../../api/keeper/type"
import {SxPropsMap} from "../../../../app/type"
import {AlertCentered} from "../../../view/box/AlertCentered"
import {AutoIconButton} from "../../../view/button/IconButtons"
import {Access} from "../../../widgets/access/Access"
import {Options} from "../../../widgets/options/Options"

const SX: SxPropsMap = {
    dialog: {minWidth: "1010px"},
    content: {display: "flex", flexDirection: "column", gap: 1, padding: "0 24px"},
    center: {display: "flex", justifyContent: "center", gap: 3},
    node: {display: "flex", gap: 2},
}

const InitialRequest: AutoRequest = {
    name: "",
    plugins: {
        database: DbPlugin.POSTGRES,
        keeper: KeeperPlugin.PATRONI,
    },
    tls: {keeper: false, database: false},
    certs: {},
    vaults: {},
    tags: [],
    host: "",
    port: 8008,
}

type Props = {
    size?: number,
}

export function ListCreateAuto(props: Props) {
    const {size} = props
    const [request, setRequest] = useState(InitialRequest)
    const [open, setOpen] = useState(false)
    const updateCluster = useRouterClusterCreateAuto(handleSuccessUpdate)

    return (
        <Access feature={Feature.ManageClusterCreate}>
            <AutoIconButton tooltip={"Add Cluster Automatically"} size={size} onClick={() => setOpen(!open)}/>
            <Dialog sx={SX.dialog} open={open} onClose={() => setOpen(false)}>
                <DialogTitle sx={SX.center}>Cluster Auto Detection</DialogTitle>
                <DialogContent sx={SX.content}>
                    <AlertCentered text={"Specify only one node and the others will be detected automatically."}/>
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
                    <Options cluster={request} onUpdate={(opt) => setRequest({...request, ...opt})}/>
                </DialogContent>
                <DialogActions sx={SX.center}>
                    <Button color={"inherit"} onClick={() => setOpen(false)}>Cancel</Button>
                    <Button
                        loading={updateCluster.isPending}
                        onClick={() => updateCluster.mutate(request)}
                        disabled={!request.name || !request.host || !request.port}
                    >
                        Create
                    </Button>
                </DialogActions>
            </Dialog>
        </Access>
    )

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
        setOpen(false)
        setRequest(InitialRequest)
    }
}
