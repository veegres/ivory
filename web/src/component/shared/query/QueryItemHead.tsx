import select from "../../../style/select.module.css";
import {Box} from "@mui/material";
import {SxPropsMap} from "../../../type/common";
import {ReactNode} from "react";

const SX: SxPropsMap = {
    head: {display: "flex", padding: "5px 15px"},
    title: {flexGrow: 1, display: "flex", alignItems: "center", gap: 1},
    buttons: {display: "flex", alignItems: "center", gap: 1},
}


type Props = {
    renderTitle: ReactNode,
    renderButtons: ReactNode,
    onClick?: () => void,
}

export function QueryItemHead(props: Props) {
    const {renderButtons, renderTitle, onClick} = props
    return (
        <Box sx={SX.head} className={select.none}>
            <Box sx={SX.title} onClick={onClick}>
                {renderTitle}
            </Box>
            <Box sx={SX.buttons}>
                {renderButtons}
            </Box>
        </Box>
    )
}
