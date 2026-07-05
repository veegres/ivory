import {Box} from "@mui/material"
import {Children, ReactNode} from "react"

import {SettingsWrapperScroll} from "../../../core/widgets/settings/SettingsWrapperScroll"
import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1, height: "100%", overflow: "hidden"},
}

type Props = {
    scrollElement?: number,
    children: ReactNode,
}

export function LastElementScrolling(props: Props) {
    const {children, scrollElement} = props
    const scroll = scrollElement ?? Children.count(children) - 1
    return (
        <Box sx={SX.box}>
            {Children.map(children, (child, index) => index === scroll ? (
                <SettingsWrapperScroll>{child}</SettingsWrapperScroll>
            ) : (
                <Box>{child}</Box>
            ))}
        </Box>
    )
}
