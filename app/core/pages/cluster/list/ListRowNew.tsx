import {FormControl, OutlinedInput} from "@mui/material"
import {useState} from "react"

import {getNodeConfigs, KeeperPluginOptions} from "../../../../shared/helper/HelperUtils"
import {useStore} from "../../../../shared/provider/StoreProvider"
import {ListCellUpdate} from "./ListCellUpdate"
import {ListNodeInput} from "./ListNodeInput"
import {ListRowLayout} from "./ListRowLayout"

type Props = {
    show: boolean,
    close: () => void
}

export function ListRowNew(props: Props) {
    const {show, close} = props
    const activeTags = useStore(s => s.activeTags)
    const keeper = useStore(s => s.activeClusterKeeperPlugin)
    const database = KeeperPluginOptions[keeper].dbPlugin
    const [stateName, setStateName] = useState("")
    const [stateNodes, setStateNodes] = useState([""])

    if (!show) return null

    return (
        <ListRowLayout
            renderName={renderName()}
            renderNodes={renderNodes()}
            renderActions={renderActions()}
        />
    )

    function renderName() {
        return (
            <FormControl fullWidth>
                <OutlinedInput
                    placeholder={"Name"}
                    value={stateName}
                    onChange={(event) => setStateName(event.target.value)}
                />
            </FormControl>
        )
    }

    function renderNodes() {
        return (
            <ListNodeInput inputs={stateNodes} editable={true} onChange={n => setStateNodes(n)}/>
        )
    }

    function renderActions() {
        return (
            <ListCellUpdate
                cluster={{
                    name: stateName,
                    plugins: {database, keeper},
                    nodes: getNodeConfigs(stateNodes),
                    tags: activeTags.filter(t => t !== "ALL"),
                    certs: {},
                    vaults: {},
                    tls: {keeper: false, database: false},
                }}
                toggle={toggle}
                onUpdate={clean}
            />
        )
    }

    function toggle() {
        close()
        clean()
    }

    function clean() {
        setStateName("")
        setStateNodes([""])
    }
}
