import select from "../../../style/select.module.css";
import {Box} from "@mui/material";
import {SxPropsMap} from "../../../type/common";
import {ReactNode} from "react";

const SX: SxPropsMap = {
    head: {display: "flex", padding: "5px 15px"},
    title: {
        flexGrow: 1, display: "flex", alignItems: "center", gap: 1,
        whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis",
    },
    buttons: {display: "flex", alignItems: "center"},
}


type Props = {
    renderTitle: ReactNode,
    renderButtons: ReactNode,
    onClick?: () => void,
    onMouseOver?: () => void,
    onMouseOut?: () => void,
}

export function QueryHead(props: Props) {
    const {renderButtons, renderTitle, onClick, onMouseOver, onMouseOut} = props
    return (
        <Box sx={SX.head} className={select.none} onMouseOver={onMouseOver} onMouseOut={onMouseOut}>
            <Box sx={SX.title} onClick={onClick}>
                {renderTitle}
            </Box>
            <Box sx={SX.buttons}>
                {renderButtons}
            </Box>
        </Box>
    )
}
