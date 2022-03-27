import {Box, Chip, FormControl, OutlinedInput} from "@mui/material";
import {useStore} from "../../provider/StoreProvider";

const SX = {
    nodesCellIcon: {fontSize: 18},
    nodesCellButton: {height: '22px', width: '22px'},
    nodesCellBox: {minWidth: '150px'},
    chipSize: {width: '100%'},
    nodesCellInput: {height: '32px'},
}

type Props = {
    nodes: string[],
    isReadOnly: boolean,
    onChange: (nodes: string[]) => void,
}

export function ClustersRowNodes({ nodes, isReadOnly, onChange } : Props) {
    const {store, setStore} = useStore()

    return (
        <Box display="grid" gridTemplateColumns={`repeat(auto-fill, minmax(${SX.nodesCellBox.minWidth}, 1fr))`} gap={1}>
            {isReadOnly ? renderChips(nodes) : renderInputs(nodes)}
        </Box>
    )

    function renderInputs(nodes: string[]) {
        const nodesWithEmptyElement = nodes[nodes.length - 1] === '' ? nodes : [...nodes, '']

        return nodesWithEmptyElement.map((node, index) => (
            <FormControl key={index}>
                <OutlinedInput
                    sx={SX.nodesCellInput}
                    type="string"
                    placeholder={`Node ${index}`}
                    size="small"
                    value={node}
                    onChange={(event) => handleChange(index, event.target.value)}
                />
            </FormControl>
        ))
    }

    function renderChips(nodes: string[]) {
        return nodes.map((node, index) => {
            const isNodeActive = !!node && store.activeNode === node

            return (
                <Chip
                    key={index}
                    sx={SX.chipSize}
                    label={node ? node : `Node ${index}`}
                    disabled={!node}
                    variant="outlined"
                    color={isNodeActive ? "primary" : "default"}
                    onClick={handleSelectedNodeChange(node, isNodeActive)}
                />
            )
        })
    }

    function handleSelectedNodeChange(selectedNode: string, isNodeActive: boolean) {
        return () => setStore({ activeNode: isNodeActive ? '' : selectedNode })
    }

    function handleChange(index: number, value: string) {
        nodes[index] = value
        onChange(nodes.filter(node => !!node))
    }
}
