import {TransitionGroup} from "react-transition-group";
import {Collapse} from "@mui/material";
import {OverviewBloatJobItem} from "./OverviewBloatJobItem";
import {InfoAlert} from "../../../view/box/InfoAlert";
import {StylePropsMap} from "../../../../api/management/type";
import {Type} from "../../../../api/bloat/type";

const style: StylePropsMap = {
    transition: {display: "flex", flexDirection: "column", gap: "10px"}
}

type Props = {
    cluster: string,
    list: Type[],
}

export function OverviewBloatJob(props: Props) {
    const {list, cluster} = props
    if (list.length === 0) return <InfoAlert text={"There is no jobs yet"}/>

    return (
        <TransitionGroup style={style.transition} appear={false}>
            {list.map((value) => (
                <Collapse key={value.uuid}>
                    <OverviewBloatJobItem key={value.uuid} item={value} cluster={cluster}/>
                </Collapse>
            ))}
        </TransitionGroup>
    )
}
