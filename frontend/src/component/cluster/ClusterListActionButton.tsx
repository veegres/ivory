import {cloneElement, ReactElement} from "react";
import {CircularProgress, IconButton} from "@mui/material";

const SX = {
    icon: { fontSize: 18 },
    button: { height: '32px', width: '32px' },
}

export function ClusterListActionButton(props: { icon: ReactElement, loading: boolean, disabled: boolean, onClick: () => void }) {
    const { loading, icon, onClick, disabled } = props;
    return (
        <IconButton sx={SX.button} disabled={loading || disabled} onClick={onClick}>
            {loading ?
                <CircularProgress size={SX.icon.fontSize} /> :
                cloneElement(icon, { sx: SX.icon })
            }
        </IconButton>
    )
}
