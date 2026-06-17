import {Box, Collapse, Skeleton} from "@mui/material"
import {TransitionGroup} from "react-transition-group"

import {ErrorSmart} from "../../../shared/component/box/ErrorSmart"
import {StylePropsMap} from "../../../shared/helper/type"
import {Feature} from "../../feature"
import {ManageAccess} from "../../management/component/ManageAccess"
import {useRouterQueryList} from "../api/hook"
import {Connection, Type} from "../api/type"
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
            <ManageAccess feature={Feature.ManageQueryCrudCreate}>
                <QueryTemplateNew type={type} connection={connection}/>
            </ManageAccess>
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
