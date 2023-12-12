import {Database, StylePropsMap} from "../../../type/common";
import {Box, Collapse, Skeleton} from "@mui/material";
import {QueryItem} from "./QueryItem";
import {useQuery} from "@tanstack/react-query";
import {generalApi, queryApi} from "../../../app/api";
import {ErrorSmart} from "../../view/box/ErrorSmart";
import {TransitionGroup} from "react-transition-group";
import {QueryType} from "../../../type/query";
import {QueryItemNew} from "./QueryItemNew";

const style: StylePropsMap = {
    box: {display: "flex", flexDirection: "column", gap: "8px"},
}

type Props = {
    type: QueryType,
    credentialId?: string,
    db: Database,
}

export function Query(props: Props) {
    const {type, credentialId, db} = props
    const query = useQuery({queryKey: ["query", "map", type], queryFn: () => queryApi.list(type)})
    const info = useQuery({queryKey: ["info"], queryFn: () => generalApi.info()})

    if (query.isPending) return renderLoading()
    if (query.error) return <ErrorSmart error={query.error}/>

    const isManual = info.data?.availability.manualQuery ?? false

    return (
        <Box style={style.box}>
            <TransitionGroup style={style.box} appear={false}>
                {(query.data ?? []).map((value) => (
                    <Collapse key={value.id}>
                        <QueryItem
                            key={value.id}
                            query={value}
                            credentialId={credentialId}
                            db={db} type={type}
                            editable={isManual}
                        />
                    </Collapse>
                ))}
            </TransitionGroup>
            {isManual && <QueryItemNew type={type}/>}
        </Box>
    )

    function renderLoading() {
        return (
            <Box>
                <Skeleton width={"100%"} height={42}/>
                <Skeleton width={"100%"} height={42}/>
                <Skeleton width={"100%"} height={42}/>
                <Skeleton width={"100%"} height={42}/>
                <Skeleton width={"100%"} height={42}/>
            </Box>
        )
    }
}
