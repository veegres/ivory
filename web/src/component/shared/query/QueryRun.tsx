import {Box, Checkbox, OutlinedInput, Tooltip} from "@mui/material";
import {SxPropsMap} from "../../../type/general";
import {QueryConnection, QueryVariety} from "../../../type/query";
import {RefreshIconButton} from "../../view/button/IconButtons";
import {QueryVarieties} from "./QueryVarieties";
import {useRouterQueryRun} from "../../../router/query";
import {getPostgresUrl, SxPropsFormatter} from "../../../app/utils";
import {QueryTable} from "./QueryTable";
import {useMemo, useState} from "react";
import {HealthAndSafety, Wash, WashOutlined} from "@mui/icons-material";
import {MenuButton} from "../../view/button/MenuButton";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    info: {display: "flex", alignItems: "center", justifyContent: "space-between", gap: 2},
    label: {color: "text.secondary", cursor: "pointer", fontSize: "13.5px", whiteSpace: "nowrap"},
    word: {whiteSpace: "wrap", wordBreak: "break-all"},
    infoRight: {display: "flex", alignItems: "center", gap: 1},
    buttons: {display: "flex", alignItems: "center"},
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
                <Tooltip title={"SENT TO"} placement={"right"} arrow={true}>
                    <Box sx={SxPropsFormatter.merge(SX.label, SX.word)}>
                        [ {getPostgresUrl(connection)} ]
                    </Box>
                </Tooltip>
                <Box sx={SX.infoRight}>
                    {varieties && <QueryVarieties varieties={varieties}/>}
                    {!isFetching && data && (<>
                        <Tooltip title={"LIMIT"} placement={"top"}>
                            <Box sx={SX.label}>[ {data.options?.limit ?? "NO LIMIT"} ]</Box>
                        </Tooltip>
                        <Tooltip title={"DURATION"} placement={"top"}>
                            <Box sx={SX.label}>[ {(data.endTime - data.startTime) / 1000}s ]</Box>
                        </Tooltip>
                    </>)}
                    <Box sx={SX.buttons}>
                        <RefreshIconButton color={"success"} loading={isFetching} onClick={() => refetch()}/>
                        <MenuButton>
                            <Tooltip title={"Trim and Clean Comments"} placement={"top"}>
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
                            <OutlinedInput
                                sx={SX.limit}
                                size={"small"}
                                value={limit}
                                placeholder={"no limit"}
                                onChange={(e) => setLimit(e.target.value)}
                                disabled={!trim || !limit}
                            />
                        </MenuButton>
                    </Box>
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
