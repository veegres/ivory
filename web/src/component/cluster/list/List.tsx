import {useQuery} from "@tanstack/react-query";
import {clusterApi} from "../../../app/api";
import {useEffect, useMemo} from "react";
import {PageBlock} from "../../view/PageBlock";
import {ListTags} from "./ListTags";
import {useStore} from "../../../provider/StoreProvider";
import {ListTable} from "./ListTable";

export function List() {
    const {store: {activeTags}} = useStore()
    const tags = activeTags[0] === "ALL" ? undefined : activeTags
    const query = useQuery(["cluster/list"], () => clusterApi.list(tags))
    const {data, isLoading, isFetching, error} = query
    const rows = useMemo(() => Object.entries(data ?? {}), [data])

    // NOTE: we don't need to check update of `query` because it won't change
    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(() => { query.refetch().then() }, [tags])

    return (
        <PageBlock withMarginTop={"40px"}>
            <ListTags/>
            <ListTable rows={rows} error={error} isFetching={isFetching} isLoading={isLoading}/>
        </PageBlock>
    )
}
