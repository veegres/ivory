import {Box, Button, Dialog, DialogActions, DialogContent, DialogTitle, Divider, TextField} from "@mui/material"
import {useState} from "react"

import {useRouterClusterCreateAuto} from "../../../../api/cluster/hook"
import {ClusterAuto} from "../../../../api/cluster/type"
import {SxPropsMap} from "../../../../app/type"
import {AlertCentered} from "../../../view/box/AlertCentered"
import {AutoIconButton} from "../../../view/button/IconButtons"
import {Options} from "../../../widgets/options/Options"

const SX: SxPropsMap = {
    dialog: {minWidth: "1010px"},
    content: {display: "flex", flexDirection: "column", gap: 2, padding: "0 24px"},
    center: {display: "flex", justifyContent: "center", gap: 3},
    instance: {display: "flex", gap: 2},
}

const InitialClusterAuto: ClusterAuto = {
    name: "", tls: {sidecar: false, database: false}, certs: {}, credentials: {}, tags: [], instance: {host: "", port: 8008},
}

type Props = {
    size?: number,
}

export function ListCreateAuto(props: Props) {
    const {size} = props
    const [cluster, setCluster] = useState(InitialClusterAuto)
    const [open, setOpen] = useState(false)
    const updateCluster = useRouterClusterCreateAuto(handleSuccessUpdate)

    return (
        <>
            <AutoIconButton tooltip={"Add Cluster Automatically"} size={size} onClick={() => setOpen(!open)}/>
            <Dialog sx={SX.dialog} open={open} onClose={() => setOpen(false)}>
                <DialogTitle sx={SX.center}>Cluster Auto Detection</DialogTitle>
                <DialogContent sx={SX.content}>
                    <AlertCentered text={"Specify only one instance and the others will be detected automatically."}/>
                    <TextField
                        size={"small"}
                        label={"Name"}
                        required
                        value={cluster.name}
                        onChange={(e) => handleNameUpdate(e.target.value)}
                    />
                    <Box sx={SX.instance}>
                        <TextField
                            fullWidth
                            size={"small"}
                            label={"Domain"}
                            required
                            value={cluster.instance.host}
                            onChange={(e) => handleHostUpdate(e.target.value)}
                        />
                        <TextField
                            type={"number"}
                            size={"small"}
                            label={"Port"}
                            required
                            value={cluster.instance.port}
                            onChange={(e) => handlePortUpdate(parseInt(e.target.value))}
                        />
                    </Box>
                    <Divider variant={"middle"}/>
                    <Options cluster={cluster} onUpdate={(opt) => setCluster({...cluster, ...opt})}/>
                </DialogContent>
                <DialogActions sx={SX.center}>
                    <Button color={"inherit"} onClick={() => setOpen(false)}>Cancel</Button>
                    <Button
                        loading={updateCluster.isPending}
                        onClick={() => updateCluster.mutate(cluster)}
                    >
                        Save
                    </Button>
                </DialogActions>
            </Dialog>
        </>
    )

    function handleNameUpdate(v: string) {
        setCluster(c => ({...c, name: v}))
    }

    function handleHostUpdate(v: string) {
        setCluster(c => ({...c, instance: {...c.instance, host: v}}))
    }

    function handlePortUpdate(v: number) {
        setCluster(c => ({...c, instance: {...c.instance, port: v}}))
    }

    function handleSuccessUpdate() {
        setOpen(false)
        setCluster(InitialClusterAuto)
    }
}
