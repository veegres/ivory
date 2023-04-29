import {Database, StylePropsMap, SxPropsMap} from "../../../type/common";
import {Box, Collapse, Skeleton, ToggleButton, Tooltip} from "@mui/material";
import {QueryItem} from "./QueryItem";
import {useQuery} from "@tanstack/react-query";
import {queryApi} from "../../../app/api";
import {ErrorSmart} from "../../view/box/ErrorSmart";
import {useState} from "react";
import {TransitionGroup} from "react-transition-group";
import {QueryType} from "../../../type/query";
import {DatabaseBox} from "../../view/box/DatabaseBox";
import {QueryNew} from "./QueryNew";
import {Add} from "@mui/icons-material";

const SX: SxPropsMap = {
    info: {display: "flex", justifyContent: "space-between", alignItems: "center", margin: "0 5px"},
    toggle: {padding: "0px"},
}

const style: StylePropsMap = {
    box: {display: "flex", flexDirection: "column", gap: "8px"},
}

type Props = {
    type: QueryType,
    credentialId?: string,
    db: Database,
}

export function Query(props: Props) {
    const {type, credentialId, db} = props
    const [show, setShow] = useState(false)
    const query = useQuery(["query", "map", type], () => queryApi.list(type))

    if (query.isLoading) return renderLoading()
    if (query.error) return <ErrorSmart error={query.error}/>

    return (
        <Box style={style.box}>
            <Box sx={SX.info}>
                <DatabaseBox db={db}/>
                <ToggleButton sx={SX.toggle} value={"add"} size={"small"} selected={show} onClick={() => setShow(!show)}>
                    <Tooltip title={"New Custom Query"} placement={"top"}>
                        <Add/>
                    </Tooltip>
                </ToggleButton>
            </Box>
            <Collapse in={show}>
                <QueryNew type={type}/>
            </Collapse>
            <TransitionGroup style={style.box} appear={false}>
                {(query.data ?? []).map((value) => (
                    <Collapse key={value.id}>
                        <QueryItem query={value} credentialId={credentialId} db={db} type={type}/>
                    </Collapse>
                ))}
            </TransitionGroup>
        </Box>
    )

    function renderLoading() {
        return (
            <Box style={style.transition}>
                <Skeleton width={"100%"} height={42}/>
                <Skeleton width={"100%"} height={42}/>
                <Skeleton width={"100%"} height={42}/>
                <Skeleton width={"100%"} height={42}/>
                <Skeleton width={"100%"} height={42}/>
            </Box>
        )
    }
}
