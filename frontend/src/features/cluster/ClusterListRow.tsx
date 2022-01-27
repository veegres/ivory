import {Box, IconButton, TableCell, TableRow, TextField} from "@mui/material";
import {Add, Delete, Remove, Save} from "@mui/icons-material";
import {useState} from "react";
import {useMutation, useQueryClient} from "react-query";
import {clusterApi} from "../../app/api";

export function ClusterListRow(props: { nodes?: string[], name?: string }) {
    const [name, setName] = useState(props.name ?? '');
    const [nodes, setNodes] = useState(props.nodes && props.nodes.length ? props.nodes : ['']);

    const queryClient = useQueryClient();
    const { mutate } = useMutation(clusterApi.create, {
        onSuccess: () => queryClient.invalidateQueries('cluster/list')
    })

    return (
        <TableRow>
            <TableCell sx={{ verticalAlign: "top", width: '15%' }}>
                <Box display="flex" flexDirection="row" alignItems="center" justifyContent="space-between">
                    <TextField size="small" fullWidth placeholder="Cluster Name" value={name} onChange={(event) => setName(event.target.value)}/>
                    <Box display="flex" flexDirection="column" alignItems="center">
                        <IconButton size="small" onClick={handleAdd}><Add sx={{ fontSize: 10 }} /></IconButton>
                        <IconButton size="small" onClick={handleRemove}><Remove sx={{ fontSize: 10 }} /></IconButton>
                    </Box>
                </Box>
            </TableCell>
            <TableCell sx={{ verticalAlign: "top" }}>
                <Box display="grid" gridTemplateColumns="repeat(auto-fill, minmax(150px, 1fr))" gap={1}>
                    {nodes.map((name, index) => (
                        <Box key={index}>
                            <TextField
                                type="string"
                                label={`Node ${index}`}
                                size="small"
                                value={name}
                                onChange={(event) => handleChange(index, event.target.value)}
                            />
                        </Box>
                    ))}
                </Box>
            </TableCell>
            <TableCell sx={{ verticalAlign: "top", width: '1%', whiteSpace: 'nowrap' }}>
                <IconButton onClick={() => {}}><Delete fontSize="small" /></IconButton>
                <IconButton onClick={handleUpdate}><Save fontSize="small" /></IconButton>
            </TableCell>
        </TableRow>
    )

    function handleChange(index: number, value: string) {
        setNodes((previous) => { previous[index] = value; return [...previous] })
    }

    function handleAdd() {
        nodes.push("")
        setNodes([...nodes])
    }

    function handleRemove() {
        if (nodes.length <= 1) return
        nodes.pop()
        setNodes([...nodes])
    }

    function handleUpdate() {
        mutate({ name, nodes })
    }
}
