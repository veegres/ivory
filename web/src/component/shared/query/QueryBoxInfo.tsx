import {SxPropsMap} from "../../../type/common";
import {QueryBoxWrapper} from "./QueryBoxWrapper";
import {Box} from "@mui/material";
import {ReactNode} from "react";
import {QueryType} from "../../../type/query";
import {InfoBox} from "../../view/box/InfoBox";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    type: {color: "secondary.main"},
    options: {display: "flex", gap: 1},
    params: {flexGrow: 1},
}

type Props = {
    type: QueryType,
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
            <Box sx={SX.options}>
                <QueryBoxWrapper sx={SX.params} editable={editable}>
                    {props.renderParams}
                </QueryBoxWrapper>
                <QueryBoxWrapper editable={editable}>
                    {props.renderVarieties}
                </QueryBoxWrapper>
                <QueryBoxWrapper editable={editable}>
                    <InfoBox tooltip={"Type"} withPadding>
                        <Box sx={SX.type}>{QueryType[type]}</Box>
                    </InfoBox>
                </QueryBoxWrapper>
            </Box>
            <QueryBoxWrapper editable={editable}>
                {props.renderDescription}
            </QueryBoxWrapper>
            <QueryBoxWrapper editable={editable}>
                {props.renderQuery}
            </QueryBoxWrapper>
        </Box>
    )
}
