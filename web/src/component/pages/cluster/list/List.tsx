import {useEffect} from "react"

import {useRouterClusterList} from "../../../../api/cluster/hook"
import {useStore} from "../../../../provider/StoreProvider"
import {PageMainBox} from "../../../view/box/PageMainBox"
import {ListTable} from "./ListTable"
import {ListTags} from "./ListTags"

export function List() {
    const activeTags = useStore(s => s.activeTags)
    const query = useRouterClusterList(activeTags)
    const {data, isPending, isFetching, error} = query

    // NOTE: we don't use queryKey to update it, because it will create separate request and cause new fetching
    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(() => { query.refetch().then() }, [activeTags])

    return (
        <PageMainBox withMarginTop={"40px"}>
            <ListTags/>
            <ListTable map={data} error={error} fetching={isFetching} pending={isPending}/>
        </PageMainBox>
    )
}
