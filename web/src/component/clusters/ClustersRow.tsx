import {Chip, FormControl, OutlinedInput, TableCell, TableRow} from "@mui/material";
import {Cancel, CheckCircle, Delete, Edit} from "@mui/icons-material";
import {useState} from "react";
import {useMutation, useQueryClient} from "react-query";
import {clusterApi} from "../../app/api";
import {ClusterMap, Style} from "../../app/types";
import {ClustersRowButton} from "./ClustersRowButton";
import {ClustersRowNodes} from "./ClustersRowNodes";

const SX = {
    nodesCellIcon: {fontSize: 18},
    nodesCellButton: {height: '22px', width: '22px'},
    nodesCellBox: {minWidth: '150px'},
    nodesCellInput: {height: '32px'},
    chipSize: {width: '100%'},
    cell: {verticalAlign: "top"},
}

const style: Style = {
    thirdCell: {whiteSpace: 'nowrap'}
}

type Props = {
    name: string,
    nodes: string[],
    edit?: { isReadOnly?: boolean, toggleEdit?: () => void, closeNewElement?: () => void }
}

export function ClustersRow({name, nodes, edit = {}}: Props) {
    const {isReadOnly = false, toggleEdit = () => {}, closeNewElement = () => {}} = edit
    const isNewElement = !name

    const [stateName, setStateName] = useState(name);
    const [stateNodes, setStateNodes] = useState(nodes);

    const queryClient = useQueryClient();
    const updateCluster = useMutation(clusterApi.update, {
        onSuccess: (data) => {
            const map = queryClient.getQueryData<ClusterMap>('cluster/list') ?? {} as ClusterMap
            map[data.name] = data.nodes
            queryClient.setQueryData<ClusterMap>('cluster/list', map)
        }
    })
    const deleteCluster = useMutation(clusterApi.delete, {
        onSuccess: (_, newName) => {
            const map = queryClient.getQueryData<ClusterMap>('cluster/list') ?? {} as ClusterMap
            delete map[newName]
            queryClient.setQueryData<ClusterMap>('cluster/list', map)
        }
    })

    return (
        <TableRow>
            <TableCell sx={SX.cell}>
                {renderClusterNameCell()}
            </TableCell>
            <TableCell sx={SX.cell}>
                <ClustersRowNodes nodes={stateNodes} isReadOnly={isReadOnly} onChange={n => setStateNodes(n)} />
            </TableCell>
            <TableCell sx={SX.cell} style={style.thirdCell}>
                {renderClusterActionsCell()}
            </TableCell>
        </TableRow>
    )

    function renderClusterNameCell() {
        if (!isNewElement) {
            return <Chip sx={SX.chipSize} label={stateName}/>
        }

        return (
            <FormControl fullWidth>
                <OutlinedInput
                    sx={SX.nodesCellInput}
                    placeholder="Name"
                    value={stateName}
                    onChange={(event) => setStateName(event.target.value)}
                />
            </FormControl>
        )
    }

    function renderClusterActionsCell() {
        const isDisabled = !stateName
        return isReadOnly ? (
            <>
                <ClustersRowButton icon={<Edit />} tooltip={'Edit'} loading={updateCluster.isLoading} disabled={isDisabled} onClick={toggleEdit} />
                <ClustersRowButton icon={<Delete />} tooltip={'Delete'} loading={deleteCluster.isLoading} disabled={isDisabled} onClick={handleDelete} />
            </>
        ) : (
            <>
                <ClustersRowButton icon={<Cancel/>} tooltip={'Cancel'} loading={updateCluster.isLoading} onClick={handleCancel}/>
                <ClustersRowButton icon={<CheckCircle/>} tooltip={'Save'} loading={updateCluster.isLoading} disabled={isDisabled} onClick={handleUpdate}/>
            </>
        )
    }

    function handleUpdate() {
        if (isNewElement) closeNewElement(); else toggleEdit()
        updateCluster.mutate({name: stateName, nodes: stateNodes})
    }

    function handleCancel() {
        if (isNewElement) closeNewElement(); else toggleEdit()
        setStateNodes([...nodes])
    }

    function handleDelete() {
        deleteCluster.mutate(stateName)
    }
}
