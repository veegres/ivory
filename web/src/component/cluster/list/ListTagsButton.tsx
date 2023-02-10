import {ToggleButton} from "@mui/material";
import {useTheme} from "../../../provider/ThemeProvider";
import {SxPropsMap} from "../../../app/types";

const SX: SxPropsMap = {
    element: {padding: "3px 7px", borderRadius: "3px", lineHeight: "1.1"},
}

type Props = {
    tag: string,
    selected: boolean
    onClick?: (value: string) => void,
}
export function ListTagsButton(props: Props) {
    const { tag, selected, onClick } = props
    const { info } = useTheme()

    return (
        <ToggleButton
            sx={{...SX.element, border: `1px solid ${info?.palette.divider}`}}
            color={"secondary"}
            size={"small"}
            key={tag}
            selected={selected}
            disabled={!onClick}
            value={tag}
            onClick={() => onClick!!(tag)}
        >
            {tag}
        </ToggleButton>
    )
}
