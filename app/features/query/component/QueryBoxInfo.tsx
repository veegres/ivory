import {Box} from "@mui/material"
import {ReactNode} from "react"

import {InfoBox, Padding} from "../../../shared/component/box/InfoBox"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {Type} from "../api/QueryType"
import {QueryBoxWrapper} from "./QueryBoxWrapper"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    type: {color: "secondary.main"},
    options: {display: "flex", gap: 1},
    params: {flexGrow: 1},
}

type Props = {
    type: Type,
    editable: boolean,
    renderVarieties: ReactNode,
    renderDescription: ReactNode,
    renderParams: ReactNode,
    renderQuery: ReactNode,
}

export function QueryBoxInfo(props: Props) {
    const {type, editable} = props
    return (
        <Box sx={SX.box}>
            <QueryBoxWrapper editable={editable}>
                {props.renderDescription}
            </QueryBoxWrapper>
            <Box sx={SX.options}>
                <QueryBoxWrapper sx={SX.params} editable={editable}>
                    {props.renderParams}
                </QueryBoxWrapper>
                <QueryBoxWrapper editable={editable}>
                    {props.renderVarieties}
                </QueryBoxWrapper>
                <QueryBoxWrapper editable={editable}>
                    <InfoBox tooltip={"Type"} padding={Padding.Small} height={editable ? "26px" : undefined}>
                        <Box sx={SX.type}>{Type[type]}</Box>
                    </InfoBox>
                </QueryBoxWrapper>
            </Box>
            <QueryBoxWrapper editable={editable}>
                {props.renderQuery}
            </QueryBoxWrapper>
        </Box>
    )
}
