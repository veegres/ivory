import {Box, Collapse, Skeleton} from "@mui/material"
import {TransitionGroup} from "react-transition-group"

import {Feature} from "../../../features/feature"
import {useRouterQueryList} from "../../../features/query/hook"
import {Connection, Type} from "../../../features/query/type"
import {ErrorSmart} from "../../../shared/component/box/ErrorSmart"
import {StylePropsMap} from "../../../shared/helper/type"
import {Access} from "../access/Access"
import {QueryTemplateNew} from "./QueryTemplateNew"
import {QueryTemplateView} from "./QueryTemplateView"

const style: StylePropsMap = {
    box: {display: "flex", flexDirection: "column", gap: "8px"},
}

type Props = {
    type: Type,
    connection: Connection,
}

export function Query(props: Props) {
    const {type, connection} = props
    const query = useRouterQueryList(type)

    return (
        <Box style={style.box}>
            <Access feature={Feature.ManageQueryCrudCreate}>
                <QueryTemplateNew type={type} connection={connection}/>
            </Access>
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
