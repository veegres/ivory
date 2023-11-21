import {useQuery} from "@tanstack/react-query";
import {clusterApi} from "../../../app/api";
import {useEffect} from "react";
import {PageMainBox} from "../../view/box/PageMainBox";
import {ListTags} from "./ListTags";
import {useStore} from "../../../provider/StoreProvider";
import {ListTable} from "./ListTable";

export function List() {
    const {activeTags} = useStore()
    const tags = activeTags[0] === "ALL" ? undefined : activeTags
    const query = useQuery({queryKey: ["cluster/list"], queryFn: () => clusterApi.list(tags)})
    const {data, isPending, isFetching, error} = query

    // NOTE: we don't need to check update of `query` because it won't change
    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(() => { query.refetch().then() }, [tags])

    return (
        <PageMainBox withMarginTop={"40px"}>
            <ListTags/>
            <ListTable map={data} error={error} isFetching={isFetching} isLoading={isPending}/>
        </PageMainBox>
    )
}
