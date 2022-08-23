import {Box} from "@mui/material";

const SX = {
    row: { display: "flex", gap: 1, justifyContent: "space-between" },
    name: { fontWeight: "bold" }
}

type Item = { name: string, value: string, color?: string }
type Props = {
    items: Item[]
}

export function InfoTitle(props: Props) {
    const { items } = props

    return (
        <Box>
            {items.map((item, index) => (
                <Box key={index} sx={SX.row}>
                    <Box sx={SX.name}>{item.name}:</Box>
                    <Box sx={{ color: item.color }}>{item.value}</Box>
                </Box>
            ))}
        </Box>
    )
}
