import {Box, Chip, Tooltip} from "@mui/material";
import {AutoRefreshIconButton, RefreshIconButton} from "../../../view/button/IconButtons";
import {purple} from "@mui/material/colors";
import {getDomain, InstanceColor} from "../../../../app/utils";
import {InfoColorBoxList} from "../../../view/box/InfoColorBoxList";
import {useStoreAction} from "../../../../provider/StoreProvider";
import {Cluster, InstanceDetection} from "../../../../api/cluster/type";
import {SxPropsMap} from "../../../../app/type";

const SX: SxPropsMap = {
    chip: {width: "100%"},
    clusterName: {display: "flex", justifyContent: "center", alignItems: "center", gap: "3px"}
}

type Props = {
    cluster: Cluster,
    instanceDetection: InstanceDetection,
}

export function ListCellChip(props: Props) {
    const {cluster, instanceDetection} = props
    const {defaultInstance, combinedInstanceMap, warning} = instanceDetection
    const {fetching, active, detection, refetch} = instanceDetection

    const {setCluster} = useStoreAction

    return (
        <Box sx={SX.clusterName}>
            <Tooltip title={renderChipTooltip()} placement={"right-start"} arrow disableInteractive>
                <Chip
                    sx={SX.chip}
                    color={active ? "primary" : "default"}
                    label={cluster.name}
                    onClick={handleClick}
                />
            </Tooltip>
            {instanceDetection.detection === "auto" ? (
                <AutoRefreshIconButton loading={fetching}  onClick={refetch}/>
            ) : (
                <RefreshIconButton loading={fetching} onClick={refetch}/>
            )}
        </Box>
    )

    // TODO we have same component in Overview page probably we should combine them
    function renderChipTooltip() {
        const items = [
            {title: "Detection", label: detection, bgColor: purple[400]},
            {title: "Instance", label: getDomain(defaultInstance.sidecar), bgColor: InstanceColor[defaultInstance.role]},
        ]

        return <InfoColorBoxList items={items} label={"Cluster Detection"}/>
    }

    function handleClick() {
        if (active) {
            setCluster(undefined)
        } else {
            setCluster({cluster, defaultInstance, combinedInstanceMap, warning, detection})
        }
    }
}
