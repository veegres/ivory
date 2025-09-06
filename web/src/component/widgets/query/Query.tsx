import {StylePropsMap} from "../../../api/management/type";
import {Box, Collapse, Skeleton} from "@mui/material";
import {ErrorSmart} from "../../view/box/ErrorSmart";
import {TransitionGroup} from "react-transition-group";
import {QueryConnection, QueryType} from "../../../api/query/type";
import {QueryTemplateNew} from "./QueryTemplateNew";
import {QueryTemplateView} from "./QueryTemplateView";
import {useRouterInfo} from "../../../api/management/hook";
import {useRouterQueryList} from "../../../api/query/hook";

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
            {isManual && <QueryTemplateNew type={type} connection={connection}/>}
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
                        <QueryTemplateView key={q.id} connection={connection} query={q} manual={isManual}/>
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
