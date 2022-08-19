import {FormControl, OutlinedInput, TableRow} from "@mui/material";
import {DynamicInputs} from "../view/DynamicInputs";
import {useState} from "react";
import {ClustersRowUpdate} from "./ClustersRowUpdate";
import {ClustersCell} from "./ClustersCell";

const SX = {
    nodesCellInput: {height: '32px'},
}

type Props = {
    show: boolean,
    close: () => void
}

export function ClustersRowNew(props: Props) {
    const { show, close } = props
    const [stateName, setStateName] = useState('');
    const [stateNodes, setStateNodes] = useState(['']);

    if (!show) return null

    return (
        <TableRow>
            <ClustersCell>
                <FormControl fullWidth>
                    <OutlinedInput
                        sx={SX.nodesCellInput}
                        placeholder="Name"
                        value={stateName}
                        onChange={(event) => setStateName(event.target.value)}
                    />
                </FormControl>
            </ClustersCell>
            <ClustersCell>
                <DynamicInputs
                    inputs={stateNodes}
                    editable={true}
                    placeholder={`Instance`}
                    onChange={n => setStateNodes(n)}
                />
            </ClustersCell>
            <ClustersCell>
                <ClustersRowUpdate
                    name={stateName}
                    nodes={stateNodes}
                    toggle={toggle}
                    onUpdate={clean}
                />
            </ClustersCell>
        </TableRow>
    )

    function toggle() {
        close()
        clean()
    }

    function clean() {
        setStateName('')
        setStateNodes([''])
    }
}
