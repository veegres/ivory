import {Box, Tooltip} from "@mui/material";
import {SxPropsMap} from "../../../type/general";
import {QueryRunRequest, QueryVariety} from "../../../type/query";
import {RefreshIconButton} from "../../view/button/IconButtons";
import {QueryVarieties} from "./QueryVarieties";
import {useRouterQueryRun} from "../../../router/query";
import {getPostgresUrl, SxPropsFormatter} from "../../../app/utils";
import {QueryBodyTable} from "./QueryBodyTable";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    info: {display: "flex", alignItems: "center", justifyContent: "space-between", gap: 2},
    label: {color: "text.secondary", cursor: "pointer", fontSize: "13.5px", whiteSpace: "nowrap"},
    word: {whiteSpace: "wrap", wordBreak: "break-all"},
    buttons: {display: "flex", alignItems: "center", gap: 1},
}

type Props = {
    varieties?: QueryVariety[],
    request: QueryRunRequest,
}

export function QueryBodyRun(props: Props) {
    const {varieties, request} = props
    const {connection, queryUuid} = request

    const {data, error, isFetching, refetch}  = useRouterQueryRun(request)

    return (
        <Box sx={SX.box}>
            {renderInfo()}
            <QueryBodyTable
                connection={connection}
                data={data}
                error={error}
                queryUuid={queryUuid}
                loading={isFetching}
            />
        </Box>
    )

    function renderInfo() {
        return (
            <Box sx={SX.info}>
                <Tooltip title={"SENT TO"} placement={"right"} arrow={true}>
                    <Box sx={SxPropsFormatter.merge(SX.label, SX.word)}>
                        [ {getPostgresUrl(connection)} ]
                    </Box>
                </Tooltip>
                <Box sx={SX.buttons}>
                    {!isFetching && data && (
                        <Tooltip title={"DURATION"} placement={"left"} arrow={true}>
                            <Box sx={SX.label}>[ {(data.endTime - data.startTime) / 1000}s ]</Box>
                        </Tooltip>
                    )}
                    {varieties && <QueryVarieties varieties={varieties}/>}
                    <RefreshIconButton color={"success"} loading={isFetching} onClick={() => refetch()}/>
                </Box>
            </Box>
        )
    }
}
