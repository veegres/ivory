import {FormControl, OutlinedInput, TableRow} from "@mui/material"
import {useState} from "react"

import {Type as DbType} from "../../../../api/database/type"
import {KeeperType} from "../../../../api/keeper/type"
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
                        dbType: DbType.POSTGRES,
                        keeperType: KeeperType.PATRONI,
                        nodes: getNodeConnections(KeeperType.PATRONI, DbType.POSTGRES, stateNodes),
                        tags: activeTags.filter(t => t !== "ALL"),
                        certs: {},
                        vaults: {},
                        tls: {keeper: false, database: false},
                        nodesOverview: {},
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
