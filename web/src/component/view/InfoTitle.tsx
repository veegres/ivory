import {Box} from "@mui/material";

const SX = {
    row: { display: "flex", gap: 1, justifyContent: "space-between", margin: "3px 0" },
    name: { fontWeight: "bold" },
    value: { borderRadius: "3px", padding: "0 3px" }
}

type Item = { label: string, value: string, bgColor?: string }
type Props = {
    items: Item[]
}

export function InfoTitle(props: Props) {
    const { items } = props

    return (
        <Box>
            {items.map((item, index) => renderItem(item, index))}
        </Box>
    )

    function renderItem(item: Item, index: number) {
        const { bgColor, label, value } = item

        return (
            <Box key={index} sx={SX.row}>
                <Box sx={SX.name}>{label}:</Box>
                <Box sx={SX.value} bgcolor={bgColor}>{value}</Box>
            </Box>
        )
    }
}
