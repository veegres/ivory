import {Box, Collapse, Skeleton} from "@mui/material"
import {TransitionGroup} from "react-transition-group"

import {useRouterQueryList} from "../../../api/query/hook"
import {QueryConnection, QueryType} from "../../../api/query/type"
import {StylePropsMap} from "../../../app/type"
import {ErrorSmart} from "../../view/box/ErrorSmart"
import {QueryTemplateNew} from "./QueryTemplateNew"
import {QueryTemplateView} from "./QueryTemplateView"

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

    return (
        <Box style={style.box}>
            <QueryTemplateNew type={type} connection={connection}/>
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
                        <QueryTemplateView key={q.id} connection={connection} query={q}/>
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
