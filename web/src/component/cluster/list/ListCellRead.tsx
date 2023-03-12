import {DeleteIconButton, EditIconButton} from "../../view/IconButtons";
import {useMutation} from "@tanstack/react-query";
import {clusterApi} from "../../../app/api";
import {Box} from "@mui/material";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {SxPropsMap} from "../../../type/common";

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
    const deleteCluster = useMutation(clusterApi.delete, deleteMutationOptions)
    const {isLoading, isSuccess} = deleteCluster

    return (
        <Box sx={SX.box}>
            <EditIconButton disabled={isLoading || isSuccess} onClick={toggle}/>
            <DeleteIconButton loading={isLoading} disabled={isSuccess} onClick={handleDelete}/>
        </Box>
    )

    function handleDelete() {
        deleteCluster.mutate(name)
    }
}
