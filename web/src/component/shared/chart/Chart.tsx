import {Database, SxPropsMap} from "../../../type/common";
import {Box} from "@mui/material";
import {ChartRow} from "./ChartRow";
import {ChartCommon} from "./ChartCommon";
import {ChartDatabase} from "./ChartDatabase";
import {DatabaseBlock} from "../../view/DatabaseBlock";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", justifyContent: "center", gap: 2},
    info: {display: "flex"},
}

type Props = {
    cluster: string,
    db: Database,
}

export function Chart(props: Props) {
    const {cluster, db} = props

    return (
        <Box sx={SX.box}>
            <Box sx={SX.info}>
                <DatabaseBlock db={db}/>
            </Box>
            <ChartRow>
                <ChartCommon cluster={cluster} db={db}/>
            </ChartRow>
            <ChartRow label={db.database && `${db.database}`}>
                <ChartDatabase cluster={cluster} db={db}/>
            </ChartRow>
        </Box>
    )
}
