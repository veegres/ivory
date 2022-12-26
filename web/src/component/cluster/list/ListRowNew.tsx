import {FormControl, OutlinedInput, TableRow} from "@mui/material";
import {DynamicInputs} from "../../view/DynamicInputs";
import {useState} from "react";
import {ListRowUpdate} from "./ListRowUpdate";
import {ListCell} from "./ListCell";

const SX = {
    nodesCellInput: {height: '32px'},
}

type Props = {
    show: boolean,
    close: () => void
}

export function ListRowNew(props: Props) {
    const { show, close } = props
    const [stateName, setStateName] = useState('');
    const [stateNodes, setStateNodes] = useState(['']);

    if (!show) return null

    return (
        <TableRow>
            <ListCell>
                <FormControl fullWidth>
                    <OutlinedInput
                        sx={SX.nodesCellInput}
                        placeholder="Name"
                        value={stateName}
                        onChange={(event) => setStateName(event.target.value)}
                    />
                </FormControl>
            </ListCell>
            <ListCell>
                <DynamicInputs
                    inputs={stateNodes}
                    editable={true}
                    placeholder={`Instance`}
                    onChange={n => setStateNodes(n)}
                />
            </ListCell>
            <ListCell>
                <ListRowUpdate
                    name={stateName}
                    nodes={stateNodes}
                    toggle={toggle}
                    onUpdate={clean}
                />
            </ListCell>
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
