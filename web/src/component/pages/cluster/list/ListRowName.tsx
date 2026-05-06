import {Box, Chip} from "@mui/material"
import {useEffect} from "react"

import {Cluster} from "../../../../api/cluster/type"
import {SxPropsMap} from "../../../../app/type"
import {useStore, useStoreAction} from "../../../../provider/StoreProvider"
import {AutoRefreshIconButton, RefreshIconButton} from "../../../view/button/IconButtons"

const SX: SxPropsMap = {
    chip: {width: "100%"},
    clusterName: {display: "flex", justifyContent: "center", alignItems: "center", gap: "3px"}
}

type Props = {
    cluster: Cluster,
    active: boolean,
    refresh: () => void,
    loading: boolean,
}

export function ListRowName(props: Props) {
    const {cluster, loading, refresh, active} = props
    const manualKeeper = useStore(s => s.manualKeeper)
    const {setCluster} = useStoreAction

    useEffect(handleEffectActiveClusterUpdate, [active, cluster, setCluster])

    return (
        <Box sx={SX.clusterName}>
            <Chip
                sx={SX.chip}
                color={active ? "primary" : "default"}
                label={cluster.name}
                onClick={handleClick}
            />
            {renderRefresh()}
        </Box>
    )

    function renderRefresh() {
        return !active || !manualKeeper ? (
            <AutoRefreshIconButton
                placement={"right-start"}
                arrow={true}
                loading={loading}
                onClick={refresh}
            />
        ) : (
            <RefreshIconButton
                placement={"right-start"}
                arrow={true}
                loading={loading}
                onClick={refresh}
            />
        )
    }

    // NOTE: we could potentially take the data from tanstack.query, but then we couldn't
    //  see the selected cluster which is not on the list
    function handleEffectActiveClusterUpdate() {
        if (!active) return
        setCluster(cluster)
    }

    function handleClick() {
        if (active) setCluster(undefined)
        else setCluster(cluster)
    }
}
