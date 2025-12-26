import {HealthAndSafety, Wash, WashOutlined} from "@mui/icons-material"
import {Box, Checkbox, OutlinedInput, Tooltip} from "@mui/material"
import {useMemo, useState} from "react"

import {ConnectionRequest, QueryVariety} from "../../../api/postgres"
import {useRouterQueryRun} from "../../../api/query/hook"
import {QueryApi} from "../../../api/query/router"
import {SxPropsMap} from "../../../app/type"
import {getPostgresUrl} from "../../../app/utils"
import {MenuButton} from "../../view/button/MenuButton"
import {Refresher} from "../refresher/Refresher"
import {QueryResponseInfo} from "./QueryResponseInfo"
import {QueryTable} from "./QueryTable"
import {QueryVarieties} from "./QueryVarieties"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    info: {display: "flex", alignItems: "center", justifyContent: "space-between", gap: 1, fontSize: "13.5px"},
    buttons: {display: "flex", alignItems: "center", fontSize: "10px", gap: "4px"},
    limit: {height: "26px", width: "90px", fontSize: "14px"},
    checkbox: {height: "26px", width: "26px"},
}

type Props = {
    varieties?: QueryVariety[],
    connection: ConnectionRequest,
    queryUuid?: string,
    query: string,
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
    const queryKey = queryUuid ? QueryApi.template.key(queryUuid) : QueryApi.console.key()

    return (
        <Box sx={SX.box}>
            {renderInfo()}
            <QueryTable
                connection={connection}
                data={data}
                error={error}
                refetch={refetch}
                loading={isFetching}
            />
        </Box>
    )

    function renderInfo() {
        return (
            <Box sx={SX.info}>
                <QueryResponseInfo
                    url={getPostgresUrl(connection)}
                    options={!error ? data?.options : undefined}
                    time={{start: !error ? data?.startTime : undefined, end: !error ? data?.endTime : undefined}}
                />
                <Box sx={SX.buttons}>
                    {varieties && <QueryVarieties varieties={varieties}/>}
                    <Refresher queryKeys={[queryKey]}/>
                    <MenuButton>
                        <Tooltip title={"Trim and Remove Comments"} placement={"top"}>
                            <Checkbox
                                sx={SX.checkbox}
                                size={"small"}
                                color={"secondary"}
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
                                color={"secondary"}
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
