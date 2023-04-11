import {TableRow} from "@mui/material";
import {useEffect, useRef, useState} from "react";
import {Cluster} from "../../../type/cluster";
import {DynamicInputs} from "../../view/input/DynamicInputs";
import {ListCellRead} from "./ListCellRead";
import {ListCellUpdate} from "./ListCellUpdate";
import {ListCell} from "./ListCell";
import {getDomains} from "../../../app/utils";
import {ListCellChip} from "./ListCellChip";
import {useInstanceDetection} from "../../../hook/InstanceDetection";

type Props = {
    cluster: Cluster,
    editable: boolean,
    toggle: () => void,
}

export function ListRow(props: Props) {
    const {cluster, editable, toggle} = props
    const {name, instances} = cluster

    const ref = useRef<HTMLTableRowElement | null>(null)
    const [stateNodes, setStateNodes] = useState(getDomains(instances))
    const instanceDetection = useInstanceDetection(cluster, instances)

    useEffect(handleEffectScroll, [instanceDetection.active])

    return (
        <TableRow ref={ref}>
            <ListCell>
                <ListCellChip cluster={cluster} instanceDetection={instanceDetection}/>
            </ListCell>
            <ListCell>
                <DynamicInputs
                    inputs={stateNodes}
                    colors={instanceDetection.colors}
                    editable={editable}
                    placeholder={`Instance`}
                    onChange={n => setStateNodes(n)}
                />
            </ListCell>
            <ListCell>
                {!editable ? (
                    <ListCellRead name={name} toggle={toggle}/>
                ) : (
                    <ListCellUpdate
                        name={name}
                        nodes={stateNodes}
                        toggle={toggle}
                        onUpdate={instanceDetection.refetch}
                        onClose={() => setStateNodes(getDomains(instances))}
                    />
                )}
            </ListCell>
        </TableRow>
    )

    function handleEffectScroll() {
        if (instanceDetection.active) ref.current?.scrollIntoView({block: "nearest"})
    }
}
