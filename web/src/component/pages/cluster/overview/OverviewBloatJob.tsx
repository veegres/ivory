import {Collapse} from "@mui/material"
import {TransitionGroup} from "react-transition-group"

import {Bloat} from "../../../../api/bloat/type"
import {StylePropsMap} from "../../../../app/type"
import {NoBox} from "../../../view/box/NoBox"
import {OverviewBloatJobItem} from "./OverviewBloatJobItem"

const style: StylePropsMap = {
    transition: {display: "flex", flexDirection: "column", gap: "10px"}
}

type Props = {
    cluster: string,
    list: Bloat[],
    refetchList: () => void,
}

export function OverviewBloatJob(props: Props) {
    const {list, cluster, refetchList} = props
    if (list.length === 0) return <NoBox text={"There is no jobs yet"}/>

    return (
        <TransitionGroup style={style.transition} appear={false}>
            {list.map((value) => (
                <Collapse key={value.uuid}>
                    <OverviewBloatJobItem key={value.uuid} item={value} cluster={cluster} refetchList={refetchList}/>
                </Collapse>
            ))}
        </TransitionGroup>
    )
}
