import {Box, Checkbox, OutlinedInput, Tooltip} from "@mui/material";
import {SxPropsMap} from "../../../type/general";
import {QueryConnection, QueryVariety} from "../../../type/query";
import {RefreshIconButton} from "../../view/button/IconButtons";
import {QueryVarieties} from "./QueryVarieties";
import {useRouterQueryRun} from "../../../router/query";
import {QueryTable} from "./QueryTable";
import {useMemo, useState} from "react";
import {HealthAndSafety, Wash, WashOutlined} from "@mui/icons-material";
import {MenuButton} from "../../view/button/MenuButton";
import {QueryResponseInfo} from "./QueryResponseInfo";
import {getPostgresUrl} from "../../../app/utils";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    info: {display: "flex", alignItems: "center", justifyContent: "space-between", gap: 1, fontSize: "13.5px"},
    label: {color: "text.secondary", cursor: "pointer", fontSize: "13.5px", whiteSpace: "nowrap"},
    word: {whiteSpace: "wrap", wordBreak: "break-all"},
    infoRight: {display: "flex", alignItems: "center", gap: 1},
    buttons: {display: "flex", alignItems: "center", fontSize: "10px"},
    limit: {height: "26px", width: "90px", fontSize: "14px"},
    checkbox: {height: "26px", width: "26px"},
}

type Props = {
    varieties?: QueryVariety[],
    connection: QueryConnection,
    queryUuid?: string,
    query?: string,
    params?: string[],
}

export function QueryRun(props: Props) {
    const {varieties, connection, queryUuid, query, params} = props
    const [limit, setLimit] = useState("500")
    const [trim, setTrim] = useState(true)

    const queryRunRequest = useMemo(
        () => ({queryUuid, query, connection, queryOptions: {params, trim, limit: limit ? limit : undefined}}),
        [connection, limit, params, query, queryUuid, trim]
    )
    const {data, error, isFetching, refetch} = useRouterQueryRun(queryRunRequest)

    return (
        <Box sx={SX.box}>
            {renderInfo()}
            <QueryTable
                connection={connection}
                data={data}
                error={error}
                queryKey={queryUuid ?? "console"}
                loading={isFetching}
            />
        </Box>
    )

    function renderInfo() {
        return (
            <Box sx={SX.info}>
                <QueryResponseInfo
                    url={getPostgresUrl(connection)}
                    options={data?.options}
                    time={{start: data?.startTime, end: data?.endTime}}
                />
                <Box sx={SX.buttons}>
                    {varieties && <QueryVarieties varieties={varieties}/>}
                    <RefreshIconButton color={"success"} loading={isFetching} onClick={() => refetch()}/>
                    <MenuButton>
                        <Tooltip title={"Trim and Remove Comments"} placement={"top"}>
                            <Checkbox
                                sx={SX.checkbox}
                                size={"small"}
                                color="secondary"
                                icon={<WashOutlined/>}
                                checkedIcon={<Wash/>}
                                onClick={handleTrimClick}
                                checked={trim}
                            />
                        </Tooltip>
                        <Tooltip title={"Safe request with LIMIT"} placement={"top"}>
                            <Checkbox
                                sx={SX.checkbox}
                                size={"small"}
                                color="secondary"
                                icon={<HealthAndSafety/>}
                                checkedIcon={<HealthAndSafety/>}
                                onClick={handleLimitClick}
                                checked={!!limit}
                                disabled={!trim}
                            />
                        </Tooltip>
                        <Tooltip title={"Limit Number"} placement={"top"}>
                            <OutlinedInput
                                sx={SX.limit}
                                size={"small"}
                                value={limit}
                                placeholder={"no limit"}
                                onChange={(e) => setLimit(e.target.value)}
                                disabled={!trim || !limit}
                            />
                        </Tooltip>
                    </MenuButton>
                </Box>
            </Box>
        )
    }

    function handleTrimClick() {
        setTrim(!trim)
        setLimit(!trim ? "500" : "")
    }

    function handleLimitClick() {
        setLimit(!limit ? "500" : "")
    }
}
