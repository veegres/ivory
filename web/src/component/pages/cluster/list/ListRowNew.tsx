import {FormControl, OutlinedInput, TableRow} from "@mui/material";
import {DynamicInputs} from "../../../view/input/DynamicInputs";
import {useState} from "react";
import {ListCellUpdate} from "./ListCellUpdate";
import {ListCell} from "./ListCell";
import {SxPropsMap} from "../../../../type/general";
import {getSidecars} from "../../../../app/utils";
import {useStore} from "../../../../provider/StoreProvider";

const SX: SxPropsMap = {
    nodesCellInput: {height: '32px'},
}

type Props = {
    show: boolean,
    close: () => void
}

export function ListRowNew(props: Props) {
    const {show, close} = props
    const {activeTags} = useStore()
    const [stateName, setStateName] = useState('');
    const [stateNodes, setStateNodes] = useState(['']);

    if (!show) return null

    return (
        <TableRow>
            <ListCell>
                <FormControl fullWidth>
                    <OutlinedInput
                        sx={SX.nodesCellInput}
                        placeholder={"Name"}
                        value={stateName}
                        onChange={(event) => setStateName(event.target.value)}
                    />
                </FormControl>
            </ListCell>
            <ListCell>
                <DynamicInputs
                    inputs={stateNodes}
                    editable={true}
                    placeholder={`Instance `}
                    onChange={n => setStateNodes(n)}
                />
            </ListCell>
            <ListCell>
                <ListCellUpdate
                    cluster={{
                        name: stateName,
                        instances: getSidecars(stateNodes),
                        tags: activeTags.filter(t => t !== "ALL"),
                        certs: {},
                        credentials: {},
                        tls: {sidecar: false, database: false}
                    }}
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
