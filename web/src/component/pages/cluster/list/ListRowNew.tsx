import {FormControl, OutlinedInput, TableRow} from "@mui/material"
import {useState} from "react"

import {SxPropsMap} from "../../../../app/type"
import {getSidecars} from "../../../../app/utils"
import {useStore} from "../../../../provider/StoreProvider"
import {DynamicInputs} from "../../../view/input/DynamicInputs"
import {ListCell} from "./ListCell"
import {ListCellUpdate} from "./ListCellUpdate"

const SX: SxPropsMap = {
    nodesCellInput: {height: "32px"},
}

type Props = {
    show: boolean,
    close: () => void
}

export function ListRowNew(props: Props) {
    const {show, close} = props
    const activeTags = useStore(s => s.activeTags)
    const [stateName, setStateName] = useState("")
    const [stateNodes, setStateNodes] = useState([""])

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
                    placeholder={"Instance "}
                    onChange={n => setStateNodes(n)}
                />
            </ListCell>
            <ListCell>
                <ListCellUpdate
                    cluster={{
                        name: stateName,
                        sidecars: getSidecars(stateNodes),
                        tags: activeTags.filter(t => t !== "ALL"),
                        certs: {},
                        credentials: {},
                        tls: {sidecar: false, database: false},
                        unknownInstances: {},
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
        setStateName("")
        setStateNodes([""])
    }
}
