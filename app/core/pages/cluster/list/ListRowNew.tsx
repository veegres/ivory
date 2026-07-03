import {FormControl, OutlinedInput, TableRow} from "@mui/material"
import {useState} from "react"

import {KeeperPlugin} from "../../../../features/node/api/type"
import {DbPlugin} from "../../../../features/query/api/type"
import {SxPropsMap} from "../../../../shared/helper/type"
import {getNodeConfigs} from "../../../../shared/helper/utils"
import {useStore} from "../../../../shared/provider/StoreProvider"
import {ListCell} from "./ListCell"
import {ListCellUpdate} from "./ListCellUpdate"
import {ListNodeInput} from "./ListNodeInput"

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
                <ListNodeInput inputs={stateNodes} editable={true} onChange={n => setStateNodes(n)}/>
            </ListCell>
            <ListCell width={"130px"}>
                <ListCellUpdate
                    cluster={{
                        name: stateName,
                        plugins: {database: DbPlugin.POSTGRES, keeper: KeeperPlugin.PATRONI},
                        nodes: getNodeConfigs(stateNodes),
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
