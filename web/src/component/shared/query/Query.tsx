import {QueryMap} from "../../../app/types";
import {Stack} from "@mui/material";
import {QueryItem} from "./QueryItem";

type Props = {
    map: QueryMap
}

export function Query(props: Props) {
    return (
        <Stack gap={1}>
            {Object.entries(props.map).map(([key, value]) => (
                <QueryItem key={key} id={key} query={value}/>
            ))}
        </Stack>
    )
}
