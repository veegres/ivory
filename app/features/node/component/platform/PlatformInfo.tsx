import {Box, Tooltip} from "@mui/material"
import {useEffect, useState} from "react"

import {ErrorSmart} from "../../../../shared/component/box/ErrorSmart"
import {InfoStatusItem, InfoStatusList} from "../../../../shared/component/box/InfoStatusList"
import {SkeletonGroup} from "../../../../shared/component/progress/SkeletonGroup"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {useRouterNodePlatformInfo} from "../../api/NodeHook"
import {PlatformVaultConnection} from "../../api/NodeType"

const SX: SxPropsMap = {
    // Wraps instead of the nowrap pill look elsewhere, so a long value (e.g.
    // a full host model string) breaks onto a second line instead of
    // overflowing its grid column - but clamped to 2 lines with an ellipsis,
    // with the full value available via tooltip, so one very long value
    // can't blow out the row height for every item next to it.
    value: {
        fontSize: "0.8rem", wordBreak: "break-word",
        display: "-webkit-box", WebkitLineClamp: 2, WebkitBoxOrient: "vertical",
        overflow: "hidden", textOverflow: "ellipsis",
    },
}

type Props = {
    connection: PlatformVaultConnection,
}

export function PlatformInfo(props: Props) {
    const {connection} = props
    const [cachedError, setCachedError] = useState<Error>()
    const info = useRouterNodePlatformInfo(connection)

    useEffect(() => {
        if (info.data) setCachedError(undefined)
        if (info.error) setCachedError(info.error)
    }, [info.error, info.data])

    if (cachedError) return <ErrorSmart error={cachedError}/>

    return renderBody()

    function renderBody() {
        if (info.isLoading) return <SkeletonGroup count={4}/>
        if (!info.data || info.data.length === 0) return undefined

        return (
            <InfoStatusList columns={true}>
                {info.data.map((item) => (
                    <InfoStatusItem key={item.key} label={item.key} columns={true}>
                        {renderValue(item.value)}
                    </InfoStatusItem>
                ))}
            </InfoStatusList>
        )
    }

    function renderValue(value: string) {
        return (
            <Tooltip title={value} placement={"top"}>
                <Box sx={SX.value}>{value}</Box>
            </Tooltip>
        )
    }
}
