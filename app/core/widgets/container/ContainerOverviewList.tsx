import {Box} from "@mui/material"

import {useRouterNodePlatformList} from "../../../features/node/hook"
import {PlatformConnection} from "../../../features/node/type"
import {ErrorSmart} from "../../../shared/component/box/ErrorSmart"
import {SkeletonGroup} from "../../../shared/component/progress/SkeletonGroup"
import {SxPropsMap} from "../../../shared/helper/type"
import scroll from "../../../shared/style/scroll.module.css"

const SX: SxPropsMap = {
    box: {
        padding: "10px 10px 5px", backgroundColor: "background.default", color: "text.secondary",
        borderRadius: 2, border: 0.5, borderColor: "divider",
    },
    list: {
        display: "flex", flexDirection: "column", gap: 0.5, overflowX: "scroll",
        overflowY: "auto", maxHeight: "400px", minHeight: "100px",
    },
    pre: {fontFamily: "'Fira Code', 'Courier New', monospace", fontSize: "12px", whiteSpace: "pre"}
}

type Props = {
    connection: PlatformConnection,
}

export function ContainerOverviewList(props: Props) {
    const {connection} = props
    const list = useRouterNodePlatformList(connection)

    if (list.isError) return <ErrorSmart error={list.error}/>
    if (list.isPending) return <SkeletonGroup count={1}/>

    return (
        <Box sx={SX.box}>
            <Box sx={SX.list} className={scroll.tiny}>
                <Box sx={SX.pre}>{list.data?.join("\n")}</Box>
            </Box>
        </Box>
    )
}
