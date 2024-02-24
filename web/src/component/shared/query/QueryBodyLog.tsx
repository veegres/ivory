import {Box} from "@mui/material";
import {SxPropsMap} from "../../../type/common";
import {QueryBodyLogItem} from "./QueryBodyLogItem";
import {ClearAllIconButton, RefreshIconButton} from "../../view/button/IconButtons";
import {useMutation, useQuery} from "@tanstack/react-query";
import {queryApi} from "../../../app/api";
import {ErrorSmart} from "../../view/box/ErrorSmart";
import {useMutationOptions} from "../../../hook/QueryCustom";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    header: {display: "flex", justifyContent: "space-between", alignItems: "center"},
    no: {
        display: "flex", alignItems: "center", justifyContent: "center", textTransform: "uppercase",
        padding: "8px 16px", border: "1px solid", borderColor: "divider", lineHeight: "1.54",
    },
}

type Props = {
    queryId: string,
}

export function QueryBodyLog(props: Props) {
    const {queryId} = props

    const result = useQuery({
        queryKey: ["query", "log", queryId],
        queryFn: () => queryApi.getLog(queryId),
        enabled: true, retry: false, refetchOnWindowFocus: false,
    })

    const clearOptions = useMutationOptions([["query", "log", queryId]])
    const clear = useMutation({mutationFn: () => queryApi.deleteLog(queryId), ...clearOptions})

    return (
        <Box sx={SX.box}>
            <Box sx={SX.header}>
                <Box>Previous Responses</Box>
                <Box>
                    <ClearAllIconButton onClick={clear.mutate} loading={clear.isPending} disabled={result.data && !result.data.length}/>
                    <RefreshIconButton onClick={result.refetch} loading={result.isFetching}/>
                </Box>
            </Box>
            <Box>
                {renderBody()}
            </Box>
        </Box>
    )


    function renderBody() {
        if (result.error) return <ErrorSmart error={result.error}/>
        if (!result.data?.length) return <Box sx={SX.no}>Query log is empty</Box>

        return result.data.map((query, index) => (
            <QueryBodyLogItem key={index} index={index} query={query} />
        ))
    }
}
