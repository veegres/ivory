import {Collapse} from "@mui/material"
import {TransitionGroup} from "react-transition-group"

import {NoBox} from "../../../shared/component/box/NoBox"
import {StylePropsMap} from "../../../shared/helper/HelperType"
import {PgCompactTable} from "../api/PgCompactTableType"
import {PgCompactTableJobItem} from "./PgCompactTableJobItem"

const style: StylePropsMap = {
    transition: {display: "flex", flexDirection: "column", gap: "8px"}
}

type Props = {
    cluster: string,
    list: PgCompactTable[],
    refetchList: () => void,
}

export function PgCompactTableJob(props: Props) {
    const {list, cluster, refetchList} = props
    if (list.length === 0) return <NoBox text={"There is no jobs yet"}/>

    return (
        <TransitionGroup style={style.transition} appear={false}>
            {list.map((value) => (
                <Collapse key={value.uuid}>
                    <PgCompactTableJobItem key={value.uuid} item={value} cluster={cluster} refetchList={refetchList}/>
                </Collapse>
            ))}
        </TransitionGroup>
    )
}
