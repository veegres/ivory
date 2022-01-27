import {InputLabel} from "@mui/material";

export function Error({ error }: { error: any }) {
    return (
        <InputLabel error={true} sx={{ padding: '20px', whiteSpace: 'pre-wrap' }}>
            {JSON.stringify(error, null, 4)}
        </InputLabel>
    )
}
