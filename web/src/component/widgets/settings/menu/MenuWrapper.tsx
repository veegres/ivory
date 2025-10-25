import {Box} from "@mui/material";
import {Children, ReactNode} from "react";

import {SxPropsMap} from "../../../../app/type";
import {MenuWrapperScroll} from "./MenuWrapperScroll";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1, height: "100%", overflow: "hidden"},
}

type Props = {
    scrollElement?: number,
    children: ReactNode,
}

export function MenuWrapper(props: Props) {
    const {children, scrollElement} = props
    const scroll = scrollElement ?? Children.count(children) - 1
    return (
        <Box sx={SX.box}>
            {Children.map(children, (child, index) => index === scroll ? (
                <MenuWrapperScroll>{child}</MenuWrapperScroll>
            ) : (
                <Box>{child}</Box>
            ))}
        </Box>
    )
}
