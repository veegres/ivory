import {Box} from "@mui/material";
import {QueryLogItem} from "./QueryLogItem";
import {ClearAllIconButton, RefreshIconButton} from "../../view/button/IconButtons";
import {ErrorSmart} from "../../view/box/ErrorSmart";
import {NoBox} from "../../view/box/NoBox";
import {useRouterQueryLog, useRouterQueryLogDelete} from "../../../api/query/hook";
import {SxPropsMap} from "../../../app/type";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    header: {display: "flex", justifyContent: "space-between", alignItems: "center"},
    label: {color: "text.secondary", fontSize: "13.5px"},
    info: {display: "flex", alignItems: "center", gap: 1}
}

type Props = {
    queryId: string,
}

export function QueryLog(props: Props) {
    const {queryId} = props
    const result = useRouterQueryLog(queryId)
    const clear = useRouterQueryLogDelete(queryId)

    return (
        <Box sx={SX.box}>
            <Box sx={SX.header}>
                <Box>Previous Responses</Box>
                <Box sx={SX.info}>
                    {result.data && (<Box sx={SX.label}>[ {result.data.length} of 10 ]</Box>)}
                    <ClearAllIconButton onClick={() => clear.mutate(queryId)} loading={clear.isPending} disabled={result.data && !result.data.length}/>
                    <RefreshIconButton onClick={() => result.refetch()} loading={result.isFetching}/>
                </Box>
            </Box>
            <Box>
                {renderBody()}
            </Box>
        </Box>
    )


    function renderBody() {
        if (result.error) return <ErrorSmart error={result.error}/>
        if (!result.data?.length) return <NoBox text={"Query log is empty"}/>

        return result.data.toReversed().map((query, index) => (
            <QueryLogItem key={index} index={index} query={query} />
        ))
    }
}
