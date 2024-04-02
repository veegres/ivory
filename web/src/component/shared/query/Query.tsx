import {StylePropsMap} from "../../../type/general";
import {Box, Collapse, Skeleton} from "@mui/material";
import {ErrorSmart} from "../../view/box/ErrorSmart";
import {TransitionGroup} from "react-transition-group";
import {QueryConnection, QueryType} from "../../../type/query";
import {QueryItemNew} from "./QueryItemNew";
import {QueryItemView} from "./QueryItemView";
import {useRouterInfo} from "../../../router/general";
import {useRouterQueryList} from "../../../router/query";

const style: StylePropsMap = {
    box: {display: "flex", flexDirection: "column", gap: "8px"},
}

type Props = {
    type: QueryType,
    connection: QueryConnection,
}

export function Query(props: Props) {
    const {type, connection} = props
    const query = useRouterQueryList(type)
    const info = useRouterInfo()

    const isManual = info.data?.availability.manualQuery ?? false

    return (
        <Box style={style.box}>
            {isManual && <QueryItemNew type={type} connection={connection}/>}
            {renderList()}
        </Box>
    )

    function renderList() {
        if (query.isPending) return renderLoading()
        if (query.error) return <ErrorSmart error={query.error}/>

        return (
            <TransitionGroup style={style.box} appear={false}>
                {(query.data ?? []).map((q) => (
                    <Collapse key={q.id}>
                        <QueryItemView key={q.id} connection={connection} query={q} manual={isManual}/>
                    </Collapse>
                ))}
            </TransitionGroup>
        )
    }

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
