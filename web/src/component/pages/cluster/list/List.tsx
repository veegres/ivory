import {useEffect} from "react"

import {useRouterClusterList} from "../../../../api/cluster/hook"
import {Permission} from "../../../../api/permission/type"
import {useStore} from "../../../../provider/StoreProvider"
import {ErrorSmart} from "../../../view/box/ErrorSmart"
import {PageMainBox} from "../../../view/box/PageMainBox"
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
            <Access permission={Permission.ViewTagList}><ListTags/></Access>
            {clusters.error ? <ErrorSmart error={clusters.error}/> : (
                <ListTable map={clusters.data} fetching={clusters.isFetching} pending={clusters.isPending}/>
            )}
        </PageMainBox>
    )
}
