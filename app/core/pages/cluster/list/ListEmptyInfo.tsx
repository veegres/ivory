import {Alert, Box} from "@mui/material"

import {ClusterDeploy} from "../../../../features/cluster/component/ClusterDeploy"
import {ClusterDetect} from "../../../../features/cluster/component/ClusterDetect"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {KeeperPluginOptions} from "../../../../shared/helper/HelperUtils"
import {useStore} from "../../../../shared/provider/StoreProvider"
import {ListClusterAdd} from "./ListClusterAdd"
import {ListClusterConcepts} from "./ListClusterConcepts"

const SX: SxPropsMap = {
    alert: {"& .MuiAlert-message": {width: "100%"}},
    box: {display: "flex", flexDirection: "column", alignItems: "center", gap: 1},
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
        <Alert sx={SX.alert} severity={"info"} variant={"outlined"} icon={false}>
            <Box sx={SX.box}>
                <Box>{getHeadline()}</Box>
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
            </Box>
        </Alert>
    )

    function getHeadline() {
        const keeperLabel = KeeperPluginOptions[keeper].name
        const tagLabel = tags[0] === "ALL" ? "" : ` tagged “${tags.join(", ")}”`
        return `There are no ${keeperLabel} clusters${tagLabel} yet`
    }
}
