import {FormControl, OutlinedInput, TableRow} from "@mui/material"
import {useState} from "react"

import {Plugin as DbPlugin} from "../../../../api/database/type"
import {Plugin as KeeperPlugin} from "../../../../api/keeper/type"
import {SxPropsMap} from "../../../../app/type"
import {getNodeConnections} from "../../../../app/utils"
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
            <ListCell width={"220px"}>
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
                    placeholder={"Node "}
                    onChange={n => setStateNodes(n)}
                />
            </ListCell>
            <ListCell width={"130px"}>
                <ListCellUpdate
                    cluster={{
                        name: stateName,
                        plugins: {database: DbPlugin.POSTGRES, keeper: KeeperPlugin.PATRONI},
                        nodes: getNodeConnections(stateNodes),
                        tags: activeTags.filter(t => t !== "ALL"),
                        certs: {},
                        vaults: {},
                        tls: {keeper: false, database: false},
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
