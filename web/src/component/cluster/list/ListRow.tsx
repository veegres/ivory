import {Box, TableRow, Tooltip} from "@mui/material";
import {useEffect, useRef, useState} from "react";
import {Cluster} from "../../../type/cluster";
import {DynamicInputs} from "../../view/input/DynamicInputs";
import {ListCellRead} from "./ListCellRead";
import {ListCellUpdate} from "./ListCellUpdate";
import {ListCell} from "./ListCell";
import {getDomains, getSidecars} from "../../../app/utils";
import {ListCellChip} from "./ListCellChip";
import {useInstanceDetection} from "../../../hook/InstanceDetection";
import {WarningAmberRounded} from "@mui/icons-material";
import {SxPropsMap} from "../../../type/common";

const SX: SxPropsMap = {
    actions: {display: "flex", justifyContent: "flex-end", alignItems: "center", gap: 1},
    warning: {padding: "3px 0"},
    removed: {color: "text.secondary", fontSize: "12px", textAlign: "center", marginTop: "5px"},
}

type Props = {
    cluster: Cluster,
    editable: boolean,
    toggle?: () => void,
}

export function ListRow(props: Props) {
    const {cluster, editable, toggle} = props
    const {name, instances, credentials, certs, tags} = cluster

    const ref = useRef<HTMLTableRowElement | null>(null)
    const [stateNodes, setStateNodes] = useState(getDomains(instances))
    const instanceDetection = useInstanceDetection(cluster, instances)

    useEffect(handleEffectScroll, [instanceDetection.active])
    useEffect(handleEffectInstancesUpdate, [instances])

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
                    placeholder={`Instance `}
                    onChange={n => setStateNodes(n)}
                    variant={"outlined"}
                />
                {!toggle && (
                    <Box sx={SX.removed}>
                        This cluster is not on that list, but it is still is checked, just uncheck it.
                    </Box>
                )}
            </ListCell>
            <ListCell>
                <Box sx={SX.actions}>
                    {instanceDetection.warning && (
                        <Tooltip title={"Problems were detected, select the cluster to see it"} placement={"top"}>
                            <WarningAmberRounded color={"warning"}/>
                        </Tooltip>
                    )}
                    {renderButtons()}
                </Box>
            </ListCell>
        </TableRow>
    )

    function renderButtons() {
        if (!toggle) return
        return !editable ? (
            <ListCellRead name={name} toggle={toggle}/>
        ) : (
            <ListCellUpdate
                name={name}
                instances={getSidecars(stateNodes)}
                credentials={credentials}
                certs={certs}
                tags={tags}
                toggle={toggle}
                onUpdate={instanceDetection.refetch}
                onClose={() => setStateNodes(getDomains(instances))}
            />
        )
    }

    function handleEffectScroll() {
        if (instanceDetection.active) ref.current?.scrollIntoView({block: "nearest"})
    }

    function handleEffectInstancesUpdate() {
        setStateNodes(getDomains(instances))
    }
}
