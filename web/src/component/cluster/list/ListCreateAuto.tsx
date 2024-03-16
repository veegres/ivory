import {AutoIconButton} from "../../view/button/IconButtons";
import {useState} from "react";
import {Box, Button, Dialog, DialogActions, DialogContent, DialogTitle, Divider, TextField} from "@mui/material";
import {Options} from "../../shared/options/Options";
import {SxPropsMap} from "../../../type/common";
import {InfoAlert} from "../../view/box/InfoAlert";
import {ClusterAuto} from "../../../type/cluster";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {ClusterApi} from "../../../app/api";
import {LoadingButton} from "@mui/lab";

const SX: SxPropsMap = {
    dialog: {minWidth: "1010px"},
    content: {display: "flex", flexDirection: "column", gap: 2, padding: "0 24px"},
    center: {display: "flex", justifyContent: "center", gap: 3},
    instance: {display: "flex", gap: 2},
}

const InitialClusterAuto: ClusterAuto = {
    name: "", certs: {}, credentials: {}, tags: [], instance: {host: "", port: 8008},
}

export function ListCreateAuto() {
    const [cluster, setCluster] = useState(InitialClusterAuto)
    const [open, setOpen] = useState(false)

    const updateMutationOptions = useMutationOptions([["cluster/list"], ["tag/list"]], handleSuccessUpdate)
    const updateCluster = useMutation({mutationFn: ClusterApi.createAuto, ...updateMutationOptions})

    return (
        <>
            <AutoIconButton onClick={() => setOpen(!open)}/>
            <Dialog sx={SX.dialog} open={open} onClose={() => setOpen(false)}>
                <DialogTitle sx={SX.center}>Cluster Auto Detection</DialogTitle>
                <DialogContent sx={SX.content}>
                    <InfoAlert text={`Specify only one instance and the others will be detected automatically.`}/>
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
                    <LoadingButton
                        loading={updateCluster.isPending}
                        onClick={() => updateCluster.mutate(cluster)}
                    >
                        Save
                    </LoadingButton>
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
