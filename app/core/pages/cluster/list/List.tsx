import {useEffect} from "react"

import {useRouterClusterList} from "../../../../features/cluster/hook"
import {Feature} from "../../../../features/feature"
import {ErrorSmart} from "../../../../shared/component/box/ErrorSmart"
import {PageMainBox} from "../../../../shared/component/box/PageMainBox"
import {useStore} from "../../../../shared/provider/StoreProvider"
import {Access} from "../../../widgets/access/Access"
import {ListTable} from "./ListTable"
import {ListTags} from "./ListTags"

export function List() {
    const activeTags = useStore(s => s.activeTags)
    const clusters = useRouterClusterList(activeTags)

    // NOTE: we don't use queryKey to update it, because it will create a separate request and cause new fetching
    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(() => { clusters.refetch().then() }, [activeTags])

    return (
        <PageMainBox withMarginTop={"40px"}>
            <Access feature={Feature.ViewTagList}><ListTags/></Access>
            {clusters.error ? <ErrorSmart error={clusters.error}/> : (
                <ListTable list={Object.values(clusters.data ?? [])} fetching={clusters.isFetching} pending={clusters.isPending}/>
            )}
        </PageMainBox>
    )
}
