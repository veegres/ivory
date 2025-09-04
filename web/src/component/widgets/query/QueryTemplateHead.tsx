import select from "../../../style/select.module.css";
import {Box} from "@mui/material";
import {SxPropsMap} from "../../../type/general";
import {ReactNode, useState} from "react";

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
    showAllButtons: boolean,
    renderHiddenButtons: ReactNode,
    renderButtons: ReactNode,
}

export function QueryTemplateHead(props: Props) {
    const {renderButtons, renderTitle, renderHiddenButtons, showAllButtons} = props
    const [showMouse, setShowMouse] = useState(false)

    return (
        <Box
            sx={SX.head}
            className={select.none}
            onMouseEnter={() => setShowMouse(true)}
            onMouseLeave={() => setShowMouse(false)}
        >
            <Box sx={SX.title}>
                {renderTitle}
            </Box>
            <Box sx={SX.buttons}>
                {(showMouse || showAllButtons) && renderHiddenButtons}
                {renderButtons}
            </Box>
        </Box>
    )
}
