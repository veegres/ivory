import {useQuery} from "@tanstack/react-query";
import {queryApi} from "../../../app/api";
import {TabProps} from "./Overview";
import {Query} from "../../shared/query/Query";

export function OverviewActivities(props: TabProps){
    const {data} = useQuery(["query", "map"], () => queryApi.map())

    return (
        <Query map={data ?? {}}/>
    )
}
