import {Box, TableRow, Tooltip} from "@mui/material";
import {useEffect, useRef, useState} from "react";
import {Cluster} from "../../../../type/cluster";
import {DynamicInputs} from "../../../view/input/DynamicInputs";
import {ListCellRead} from "./ListCellRead";
import {ListCellUpdate} from "./ListCellUpdate";
import {ListCell} from "./ListCell";
import {getDomains, getSidecars, SxPropsFormatter} from "../../../../app/utils";
import {ListCellChip} from "./ListCellChip";
import {useInstanceDetection} from "../../../../hook/InstanceDetection";
import {ErrorOutlineRounded, WarningAmberRounded} from "@mui/icons-material";
import {SxPropsMap} from "../../../../type/general";

const SX: SxPropsMap = {
    actions: {display: "flex", justifyContent: "flex-end", alignItems: "center", gap: 1, minHeight: "32px"},
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
        <TableRow sx={[instanceDetection.active && SxPropsFormatter.style.bgImageSelected, !toggle && SxPropsFormatter.style.bgImageError]} ref={ref}>
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
                />
            </ListCell>
            <ListCell>
                <Box sx={SX.actions}>
                    {instanceDetection.warning && (
                        <Tooltip title={"Problems were detected, select the cluster to see it"} placement={"top"}>
                            <WarningAmberRounded color={"warning"}/>
                        </Tooltip>
                    )}
                    {!toggle && (
                        <Tooltip title={"This cluster is not in this list. It was selected that is why you see it. Just uncheck it."} placement={"top"}>
                            <ErrorOutlineRounded color={"error"}/>
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
