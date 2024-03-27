import {DeleteIconButton, EditIconButton} from "../../../view/button/IconButtons";
import {useMutation} from "@tanstack/react-query";
import {ClusterApi} from "../../../../app/api";
import {Box} from "@mui/material";
import {useMutationOptions} from "../../../../hook/QueryCustom";
import {SxPropsMap} from "../../../../type/general";

const SX: SxPropsMap = {
    box: {display: "flex", justifyContent: "flex-end"},
}

type Props = {
    name: string
    toggle: () => void
}

export function ListCellRead(props: Props) {
    const {toggle, name} = props

    const deleteMutationOptions = useMutationOptions([["cluster/list"], ["tag/list"]])
    const deleteCluster = useMutation({mutationFn: ClusterApi.delete, ...deleteMutationOptions})
    const {isPending, isSuccess} = deleteCluster

    return (
        <Box sx={SX.box}>
            <EditIconButton disabled={isPending || isSuccess} onClick={toggle}/>
            <DeleteIconButton loading={isPending} disabled={isSuccess} onClick={handleDelete}/>
        </Box>
    )

    function handleDelete() {
        deleteCluster.mutate(name)
    }
}
