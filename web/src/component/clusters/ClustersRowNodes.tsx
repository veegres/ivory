import {Box, Chip, FormControl, IconButton, OutlinedInput} from "@mui/material";
import {useStore} from "../../provider/StoreProvider";
import {Visibility} from "@mui/icons-material";
import {useEffect, useState} from "react";

const SX = {
    nodesCellIcon: {fontSize: 18},
    nodesCellButton: {height: '22px', width: '22px'},
    nodesCellBox: {minWidth: '150px'},
    chipSize: {width: '100%'},
    nodesCellInput: {height: '32px'},
}

type Props = {
    nodes: string[],
    isReadOnly: false,
    onChange: (nodes: string[]) => string[],
}

export function ClustersRowNodes({ nodes, isReadOnly, onChange } : Props) {
    const [stateNodes, setStateNodes] = useState(nodes.length ? nodes : ['']);
    const {store, setStore} = useStore()

    useEffect(() => { onChange(stateNodes) }, [stateNodes])

    return (
        <Box display="grid" gridTemplateColumns={`repeat(auto-fill, minmax(${SX.nodesCellBox.minWidth}, 1fr))`} gap={1}>
            {stateNodes.map((node, index) => renderClusterNodesCellElement(node, index))}
        </Box>
    )

    function renderClusterNodesCellElement(node: string, index: number) {
        const isNodeActive = !!node && store.activeNode === node

        return isReadOnly ? (
            <Chip
                key={index}
                sx={SX.chipSize}
                label={node ? node : `Node ${index}`}
                disabled={!node}
                variant="outlined"
                color={isNodeActive ? "primary" : "default"}
                onClick={handleSelectedNodeChange(node, isNodeActive)}
            />
        ) : (
            <FormControl key={index}>
                <OutlinedInput
                    sx={SX.nodesCellInput}
                    type="string"
                    endAdornment={renderNodeSelectIcon(node, isNodeActive)}
                    placeholder={`Node ${index}`}
                    size="small"
                    value={node}
                    onChange={(event) => handleChange(index, event.target.value)}
                />
            </FormControl>
        )
    }

    function renderNodeSelectIcon(node: string, isNodeActive: boolean) {
        if (!node) return null

        return (
            <IconButton
                sx={SX.nodesCellButton}
                color={isNodeActive ? "primary" : "default"}
                onClick={handleSelectedNodeChange(node, isNodeActive)}>
                <Visibility sx={SX.nodesCellIcon} />
            </IconButton>
        )
    }

    function handleSelectedNodeChange(selectedNode: string, isNodeActive: boolean) {
        return () => setStore({ activeNode: isNodeActive ? '' : selectedNode })
    }

    function handleChange(index: number, value: string) {
        setStateNodes((previous) => { previous[index] = value; return [...previous] })
    }
}
