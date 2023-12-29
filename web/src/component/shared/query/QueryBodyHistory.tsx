import {Box} from "@mui/material";
import {SxPropsMap} from "../../../type/common";
import {QueryBodyHistoryItem} from "./QueryBodyHistoryItem";
import {ClearAllIconButton} from "../../view/button/IconButtons";
import {useMutation, useQuery} from "@tanstack/react-query";
import {queryApi} from "../../../app/api";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    header: {display: "flex", justifyContent: "space-between", alignItems: "center", padding: "0 15px"},
    no: {
        display: "flex", alignItems: "center", justifyContent: "center", textTransform: "uppercase",
        padding: "8px 16px", border: "1px solid", borderColor: "divider", lineHeight: "1.54",
    },
}

type Props = {
    queryId: string,
}

export function QueryBodyHistory(props: Props) {
    const {queryId} = props

    const result = useQuery({
        queryKey: ["query", "history", queryId],
        queryFn: () => queryApi.getHistory(queryId),
        enabled: true, retry: false, refetchOnWindowFocus: false,
    })

    const clear = useMutation({
        mutationFn: () => queryApi.deleteHistory(queryId),
        onSuccess: () => result.refetch(),
    })

    return (
        <Box sx={SX.box}>
            <Box sx={SX.header}>
                <Box>Previous Responses</Box>
                <ClearAllIconButton onClick={clear.mutate} color={"error"} disabled={!result.data?.length}/>
            </Box>
            <Box>
                {renderBody()}
            </Box>
        </Box>
    )


    function renderBody() {
        if (!result.data?.length) {
            return <Box sx={SX.no}>Query doesn't have history yet</Box>
        }

        return result.data.map((query, index) => (
            <QueryBodyHistoryItem key={index} query={query} />
        ))
    }
}
