import {QueryType} from "../../../app/types";
import {Stack} from "@mui/material";
import {QueryItem} from "./QueryItem";
import {useQuery} from "@tanstack/react-query";
import {queryApi} from "../../../app/api";

type Props = {
    type: QueryType,
}

export function Query(props: Props) {
    const {type} = props
    const query = useQuery(["query", "map"], () => queryApi.map(type))

    return (
        <Stack gap={1}>
            {Object.entries(query.data ?? {}).map(([key, value]) => (
                <QueryItem key={key} id={key} query={value}/>
            ))}
        </Stack>
    )
}
