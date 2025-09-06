import {DeleteIconButton, EditIconButton} from "../../../view/button/IconButtons";
import {Box} from "@mui/material";
import {SxPropsMap} from "../../../../api/management/type";
import {useRouterClusterDelete} from "../../../../api/cluster/hook";

const SX: SxPropsMap = {
    box: {display: "flex", justifyContent: "flex-end"},
}

type Props = {
    name: string
    toggle: () => void
}

export function ListCellRead(props: Props) {
    const {toggle, name} = props

    const deleteCluster = useRouterClusterDelete()
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
