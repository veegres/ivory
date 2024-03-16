import {Box, Chip, Tooltip} from "@mui/material";
import {RefreshIconButton} from "../../view/button/IconButtons";
import {purple} from "@mui/material/colors";
import {getDomain, InstanceColor} from "../../../app/utils";
import {InfoColorBoxList} from "../../view/box/InfoColorBoxList";
import {useStoreAction} from "../../../provider/StoreProvider";
import {Cluster, InstanceDetection} from "../../../type/cluster";
import {SxPropsMap} from "../../../type/common";

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

    const {setCluster, setInstance} = useStoreAction()

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
            <RefreshIconButton loading={fetching} onClick={refetch}/>
        </Box>
    )

    function renderChipTooltip() {
        const items = [
            {title: "Detection", label: detection, bgColor: purple[400]},
            {title: "Instance", label: getDomain(defaultInstance.sidecar), bgColor: InstanceColor[defaultInstance.role]},
        ]

        return <InfoColorBoxList items={items}/>
    }

    function handleClick() {
        if (active) {
            setInstance(undefined)
            setCluster(undefined)
        } else {
            setInstance(undefined)
            setCluster({cluster, defaultInstance, combinedInstanceMap, warning, detection})
        }
    }
}
