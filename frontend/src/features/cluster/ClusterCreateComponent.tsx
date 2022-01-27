import {Box, IconButton, TableCell, TableRow, TextField} from "@mui/material";
import {Add, Delete, Save} from "@mui/icons-material";
import {useState} from "react";
import {useMutation, useQueryClient} from "react-query";
import {clusterApi} from "../../app/api";

export function ClusterCreateComponent(props: { nodes?: string[], name?: string }) {
    const [name, setName] = useState(props.name ?? '');
    const [nodes, setNodes] = useState(props.nodes && props.nodes.length ? props.nodes : ['']);

    const queryClient = useQueryClient();
    const { mutate } = useMutation(clusterApi.create, {
        onSuccess: () => queryClient.invalidateQueries('cluster/list')
    })

    return (
        <TableRow>
            <TableCell sx={{ verticalAlign: "top" }}>
                <TextField fullWidth placeholder="Cluster Name" size="small" value={name} onChange={(event) => setName(event.target.value)}/>
            </TableCell>
            <TableCell sx={{ verticalAlign: "top" }}>
                <Box display="grid" gridTemplateColumns="repeat(auto-fill, minmax(150px, 1fr))" gap={2} sx={{ padding: '0px 10px' }}>
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
                    <Box>
                        <IconButton onClick={handleAdd}><Add /></IconButton>
                    </Box>
                </Box>
            </TableCell>
            <TableCell sx={{ verticalAlign: "top" }}>
                <IconButton><Delete /></IconButton>
                <IconButton onClick={handleCreate}><Save /></IconButton>
            </TableCell>
        </TableRow>
    )

    function handleChange(index: number, value: string) {
        setNodes((previous) => { previous[index] = value; return [...previous] })
    }

    function handleAdd() {
        setNodes([...nodes, ""])
    }

    function handleCreate() {
        mutate({ name, nodes })
    }
}
