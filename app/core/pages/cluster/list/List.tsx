import {Box} from "@mui/material"

import {useRouterClusterList} from "../../../../features/cluster/api/ClusterHook"
import {Feature} from "../../../../features/Feature"
import {ManageAccess} from "../../../../features/management/component/ManageAccess"
import {ErrorSmart} from "../../../../shared/component/box/ErrorSmart"
import {PageMainBox} from "../../../../shared/component/box/PageMainBox"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {SxPropsFormatter} from "../../../../shared/helper/HelperUtils"
import {ListKeepers} from "./ListKeepers"
import {ListTable} from "./ListTable"
import {ListTags} from "./ListTags"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 0.5},
    filter: {display: "flex", flexDirection: "column", gap: 0.5},
}

export function List() {
    const clusters = useRouterClusterList()

    return (
        <Box sx={SX.box}>
            <Box sx={[SX.filter, SxPropsFormatter.style.pageMargin]}>
                <ListKeepers/>
                <ManageAccess feature={Feature.ViewTagList}>
                    <ListTags/>
                </ManageAccess>
            </Box>
            <PageMainBox>
                {clusters.error ? <ErrorSmart error={clusters.error}/> : (
                    <ListTable list={Object.values(clusters.data ?? [])} fetching={clusters.isFetching} pending={clusters.isPending}/>
                )}
            </PageMainBox>
        </Box>
    )
}
