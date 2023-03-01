import {Database, QueryType} from "../../../app/types";
import {Stack} from "@mui/material";
import {QueryItem} from "./QueryItem";
import {useQuery} from "@tanstack/react-query";
import {queryApi} from "../../../app/api";
import {QueryAdd} from "./QueryAdd";

type Props = {
    type: QueryType,
    cluster: string,
    db: Database,
}

export function Query(props: Props) {
    const {type, cluster, db} = props
    const query = useQuery(["query", "map"], () => queryApi.map(type))

    return (
        <Stack gap={1}>
            <QueryAdd type={type}/>
            {Object.entries(query.data ?? {}).map(([key, value]) => (
                <QueryItem key={key} id={key} query={value} cluster={cluster} db={db}/>
            ))}
        </Stack>
    )
}
