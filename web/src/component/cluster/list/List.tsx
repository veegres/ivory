import {useQuery} from "@tanstack/react-query";
import {clusterApi} from "../../../app/api";
import {useEffect} from "react";
import {PageBox} from "../../view/box/PageBox";
import {ListTags} from "./ListTags";
import {useStore} from "../../../provider/StoreProvider";
import {ListTable} from "./ListTable";

export function List() {
    const {store: {activeTags}} = useStore()
    const tags = activeTags[0] === "ALL" ? undefined : activeTags
    const query = useQuery(["cluster/list"], () => clusterApi.list(tags))
    const {data, isLoading, isFetching, error} = query

    // NOTE: we don't need to check update of `query` because it won't change
    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(() => { query.refetch().then() }, [tags])

    return (
        <PageBox withMarginTop={"40px"}>
            <ListTags/>
            <ListTable map={data} error={error} isFetching={isFetching} isLoading={isLoading}/>
        </PageBox>
    )
}
