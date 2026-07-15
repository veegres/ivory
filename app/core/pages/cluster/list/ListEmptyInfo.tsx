import {Box} from "@mui/material"

import {ClusterDeploy} from "../../../../features/cluster/component/ClusterDeploy"
import {ClusterDetect} from "../../../../features/cluster/component/ClusterDetect"
import {AlertCentered} from "../../../../shared/component/box/AlertCentered"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {KeeperPluginOptions} from "../../../../shared/helper/HelperUtils"
import {useStore} from "../../../../shared/provider/StoreProvider"
import {ListClusterAdd} from "./ListClusterAdd"
import {ListClusterConcepts} from "./ListClusterConcepts"

const SX: SxPropsMap = {
    box: {width: "100%", display: "flex", flexDirection: "column"},
    info: {display: "flex", flexDirection: "column", alignItems: "center", gap: 5},
    buttons: {
        display: "flex", alignItems: "center", justifyContent: "space-evenly", gap: 1, flexWrap: "wrap",
        color: "text.primary", width: "100%", maxWidth: 430,
    },
    button: {width: "120px", display: "flex", alignItems: "center", justifyContent: "center"},
}

type Props = {
    onAddManually: () => void,
    disabledAddManually: boolean,
}

export function ListEmptyInfo(props: Props) {
    const {onAddManually, disabledAddManually} = props
    const keeper = useStore(s => s.activeClusterKeeperPlugin)
    const tags = useStore(s => s.activeTags)

    return (
        <Box sx={SX.box}>
            <AlertCentered text={getHeadline()}/>
            <Box sx={SX.info}>
                <Box/>
                <Box sx={SX.buttons}>
                    <Box sx={SX.button}>
                        <ClusterDeploy keeper={keeper} database={KeeperPluginOptions[keeper].dbPlugin} withLabel={true}/>
                    </Box>
                    <Box sx={SX.button}>
                        <ClusterDetect keeper={keeper} database={KeeperPluginOptions[keeper].dbPlugin} withLabel={true}/>
                    </Box>
                    <Box sx={SX.button}>
                        <ListClusterAdd onClick={onAddManually} disabled={disabledAddManually} withLabel={true}/>
                    </Box>
                </Box>
                <ListClusterConcepts/>
                <Box/>
            </Box>
        </Box>
    )

    function getHeadline() {
        const keeperLabel = KeeperPluginOptions[keeper].name
        const tagLabel = tags[0] === "ALL" ? "" : ` tagged “${tags.join(", ")}”`
        return `There are no ${keeperLabel} clusters${tagLabel} yet`
    }
}
