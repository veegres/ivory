import {Skeleton} from "@mui/material";

type Props = {
    count: number
}

export function ChartLoading(props: Props) {
    return (
        <>
            {[...Array(props.count).keys()].map((key) => (
                <Skeleton key={key} width={200} height={150}/>
            ))}
        </>
    )
}
