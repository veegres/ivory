import {useEffect} from "react"

import {useRouterClusterList} from "../../../../api/cluster/hook"
import {useRouterTagList} from "../../../../api/tag/hook"
import {useStore} from "../../../../provider/StoreProvider"
import {ErrorSmart} from "../../../view/box/ErrorSmart";
import {PageMainBox} from "../../../view/box/PageMainBox"
import {ListTable} from "./ListTable"
import {ListTags} from "./ListTags"

export function List() {
    const activeTags = useStore(s => s.activeTags)
    const clusters = useRouterClusterList(activeTags)
    const tags = useRouterTagList()

    // NOTE: we don't use queryKey to update it, because it will create a separate request and cause new fetching
    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(() => { clusters.refetch().then() }, [activeTags])

    return (
        <PageMainBox withMarginTop={"40px"}>
            {renderBody()}
        </PageMainBox>
    )

    function renderBody() {
        if (clusters.error) return <ErrorSmart error={clusters.error}/>
        if (tags.error) return <ErrorSmart error={tags.error}/>

        return (
            <>
                <ListTags tags={tags.data ?? []}/>
                <ListTable map={clusters.data} fetching={clusters.isFetching} pending={clusters.isPending || tags.isPending}/>
            </>
        )

    }
}
